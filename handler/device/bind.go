package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"net/http"
)

type BindRequest struct {
	Idfa     string `form:"idfa" json:"idfa"`
	Platform string `form:"platform" json:"platform"`
	Imei     string `form:"imei" json:"imei"`
	Mac      string `form:"mac" json:"mac"`
}

func BindHandler(c *gin.Context) {
	var req common.DeviceRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	err := saveDevice(Service, req, c)
	if CheckErr(err, c) {
		return
	}
	err = saveApp(Service, req)
	if CheckErr(err, c) {
		return
	}

	db := Service.Db

	if Check(req.Idfa == "" && req.Imei == "" && req.Mac == "", "invalid request", c) {
		return
	}
	rows, _, err := db.Query(`SELECT COUNT(*) FROM tmm.devices WHERE user_id=%d`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var deviceCount int
	if len(rows) > 0 {
		deviceCount = rows[0].Int(0)
	}
	if CheckWithCode(deviceCount >= Config.MaxBindDevice, MAX_BIND_DEVICE_ERROR, "exceeded maximum binding devices", c) {
		return
	}
	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}
	_, ret, err := db.Query(`UPDATE tmm.devices SET user_id=%d WHERE id='%s' AND user_id=0`, user.Id, deviceId)
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() == 0 {
		rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, deviceId)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, OTHER_BIND_DEVICE_ERROR, "the device has been bind by others", c) {
			return
		}
	}
	inviteBonus(user, deviceId)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}

func inviteBonus(user common.User, deviceId string) error {
	db := Service.Db
	_, ret, err := db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2, tmm.invite_submissions AS iss, ucoin.users AS u SET t1.parent_id=t2.user_id, t1.grand_id=t2.parent_id WHERE (t1.parent_id!=t2.user_id OR t1.grand_id!=t2.parent_id) AND t2.user_id!=t1.user_id AND t2.parent_id!=t1.user_id AND t2.id != t1.id AND t2.id=iss.code AND t1.user_id=u.id AND iss.completed=0 AND u.country_code=86 AND iss.tel=u.mobile AND u.mobile='%s'`, db.Escape(user.Mobile))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if ret.AffectedRows() == 0 {
		return nil
	}

	_, _, err = db.Query(`UPDATE tmm.invite_submissions iss, ucoin.users AS u SET iss.completed=1 WHERE iss.tel=u.mobile AND u.country_code=86 AND u.mobile='%s'`, db.Escape(user.Mobile))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	query := `SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
ORDER BY d.lastping_at DESC LIMIT 1`
	rows, _, err := db.Query(query, user.Id)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	inviterCashBonus := decimal.New(int64(Config.InviterCashBonus), 0)
	pointPrice := common.GetPointPrice(Service, Config)
	forexRate := forex.Rate(Service, "USD", "CNY")
	pointCnyPrice := pointPrice.Mul(forexRate)
	inviterPointBonus := inviterCashBonus.Div(pointCnyPrice)
	maxInviterBonus := decimal.New(4000, 0)
	if inviterPointBonus.GreaterThanOrEqual(maxInviterBonus) {
		inviterPointBonus = maxInviterBonus
	}

	inviterDeviceId := rows[0].Str(0)
	inviterUserId := rows[0].Uint64(1)
	pointsPerTs, _ := common.GetPointsPerTs(Service)
	inviteTs := decimal.New(int64(Config.InviteBonus), 0).Div(pointsPerTs)
	inviterTs := inviterPointBonus.Div(pointsPerTs)
	//log.Warn("Inviter bonus: %s, inviter:%d", inviterPointBonus.String(), inviterUserId)
	_, ret2, err := db.Query(`UPDATE tmm.devices AS d1, tmm.devices AS d2 SET d1.points = d1.points + %d, d1.total_ts = d1.total_ts + %d, d2.points = d2.points + %s, d2.total_ts = d2.total_ts + %d WHERE d1.id='%s' AND d2.id='%s'`, Config.InviteBonus, inviteTs.IntPart(), inviterPointBonus.String(), inviterTs.IntPart(), db.Escape(deviceId), db.Escape(inviterDeviceId))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if ret2.AffectedRows() > 0 {
		_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus) VALUES (%d, %d, %d), (%d, %d, %s)`, user.Id, user.Id, Config.InviteBonus, inviterUserId, user.Id, inviterPointBonus.String())
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	_, _, err = db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.grand_id=t2.parent_id WHERE t2.user_id=t1.parent_id AND t2.parent_id!=t1.user_id AND t2.user_id=%d`, user.Id)
	if err != nil {
		log.Error(err.Error())
	}
	_, _, err = db.Query(`INSERT INTO user_settings (user_id, level)
(
SELECT
i.parent_id, ul.id
FROM tmm.user_levels AS ul
INNER JOIN (
    SELECT parent_id, COUNT(*) AS invites FROM tmm.invite_codes WHERE parent_id=%d
) AS i ON (i.invites >= ul.invites)
ORDER BY ul.id DESC LIMIT 1
) ON DUPLICATE KEY UPDATE level=VALUES(level)`, inviterUserId)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
