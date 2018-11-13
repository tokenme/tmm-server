package redeem

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ethgasstation-api"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/tools/wechatpay"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
)

type TMMWithdrawRequest struct {
	TMM      decimal.Decimal `json:"tmm" form:"tmm" binding:"required"`
	Currency string          `json:"currency" form:"currency" binding:"required"`
}

func TMMWithdrawHandler(c *gin.Context) {
	if CheckWithCode(true, FEATURE_NOT_AVAILABLE_ERROR, "feature not available", c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req TMMWithdrawRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT union_id FROM tmm.wx WHERE user_id=%d LIMIT 1`, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, WECHAT_UNAUTHORIZED_ERROR, "Wechat Unauthorized", c) {
		return
	}
	wxUnionId := rows[0].Str(0)

	recyclePrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	minTmmRequire := decimal.New(int64(Config.MinTMMRedeem), 0)
	if req.TMM.LessThan(minTmmRequire) {
		c.JSON(INVALID_MIN_TOKEN_ERROR, common.ExchangeRate{MinPoints: minTmmRequire})
		return
	}

	cash := req.TMM.Mul(recyclePrice)

	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if CheckErr(err, c) {
		return
	}
	tmmInt := req.TMM.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if Check(!ok, fmt.Sprintf("Internal Error: %s", tmmInt.Floor().String()), c) {
		return
	}

	tokenBalance, err := utils.TokenBalanceOf(token, user.Wallet)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(amount.Cmp(tokenBalance) == 1, NOT_ENOUGH_TOKEN_ERROR, "not enough token", c) {
		return
	}

	tradeNumToken, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}

	forexRate := forex.Rate(Service, "USD", "CNY")
	cny := cash.Mul(forexRate)

	payClient := wechatpay.NewClient(Config.Wechat.AppId, Config.Wechat.MchId, Config.Wechat.Key, Config.Wechat.CertCrt, Config.Wechat.CertKey)
	payParams := &wechatpay.Request{
		TradeNum:    commonutils.Md5(tradeNumToken.String()),
		Amount:      cny.Mul(decimal.New(100, 0)).IntPart(),
		CallbackURL: fmt.Sprintf("%s/wechat/pay/callback", Config.BaseUrl),
		OpenId:      wxUnionId,
		Ip:          ClientIP(c),
		Desc:        "UCoinWithdraw",
	}
	payParams.Nonce = commonutils.Md5(payParams.TradeNum)
	payRes, err := payClient.Pay(payParams)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	spew.Dump(payRes)
	if CheckWithCode(payRes.ErrCode != "", WECHAT_PAYMENT_ERROR, payRes.ErrCodeDesc, c) {
		return
	}
	return
	poolPrivKey, err := commonutils.AddressDecrypt(Config.TMMPoolWallet.Data, Config.TMMPoolWallet.Salt, Config.TMMPoolWallet.Key)
	if CheckErr(err, c) {
		return
	}
	poolPubKey, err := eth.AddressFromHexPrivateKey(poolPrivKey)
	if CheckErr(err, c) {
		return
	}

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if CheckErr(err, c) {
		return
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if CheckErr(err, c) {
		return
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, GlobalLock, agentPubKey, Config.Geth)
	if CheckErr(err, c) {
		return
	}
	var gasPrice *big.Int
	gas, err := ethgasstation.Gas()
	if err != nil {
		gasPrice = nil
	} else {
		gasPrice = new(big.Int).Mul(big.NewInt(gas.SafeLow.Div(decimal.New(10, 0)).IntPart()), big.NewInt(params.GWei))
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	_, err = utils.TransferProxy(token, transactor, user.Wallet, poolPubKey, amount)
	if CheckErr(err, c) {
		return
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, GlobalLock, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	/*
		_, _, err = db.Query(`INSERT INTO tmm.exchange_records (tx, status, user_id, device_id, tmm, points, direction) VALUES ('%s', %d, %d, '%s', '%s', '%s', %d)`, db.Escape(receipt.Receipt), receipt.Status, user.Id, db.Escape(req.DeviceId), receipt.Value.String(), req.Points.String(), req.Direction)
		if CheckErr(err, c) {
			return
		}
		c.JSON(http.StatusOK, receipt)
	*/
}
