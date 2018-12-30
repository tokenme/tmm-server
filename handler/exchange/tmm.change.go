package exchange

import (
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"fmt"
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
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type TMMChangeRequest struct {
	Points    decimal.Decimal             `json:"points" form:"points" binding:"required"`
	DeviceId  string                      `json:"device_id" form:"device_id" binding:"required"`
	Direction common.TMMExchangeDirection `json:"direction" form:"direction" binding:"required"`
}

const TMMChangeRateKey = "TMMChangeRate-%d-%d"

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

	if CheckErr(user.IsBlocked(Service), c) {
		log.Error("Blocked User:%d", user.Id)
		return
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	changeRateKey := fmt.Sprintf(TMMChangeRateKey, req.Direction, user.Id)
	changeTime, err := redis.String(redisConn.Do("GET", changeRateKey))
	if CheckWithCode(err == nil, TOKEN_CHANGE_RATE_LIMIT_ERROR, "每次兑换时间间隔不能少于6小时", c) {
		log.Warn("TokenChangeRateLimit: %d, direction: %d, time: %s", user.Id, req.Direction, changeTime)
		return
	}

	exchangeRate, pointsPerTs, err := common.GetExchangeRate(Config, Service)
	if CheckErr(err, c) {
		return
	}
	/*
		minChangePointRate := decimal.New(3, -1)
		if Check(req.Direction == common.TMMExchangeOut && exchangeRate.Rate.LessThan(minChangePointRate), "Service not available", c) {
			return
			exchangeRate.Rate = minChangePointRate
		}
	*/
	if req.Points.LessThan(exchangeRate.MinPoints) {
		c.JSON(INVALID_MIN_POINTS_ERROR, exchangeRate)
		return
	}

	tmm := req.Points.Mul(exchangeRate.Rate)
	consumedTs := req.Points.Div(pointsPerTs)

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

	var (
		gasPrice    *big.Int
		fromAddress string
		toAddress   string
	)

	db := Service.Db
	if req.Direction == common.TMMExchangeIn {
		query := `SELECT
    d.points
FROM tmm.devices AS d
WHERE d.id='%s' AND d.user_id=%d`
		rows, _, err := db.Query(query, db.Escape(req.DeviceId), user.Id)
		if CheckErr(err, c) {
			return
		}
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
		tokenBalance, err := utils.TokenBalanceOf(token, poolPubKey)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(amount.Cmp(tokenBalance) == 1, NOT_ENOUGH_TOKEN_IN_POOL_ERROR, "no enough token in pool, please wait", c) {
			return
		}
		fromAddress = poolPubKey
		toAddress = user.Wallet
	} else {
		tokenBalance, err := utils.TokenBalanceOf(token, user.Wallet)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(amount.Cmp(tokenBalance) == 1, NOT_ENOUGH_TOKEN_ERROR, "not enough token", c) {
			return
		}
		fromAddress = user.Wallet
		toAddress = agentPubKey
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if CheckErr(err, c) {
		return
	}
	gasPrice, err = Service.Geth.SuggestGasPrice(c)
	if err == nil && gasPrice.Cmp(eth.MinGas) == -1 {
		gasPrice = eth.MinGas
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.TransferProxy(token, transactor, fromAddress, toAddress, amount)
	if CheckErr(err, c) {
		return
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	receipt := common.Transaction{
		Receipt:    tx.Hash().Hex(),
		Value:      tmm,
		Status:     2,
		InsertedAt: time.Now().Format(time.RFC3339),
	}

	_, _, err = db.Query(`INSERT INTO tmm.exchange_records (tx, status, user_id, device_id, tmm, points, direction) VALUES ('%s', %d, %d, '%s', '%s', '%s', %d)`, db.Escape(receipt.Receipt), receipt.Status, user.Id, db.Escape(req.DeviceId), receipt.Value.String(), req.Points.String(), req.Direction)
	if CheckErr(err, c) {
		return
	}
	var slackTitle string
	if req.Direction == common.TMMExchangeIn {
		slackTitle = "Points -> Token"
		_, _, err = db.Query(`UPDATE tmm.devices AS d SET d.points = IF(d.points - %s <= 0, 0, d.points - %s), d.consumed_ts = d.consumed_ts + %d WHERE d.id='%s' AND d.user_id=%d`, req.Points.String(), req.Points.String(), consumedTs.IntPart(), db.Escape(req.DeviceId), user.Id)
	} else {
		slackTitle = "Token -> Points"
		_, _, err = db.Query(`UPDATE tmm.devices AS d SET d.points = d.points + %s, d.total_ts = d.total_ts + %d WHERE d.id='%s' AND d.user_id=%d`, req.Points.String(), consumedTs.IntPart(), db.Escape(req.DeviceId), user.Id)
	}
	if CheckErr(err, c) {
		return
	}

	_, err = redisConn.Do("SETEX", changeRateKey, 60*60*6, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error(err.Error())
	}

	slackParams := slack.PostMessageParameters{Parse: "full", UnfurlMedia: true, Markdown: true}
	attachment := slack.Attachment{
		Color:      "success",
		AuthorName: user.ShowName,
		AuthorIcon: user.Avatar,
		Title:      slackTitle,
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
				Title: "Points",
				Value: req.Points.StringFixed(4),
				Short: true,
			},
			{
				Title: "Tokens",
				Value: tmm.StringFixed(4),
				Short: true,
			},
			{
				Title: "Gas Price",
				Value: gasPrice.String(),
				Short: true,
			},
			{
				Title: "Recept",
				Value: receipt.Receipt,
				Short: true,
			},
		},
		Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	slackParams.Attachments = []slack.Attachment{attachment}
	_, _, err = Service.Slack.PostMessage(Config.Slack.OpsChannel, "Point <-> Token", slackParams)
	if err != nil {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, receipt)
}
