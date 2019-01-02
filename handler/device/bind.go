package device

import (
	//"github.com/davecgh/go-spew/spew"
	"errors"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
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
	inviteBonus(user, deviceId, c)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}

func inviteBonus(user common.User, deviceId string, c *gin.Context) error {
	db := Service.Db
	_, ret, err := db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2, tmm.invite_submissions AS iss, ucoin.users AS u SET t1.parent_id=t2.user_id, t1.grand_id=t2.parent_id WHERE (t1.parent_id!=t2.user_id OR t1.grand_id!=t2.parent_id) AND t2.user_id!=t1.user_id AND t2.parent_id!=t1.user_id AND t2.id != t1.id AND t2.id=iss.code AND t1.user_id=u.id AND iss.completed=0 AND u.country_code=86 AND iss.tel=u.mobile AND u.mobile='%s'`, db.Escape(user.Mobile))
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return err
	}
	if ret.AffectedRows() == 0 {
		return nil
	}

	_, _, err = db.Query(`UPDATE tmm.invite_submissions iss, ucoin.users AS u SET iss.completed=1 WHERE iss.tel=u.mobile AND u.country_code=86 AND u.mobile='%s'`, db.Escape(user.Mobile))
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return err
	}

	query := `SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0 AND (IFNULL(us.blocked,0)=0 OR us.block_whitelist=1)
ORDER BY d.lastping_at DESC LIMIT 1`
	rows, _, err := db.Query(query, user.Id)
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	/*
		inviterCashBonus := decimal.New(int64(Config.InviterCashBonus), 0)
		pointPrice := common.GetPointPrice(Service, Config)
		forexRate := forex.Rate(Service, "USD", "CNY")
		pointCnyPrice := pointPrice.Mul(forexRate)
			    inviterPointBonus := inviterCashBonus.Div(pointCnyPrice)
				maxInviterBonus := decimal.New(Config.MaxInviteBonus, 0)
				if inviterPointBonus.GreaterThanOrEqual(maxInviterBonus) {
					inviterPointBonus = maxInviterBonus
				}
	*/
	inviterPointBonus := decimal.New(int64(Config.InviterBonus), 0)
	inviterDeviceId := rows[0].Str(0)
	inviterUserId := rows[0].Uint64(1)
	pointsPerTs, _ := common.GetPointsPerTs(Service)
	inviterTs := inviterPointBonus.Div(pointsPerTs)
	forexRate := forex.Rate(Service, "USD", "CNY")
	tx, tmm, err := _transferToken(user.Id, forexRate, c)
	if err != nil {
		log.Error("Bonus Transfer failed")
		raven.CaptureError(err, nil)
		return err
	}
	//log.Warn("Inviter bonus: %s, inviter:%d", inviterPointBonus.String(), inviterUserId)
	_, _, err = db.Query(`UPDATE tmm.devices AS d2 SET d2.points = d2.points + %s, d2.total_ts = d2.total_ts + %d WHERE d2.id='%s'`, inviterPointBonus.String(), inviterTs.IntPart(), db.Escape(inviterDeviceId))
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return err
	}
	_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, tmm, tmm_tx) VALUES (%d, %d, 0, %s, '%s'), (%d, %d, %s, 0, '')`, user.Id, user.Id, tmm.String(), db.Escape(tx), inviterUserId, user.Id, inviterPointBonus.String())
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return err
	}
	_, _, err = db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.grand_id=t2.parent_id WHERE t2.user_id=t1.parent_id AND t2.parent_id!=t1.user_id AND t2.user_id=%d`, user.Id)
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
	}
	_, _, err = db.Query(`INSERT INTO user_settings (user_id, level)
(
SELECT
i.parent_id, ul.id
FROM tmm.user_levels AS ul
INNER JOIN (
    SELECT ic.parent_id, COUNT(DISTINCT ic.user_id) AS invites
    FROM tmm.invite_codes AS ic
    LEFT JOIN tmm.user_settings AS us ON (us.user_id=ic.user_id)
    WHERE ic.parent_id=%d AND (IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1)
) AS i ON (i.invites >= ul.invites AND parent_id IS NOT NULL)
ORDER BY ul.id DESC LIMIT 1
) ON DUPLICATE KEY UPDATE level=VALUES(level)`, inviterUserId)
	if err != nil {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
	}
	return err
}

func _transferToken(userId uint64, forexRate decimal.Decimal, c *gin.Context) (receipt string, tokenAmount decimal.Decimal, err error) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT u.wallet_addr FROM ucoin.users AS u LEFT JOIN tmm.user_settings AS us ON (us.user_id=u.id)  WHERE u.id=%d AND (IFNULL(us.blocked,0)=0 OR us.block_whitelist=1)`, userId)
	if err != nil {
		return receipt, tokenAmount, err
	}
	if len(rows) == 0 {
		return receipt, tokenAmount, errors.New("not found")
	}
	userWallet := rows[0].Str(0)
	tokenPrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	tokenPriceCny := tokenPrice.Mul(forexRate)
	tokenAmount = decimal.New(3, 0).Div(tokenPriceCny)

	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if err != nil {
		return receipt, tokenAmount, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return receipt, tokenAmount, err
	}
	tmmInt := tokenAmount.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if !ok {
		return receipt, tokenAmount, nil
	}

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if err != nil {
		return receipt, tokenAmount, err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return receipt, tokenAmount, err
	}

	tokenBalance, err := utils.TokenBalanceOf(token, agentPubKey)
	if err != nil {
		return receipt, tokenAmount, err
	}
	if amount.Cmp(tokenBalance) == 1 {
		return receipt, tokenAmount, nil
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		return receipt, tokenAmount, err
	}
	gasPrice, err := Service.Geth.SuggestGasPrice(c)
	if err == nil && gasPrice.Cmp(eth.MinGas) == -1 {
		gasPrice = eth.MinGas
	} else {
		gasPrice = nil
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.Transfer(token, transactor, userWallet, amount)
	if err != nil {
		return receipt, tokenAmount, err
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	receipt = tx.Hash().Hex()
	return receipt, tokenAmount, nil
}
