package redeem

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	//"github.com/ethereum/go-ethereum/params"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nlopes/slack"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	//"github.com/tokenme/tmm/tools/ethgasstation-api"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/tools/ykt"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type TMMWithdrawRequest struct {
	TMM      decimal.Decimal `json:"tmm" form:"tmm" binding:"required"`
	Currency string          `json:"currency" form:"currency" binding:"required"`
}

const TMMWithdrawRateKey = "TMMWithdrawRate-%d"

func TMMWithdrawHandler(c *gin.Context) {
	/*
		if CheckWithCode(true, FEATURE_NOT_AVAILABLE_ERROR, "feature not available", c) {
			return
		}
	*/
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
	{
		rows, _, err := db.Query(`SELECT 1 FROM tmm.user_settings WHERE user_id=%d AND blocked=1 AND block_whitelist=0`, user.Id)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) > 0, "您的账户存在异常操作，疑似恶意邀请用户行为，不能执行提现操作。如有疑问请联系客服。", c) {
			return
		}
	}
	rows, _, err := db.Query(`SELECT wx.union_id, oi.open_id FROM tmm.wx LEFT JOIN tmm.wx_openids AS oi ON (oi.union_id=wx.union_id AND oi.app_id='%s') WHERE wx.user_id=%d LIMIT 1`, db.Escape(Config.Wechat.AppId), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, WECHAT_UNAUTHORIZED_ERROR, "Wechat Unauthorized", c) {
		return
	}
	wxUnionId := rows[0].Str(0)
	wxOpenId := rows[0].Str(1)
	if wxOpenId == "" {
		openIdReq := ykt.OpenIdRequest{
			UnionId: wxUnionId,
		}
		openIdRes, err := openIdReq.Run()
		if CheckWithCode(err != nil, WECHAT_OPENID_ERROR, "need openid", c) {
			log.Error(err.Error())
			return
		}
		wxOpenId = openIdRes.Data.OpenId
		_, _, err = db.Query(`INSERT INTO tmm.wx_openids (app_id, open_id, union_id) VALUES ('%s', '%s', '%s')`, db.Escape(Config.Wechat.AppId), db.Escape(wxOpenId), db.Escape(wxUnionId))
		if CheckErr(err, c) {
			return
		}
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	withdrawRateKey := fmt.Sprintf(TMMWithdrawRateKey, user.Id)
	withdrawTime, err := redis.String(redisConn.Do("GET", withdrawRateKey))
	if CheckWithCode(err == nil, TOKEN_WITHDRAW_RATE_LIMIT_ERROR, "每次提现时间间隔不能少于24小时", c) {
		log.Warn("WithdrawRateLimit: %d, time: %s", user.Id, withdrawTime)
		return
	}
	recyclePrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	minTmmRequire := decimal.New(int64(Config.MinTMMRedeem), 0)
	if req.TMM.LessThan(minTmmRequire) {
		c.JSON(INVALID_MIN_TOKEN_ERROR, common.ExchangeRate{MinPoints: minTmmRequire})
		return
	}

	if CheckWithCode(req.TMM.LessThan(minTmmRequire) || req.TMM.GreaterThan(decimal.New(10000, 0)), INVALID_TOKEN_WITHDRAW_AMOUNT_ERROR, fmt.Sprintf("提现UCoin超出限制。最小%s UCoin或累计超过10000 UCoin", minTmmRequire.String()), c) {
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

	forexRate := forex.Rate(Service, "USD", "CNY")
	cny := cash.Mul(forexRate)

	if CheckWithCode(cny.LessThan(decimal.New(1, 0)) || cny.GreaterThan(decimal.New(2000, 0)), WECHAT_PAYMENT_ERROR, "提现金额超出限制。最小金额1元或累计超过2000元", c) {
		log.Error("cash: %s, tmm: %s, recyclePrice:%s", cny.String(), req.TMM.String(), recyclePrice.String())
		return
	}

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
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if CheckErr(err, c) {
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
	tx, err := utils.TransferProxy(token, transactor, user.Wallet, poolPubKey, amount)
	if CheckErr(err, c) {
		return
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	_, _, err = db.Query(`INSERT INTO tmm.withdraw_txs (tx, user_id, tmm, cny, client_ip) VALUES ('%s', %d, %s, %s, '%s')`, tx.Hash().Hex(), user.Id, req.TMM.String(), cny.String(), ClientIP(c))
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	_, err = redisConn.Do("SETEX", withdrawRateKey, 60*60*24, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error(err.Error())
	}
	receipt := common.TMMWithdrawResponse{
		TMM:      req.TMM,
		Cash:     cny,
		Currency: req.Currency,
	}
	slackParams := slack.PostMessageParameters{Parse: "full", UnfurlMedia: true, Markdown: true}
	attachment := slack.Attachment{
		Color:      "success",
		AuthorName: user.ShowName,
		AuthorIcon: user.Avatar,
		Title:      "Token提现",
		Fallback:   "Fallback message",
		Fields: []slack.AttachmentField{
			{
				Title: "CountryCode",
				Value: strconv.FormatUint(uint64(user.CountryCode), 10),
				Short: true,
			},
			{
				Title: "UserID",
				Value: strconv.FormatUint(user.Id, 10),
				Short: true,
			},
			{
				Title: "CNY",
				Value: cny.StringFixed(2),
				Short: true,
			},
			{
				Title: "Tokens",
				Value: req.TMM.StringFixed(4),
				Short: true,
			},
			{
				Title: "Recept",
				Value: tx.Hash().Hex(),
				Short: true,
			},
		},
		Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	slackParams.Attachments = []slack.Attachment{attachment}
	_, _, err = Service.Slack.PostMessage(Config.Slack.OpsChannel, "Token提现", slackParams)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, receipt)
}
