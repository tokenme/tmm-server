package redeem

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	//"github.com/ethereum/go-ethereum/params"
	"github.com/FrontMage/xinge"
	"github.com/FrontMage/xinge/auth"
	xgreq "github.com/FrontMage/xinge/req"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nlopes/slack"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/wechatpay"
	//"github.com/tokenme/tmm/tools/ethgasstation-api"
	"errors"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/tools/ykt"
	commonutils "github.com/tokenme/tmm/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PointsWithdrawRequest struct {
	DeviceId string          `json:"device_id" form:"device_id" binding:"required"`
	Points   decimal.Decimal `json:"points" form:"points" binding:"required"`
	Currency string          `json:"currency" form:"currency" binding:"required"`
}

func PointsWithdrawHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req PointsWithdrawRequest
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

	recyclePrice := common.GetPointPrice(Service, Config)
	minPointsRequire := decimal.New(int64(Config.MinPointsRedeem), 0)
	if req.Points.LessThan(minPointsRequire) {
		c.JSON(INVALID_MIN_POINTS_ERROR, common.ExchangeRate{MinPoints: minPointsRequire})
		return
	}

	if CheckWithCode(req.Points.LessThan(minPointsRequire) || req.Points.GreaterThan(decimal.New(20000, 0)), INVALID_POINTS_WITHDRAW_AMOUNT_ERROR, fmt.Sprintf("积分提现超出限制。最小%s 积分或累计超过20000积分", minPointsRequire.String()), c) {
		return
	}

	cash := req.Points.Mul(recyclePrice)

	rows, _, err = db.Query(`SELECT points FROM tmm.devices WHERE id='%s' AND user_id=%d`, db.Escape(req.DeviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	accountPoints, _ := decimal.NewFromString(rows[0].Str(0))
	if CheckWithCode(accountPoints.LessThan(req.Points), NOT_ENOUGH_POINTS_ERROR, "not enough points", c) {
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
		log.Error("cash: %s, points: %s, recyclePrice:%s", cny.String(), req.Points.String(), recyclePrice.String())
		return
	}

	exceeded, _, err := common.ExceededDailyWithdraw(Service, Config)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(exceeded, EXCEEDED_DAILY_WITHDRAW_LIMIT_ERROR, "超出系统今日提现额度，请明天再试", c) {
		return
	}

	tradeNumToken, err := uuid.NewV4()
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	tradeNum := commonutils.Md5(tradeNumToken.String())
	payClient := wechatpay.NewClient(Config.Wechat.AppId, Config.Wechat.MchId, Config.Wechat.Key, Config.Wechat.CertCrt, Config.Wechat.CertKey)
	payParams := &wechatpay.Request{
		TradeNum:    tradeNum,
		Amount:      cny.Mul(decimal.New(100, 0)).IntPart(),
		CallbackURL: fmt.Sprintf("%s/wechat/pay/callback", Config.BaseUrl),
		OpenId:      wxOpenId,
		Ip:          ClientIP(c),
		Desc:        "UCoin积分提现",
	}
	payParams.Nonce = commonutils.Md5(payParams.TradeNum)
	payRes, err := payClient.Pay(payParams)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	var errMsg = payRes.ErrCodeDesc
	if payRes.ErrCode != "" && strings.Contains(payRes.ErrCodeDesc, "该用户今日付款次数超过限制") {
		errMsg = "每个微信账号每天只能提现1次"
	}
	if Check(payRes.ErrCode != "", errMsg, c) {
		log.Error(payRes.ErrCodeDesc)
		return
	}
	var (
		consumedTs decimal.Decimal
		tmm        decimal.Decimal
	)
	exchangeRate, pointsPerTs, err := common.GetExchangeRate(Config, Service)
	if err == nil {
		consumedTs = req.Points.Div(pointsPerTs)
		tmm = req.Points.Mul(exchangeRate.Rate)
	}
	_, _, err = db.Query(`UPDATE tmm.devices SET points=points-%s, consumed_ts=consumed_ts+%d WHERE id='%s' AND user_id=%d AND points>= %s`, req.Points.String(), consumedTs.IntPart(), db.Escape(req.DeviceId), user.Id, req.Points.String())
	if err != nil {
		log.Error(err.Error())
	}
	_, _, err = db.Query(`INSERT INTO tmm.point_withdraws (trade_num, user_id, device_id, points, cny) VALUES ('%s', %d, '%s', %s, %s)`, db.Escape(tradeNum), user.Id, db.Escape(req.DeviceId), req.Points.String(), cny.String())
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}

	burnPool(c, tmm)

	pushMsg(user.Id, cny)
	slackParams := slack.PostMessageParameters{Parse: "full", UnfurlMedia: true, Markdown: true}
	attachment := slack.Attachment{
		Color:      "success",
		AuthorName: user.ShowName,
		AuthorIcon: user.Avatar,
		Title:      "积分提现",
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
				Title: "Points",
				Value: req.Points.StringFixed(9),
				Short: true,
			},
			{
				Title: "TradeNum",
				Value: tradeNum,
				Short: true,
			},
		},
		Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	slackParams.Attachments = []slack.Attachment{attachment}
	_, _, err = Service.Slack.PostMessage(Config.Slack.OpsChannel, "积分提现", slackParams)
	if err != nil {
		log.Error(err.Error())
		return
	}
	receipt := common.TMMWithdrawResponse{
		TMM:      req.Points,
		Cash:     cny,
		Currency: req.Currency,
	}
	c.JSON(http.StatusOK, receipt)
}

func burnPool(c *gin.Context, tmm decimal.Decimal) error {

	poolPrivKey, err := commonutils.AddressDecrypt(Config.TMMPoolWallet.Data, Config.TMMPoolWallet.Salt, Config.TMMPoolWallet.Key)
	if err != nil {
		return err
	}
	poolPubKey, err := eth.AddressFromHexPrivateKey(poolPrivKey)
	if err != nil {
		return err
	}

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if err != nil {
		return err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return err
	}

	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if err != nil {
		return err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return err
	}

	tmmInt := tmm.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if !ok {
		return errors.New("token big.Int conversion failed")
	}
	transactor := eth.TransactorAccount(agentPrivKey)
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		return err
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
	tx, err := utils.BurnFrom(token, transactor, poolPubKey, amount)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Burn %s in pool, tx: %s, because of points withdraw", tmm.String(), tx.Hash().Hex())
	return nil
}

func pushMsg(userId uint64, cny decimal.Decimal) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT d.push_token, d.language, d.platform FROM tmm.devices AS d WHERE d.user_id=%d ORDER BY lastping_at DESC LIMIT 1`, userId)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(rows) == 0 {
		return
	}
	row := rows[0]
	language := "en"
	if strings.Contains(row.Str(1), "zh") {
		language = "zh"
	}
	deviceToken := row.Str(0)
	platform := row.Str(2)
	var (
		title   string
		content string
	)
	switch language {
	case "en":
		title = "UCoin points withdraw notify"
		content = fmt.Sprintf("You just received ¥%s from UCoin.", cny.String())
	case "zh":
		title = "UCoin 积分提现提醒"
		content = fmt.Sprintf("您刚刚提现成功 ¥%s", cny.String())
	}
	var auther auth.Auther
	var pushReq *http.Request
	switch platform {
	case "ios":
		pushReq, _ = xgreq.NewPushReq(
			&xinge.Request{},
			xgreq.Platform(xinge.PlatformiOS),
			xgreq.EnvProd(),
			xgreq.AudienceType(xinge.AdToken),
			xgreq.MessageType(xinge.MsgTypeNotify),
			xgreq.TokenList([]string{deviceToken}),
			xgreq.PushID("0"),
			xgreq.Message(xinge.Message{
				Title:   title,
				Content: content,
			}),
		)
		auther = auth.Auther{AppID: Config.IOSXinge.AppId, SecretKey: Config.IOSXinge.SecretKey}
	case "android":
		pushReq, _ = xgreq.NewPushReq(
			&xinge.Request{},
			xgreq.Platform(xinge.PlatformAndroid),
			xgreq.EnvProd(),
			xgreq.AudienceType(xinge.AdToken),
			xgreq.MessageType(xinge.MsgTypeNotify),
			xgreq.TokenList([]string{deviceToken}),
			xgreq.PushID("0"),
			xgreq.Message(xinge.Message{
				Title:   title,
				Content: content,
			}),
		)
		auther = auth.Auther{AppID: Config.AndroidXinge.AppId, SecretKey: Config.AndroidXinge.SecretKey}
	}
	auther.Auth(pushReq)
	rsp, err := http.DefaultClient.Do(pushReq)
	if err != nil {
		log.Error(err.Error())
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	var r xinge.CommonRsp
	json.Unmarshal(body, r)
	log.Info("%+v", r)
}
