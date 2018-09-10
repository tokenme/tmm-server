package exchange

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
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
	"time"
)

type TMMChangeRequest struct {
	Points   decimal.Decimal `json:"points" form:"points" binding:"required"`
	DeviceId string          `json:"device_id" form:"device_id" binding:"required"`
}

func TMMChangeHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req TMMChangeRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	db := Service.Db
	query := `SELECT
    d.points
FROM tmm.devices AS d
INNER JOIN tmm.user_devices AS du ON (du.device_id=d.id)
WHERE d.id='%s' AND du.user_id=%d`
	rows, _, err := db.Query(query, db.Escape(req.DeviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	log.Warn(query, db.Escape(req.DeviceId), user.Id)
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	points, err := decimal.NewFromString(rows[0].Str(0))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(points.LessThan(req.Points), NOT_ENOUGH_POINTS_ERROR, "not enough points", c) {
		return
	}
	exchangeRate, err := common.GetExchangeRate(Config, Service)
	if CheckErr(err, c) {
		return
	}
	if req.Points.LessThan(exchangeRate.MinPoints) {
		c.JSON(INVALID_MIN_POINTS_ERROR, exchangeRate)
		return
	}
	tmm := req.Points.Mul(exchangeRate.Rate)
	poolPrivKey, err := commonutils.AddressDecrypt(Config.TMMPoolWallet.Data, Config.TMMPoolWallet.Salt, Config.TMMPoolWallet.Key)
	if CheckErr(err, c) {
		return
	}
	poolPubKey, err := eth.AddressFromHexPrivateKey(poolPrivKey)
	if CheckErr(err, c) {
		return
	}
	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if CheckErr(err, c) {
		return
	}
	tmmInt := tmm.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if Check(!ok, fmt.Sprintf("Internal Error: %s", tmmInt.Floor().String()), c) {
		return
	}
	transactor := eth.TransactorAccount(poolPrivKey)
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, poolPubKey, Config.Geth)
	if CheckErr(err, c) {
		return
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Shannon)),
		GasLimit: 540000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.Transfer(token, transactor, user.Wallet, amount)
	if CheckErr(err, c) {
		return
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, poolPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	receipt := common.Transaction{
		Receipt:    tx.Hash().Hex(),
		Value:      tmm,
		Status:     2,
		InsertedAt: time.Now().Format(time.RFC3339),
	}
	_, _, err = db.Query(`INSERT INTO tmm.txs (tx, status, user_id, device_id, tmm, points) VALUES ('%s', %d, %d, '%s', '%s', '%s')`, db.Escape(receipt.Receipt), receipt.Status, user.Id, db.Escape(req.DeviceId), receipt.Value.String(), req.Points.String())
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.user_devices AS du SET d.points = d.points - %s WHERE du.device_id=d.id AND d.id='%s' AND du.user_id=%d AND d.points > %s`, req.Points.String(), db.Escape(req.DeviceId), user.Id, req.Points.String())
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, receipt)
}
