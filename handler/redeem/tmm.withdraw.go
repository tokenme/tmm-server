package redeem

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	//"github.com/ethereum/go-ethereum/params"
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
	if CheckErr(user.IsBlocked(Service), c) {
		log.Error("Blocked User:%d", user.Id)
		return
	}
	db := Service.Db
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

	{
		query := `SELECT 1
FROM tmm.wx AS wx
INNER JOIN tmm.point_withdraws AS pw ON (pw.user_id=wx.user_id)
WHERE (wx.user_id=%d OR wx.union_id='%s') AND pw.inserted_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)
UNION
SELECT 1
FROM tmm.wx AS wx
INNER JOIN tmm.withdraw_txs AS wt ON (wt.user_id=wx.user_id)
WHERE (wx.user_id=%d OR wx.union_id='%s') AND wt.inserted_at >= DATE_SUB(NOW(), INTERVAL 1 DAY)`
		rows, _, err := db.Query(query, user.Id, wxUnionId, user.Id, wxUnionId)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) > 0, TOKEN_WITHDRAW_RATE_LIMIT_ERROR, "每个账号或微信号每次提现时间间隔不能少于24小时", c) {
			return
		}
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
	minCash := decimal.New(5, 0)
	maxCash := decimal.New(2000, 0)
	{
		rows, _, err := db.Query(`SELECT 1 FROM tmm.withdraw_txs WHERE user_id=%d UNION SELECT 1 FROM tmm.point_withdraws WHERE user_id=%d LIMIT 1`, user.Id, user.Id)
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			minCash = decimal.New(15, 0)
		}
	}

	if CheckWithCode(cny.LessThan(minCash) || cny.GreaterThan(maxCash), WECHAT_PAYMENT_ERROR, fmt.Sprintf("提现金额超出限制。最小金额%s元或累计超过%s元", minCash.String(), maxCash.String()), c) {
		log.Error("cash: %s, tmm: %s, recyclePrice:%s", cny.String(), req.TMM.String(), recyclePrice.String())
		return
	}

	exceeded, _, err := common.ExceededDailyWithdraw(Service, Config)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(exceeded, EXCEEDED_DAILY_WITHDRAW_LIMIT_ERROR, "超出系统今日提现额度，请明天再试", c) {
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
