package bonus

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"net/http"
)

type DailyCommitRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
	Platform common.Platform `json:"platform" form:"platform"`
}

func DailyCommitHandler(c *gin.Context) {
	var req DailyCommitRequest
	if CheckErr(c.Bind(&req), c) {
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

	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if CheckErr(user.IsBlocked(Service), c) {
		log.Error("Blocked User:%d", user.Id)
		return
	}
	db := Service.Db
	_, ret, err := db.Query(`INSERT INTO tmm.daily_bonus_logs (user_id, updated_on, days) VALUES (%d, NOW(), 1) ON DUPLICATE KEY UPDATE days=IF(updated_on=DATE(DATE_SUB(NOW(), INTERVAL 1 DAY)) AND days<7, days+1, IF(updated_on=DATE(NOW()), days, 1)), updated_on=VALUES(updated_on)`, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(ret.AffectedRows() == 0, DAILY_BONUS_COMMITTED_ERROR, "Already checked in", c) {
		return
	}
	rows, _, err := db.Query(`SELECT days FROM tmm.daily_bonus_logs WHERE user_id=%d LIMIT 1`, user.Id)
	if CheckErr(err, c) {
		return
	}
	days := rows[0].Int64(0)
	points := decimal.New(days, 0)
	pointsPerTs, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := points.Div(pointsPerTs)
	_, _, err = db.Query(`UPDATE tmm.devices SET points=points+%s, total_ts=total_ts+%s WHERE id='%s' AND user_id=%d`, points.String(), ts.String(), db.Escape(deviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	var interests decimal.Decimal
	if points.GreaterThan(decimal.Zero) {
		_, interests, _ = giveDailyInterests(c, user)
	}
	c.JSON(http.StatusOK, gin.H{"days": days, "points": points, "interests": interests})
}

func giveDailyInterests(c *gin.Context, user common.User) (origin decimal.Decimal, interests decimal.Decimal, err error) {
	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if err != nil {
		return origin, interests, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		log.Error(err.Error())
		return origin, interests, err
	}
	balance, err := utils.TokenBalanceOf(token, user.Wallet)
	if err != nil {
		log.Error(err.Error())
		return origin, interests, err
	}
	origin = decimal.NewFromBigInt(balance, -1*int32(tokenDecimal))
	interests = origin.Mul(decimal.NewFromFloat(Config.DailyTMMInterestsRate))
	if interests.LessThanOrEqual(decimal.New(1, -4)) {
		log.Error("Interests: %s", interests.String())
		return
	}

	tmmInt := interests.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if !ok {
		return
	}
	poolPrivKey, err := commonutils.AddressDecrypt(Config.TMMPoolWallet.Data, Config.TMMPoolWallet.Salt, Config.TMMPoolWallet.Key)
	if err != nil {
		log.Error(err.Error())
		return
	}
	poolPubKey, err := eth.AddressFromHexPrivateKey(poolPrivKey)
	if err != nil {
		log.Error(err.Error())
		return
	}

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if err != nil {
		log.Error(err.Error())
		return
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		log.Error(err.Error())
		return
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var gasPrice *big.Int
	/*gas, err := ethgasstation.Gas()
	  if err != nil {
	      gasPrice = nil
	  } else {
	      gasPrice = new(big.Int).Mul(big.NewInt(gas.SafeLow.Div(decimal.New(10, 0)).IntPart()), big.NewInt(params.GWei))
	  }*/
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.TransferProxy(token, transactor, poolPubKey, user.Wallet, amount)
	if err != nil {
		log.Error(err.Error())
		return
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	db := Service.Db
	_, _, err = db.Query(`INSERT INTO tmm.interests (tx, status, user_id, balance, interest) VALUES ('%s', 2, %d, %s, %s)`, db.Escape(tx.Hash().Hex()), user.Id, origin.String(), interests.String())
	if err != nil {
		log.Error(err.Error())
	}
	return
}
