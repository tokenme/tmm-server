package user

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"net/http"
	"strings"
)

type UpdateRequest struct {
	Mobile        string           `form:"mobile" json:"mobile"`
	CountryCode   uint             `form:"country_code" json:"country_code"`
	VerifyCode    string           `form:"verify_code" json:"verify_code"`
	Password      string           `form:"passwd" json:"passwd"`
	RePassword    string           `form:"repasswd" json:"repasswd"`
	Realname      string           `form:"realname" json:"realname"`
	PaymentPasswd string           `form:"payment_passwd" json:"payment_passwd"`
	InviterCode   tokenUtils.Token `form:"inviter_code" json:"inviter_code"`
}

func UpdateHandler(c *gin.Context) {
	var req UpdateRequest
	if CheckWithCode(c.Bind(&req) != nil, BADREQUEST_ERROR, "invalid request", c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db

	var updateFields []string
	if req.Realname != "" {
		updateFields = append(updateFields, fmt.Sprintf("realname='%s'", db.Escape(req.Realname)))
	}
	var (
		passwd string
		mobile string
	)
	if req.Mobile != "" {
		if Check(req.Mobile == "" || req.CountryCode == 0 || req.VerifyCode == "" || req.Password == "" || req.RePassword == "", "missing params", c) {
			return
		}
		if Check(req.Password != req.RePassword, "repassword!=password", c) {
			return
		}
		passwdLength := len(req.Password)
		if Check(passwdLength < 8 || passwdLength > 64, "password length must between 8-32", c) {
			return
		}
		token, err := uuid.NewV4()
		if CheckErr(err, c) {
			return
		}
		salt := utils.Sha1(token.String())
		passwd = utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
		mobile = strings.Replace(req.Mobile, " ", "", 0)
		rows, _, err := db.Query(`SELECT 1 FROM ucoin.auth_verify_codes WHERE country_code=%d AND mobile='%s' AND code='%s' LIMIT 1`, req.CountryCode, db.Escape(mobile), db.Escape(req.VerifyCode))
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if Check(len(rows) == 0, "unverified phone number", c) {
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("country_code=%d, mobile='%s', salt='%s', passwd='%s'", req.CountryCode, db.Escape(mobile), db.Escape(salt), db.Escape(passwd)))
	} else if req.PaymentPasswd != "" {
		paymentPasswd := utils.Sha1(fmt.Sprintf("%s%s%s", user.Salt, req.PaymentPasswd, user.Salt))
		updateFields = append(updateFields, fmt.Sprintf("payment_passwd='%s'", db.Escape(paymentPasswd)))
	} else if req.InviterCode > 0 {
		query := `SELECT
d.id
FROM tmm.user_devices AS du
INNER JOIN tmm.devices AS d ON (d.id = du.device_id)
WHERE du.user_id = %d
ORDER BY d.lastping_at DESC LIMIT 1`
		rows, _, err := db.Query(query, user.Id)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		deviceId := rows[0].Str(0)
		_, ret, err := db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.parent_id=t2.user_id, t1.grand_id=t2.parent_id WHERE t2.id != t1.id AND t2.id=%d AND t1.user_id=%d`, req.InviterCode, user.Id)
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return
		}
		if CheckWithCode(ret.AffectedRows() == 0, INVALID_INVITE_CODE_ERROR, "invalid invite code", c) {
			return
		}
		query = `SELECT
d.id,
du.user_id
FROM tmm.user_devices AS du
INNER JOIN tmm.devices AS d ON (d.id = du.device_id)
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=du.user_id)
WHERE ic.user_id = %d
ORDER BY d.lastping_at DESC LIMIT 1`
		rows, _, err = db.Query(query, user.Id)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		inviterDeviceId := rows[0].Str(0)
		inviterUserId := rows[0].Uint64(1)
		_, ret2, err := db.Query(`UPDATE tmm.devices AS d1, tmm.devices AS d2 SET d1.points = d1.points + %d, d2.points = d2.points + %d WHERE d1.id='%s' AND d2.id='%s'`, Config.InviteBonus, Config.InviterBonus, db.Escape(deviceId), db.Escape(inviterDeviceId))
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return
		}
		if ret2.AffectedRows() > 0 {
			_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus) VALUES (%d, %d, %d), (%d, %d, %d)`, user.Id, user.Id, Config.InviteBonus, inviterUserId, user.Id, Config.InviterBonus)
			if err != nil {
				log.Error(err.Error())
				raven.CaptureError(err, nil)
			}
		}
		_, _, err = db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.grand_id=t2.parent_id WHERE t2.user_id=t1.parent_id AND t2.user_id=%d`, user.Id)
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return
		}
	}
	if len(updateFields) == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}

	_, _, err := db.Query(`UPDATE ucoin.users SET %s WHERE id=%d LIMIT 1`, strings.Join(updateFields, ","), user.Id)
	if CheckErr(err, c) {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
