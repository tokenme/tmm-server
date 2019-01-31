package tmmwithdraw

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FrontMage/xinge"
	"github.com/FrontMage/xinge/auth"
	xgreq "github.com/FrontMage/xinge/req"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/wechatpay"
	commonutils "github.com/tokenme/tmm/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	service                *common.Service
	config                 common.Config
	checkTxCh              chan struct{}
	checkExchangeRecordsCh chan struct{}
	checkWechatPayCh       chan struct{}
	exitCh                 chan struct{}
	canStopCh              chan struct{}
}

func NewService(service *common.Service, config common.Config) *Service {
	return &Service{
		service:                service,
		config:                 config,
		checkTxCh:              make(chan struct{}, 1),
		checkExchangeRecordsCh: make(chan struct{}, 1),
		checkWechatPayCh:       make(chan struct{}, 1),
		exitCh:                 make(chan struct{}, 1),
		canStopCh:              make(chan struct{}, 1),
	}
}

func (this *Service) Start() {
	shouldStop := false
	ctx, cancel := context.WithCancel(context.Background())
	go this.CheckTx(ctx)
	go this.WechatPay()
	go this.CheckExchangeRecords(ctx)
	for !shouldStop {
		select {
		case <-this.checkTxCh:
			go this.CheckTx(ctx)
		case <-this.checkWechatPayCh:
			go this.WechatPay()
		case <-this.checkExchangeRecordsCh:
			go this.CheckExchangeRecords(ctx)
		case <-this.exitCh:
			shouldStop = true
			cancel()
			this.canStopCh <- struct{}{}
			break
		}
	}
}

func (this *Service) Stop() {
	this.exitCh <- struct{}{}
	<-this.canStopCh
}

func (this *Service) CheckTx(ctx context.Context) error {
	defer func() {
		time.Sleep(10 * time.Second)
		this.checkTxCh <- struct{}{}
	}()
	db := this.service.Db
	rows, _, err := db.Query(`SELECT tx FROM tmm.withdraw_txs WHERE tx_status=2 ORDER BY inserted_at ASC LIMIT 1000`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(rows) == 0 {
		log.Warn("no withdraw tx")
		return nil
	}
	for _, row := range rows {
		txHex := row.Str(0)
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, txHex)
		if err != nil {
			log.Error("%s, %s", err.Error(), txHex)
			continue
		}
		if receipt == nil {
			continue
		}
		_, _, err = db.Query(`UPDATE tmm.withdraw_txs SET tx_status=%d WHERE tx='%s' AND tx_status=2`, receipt.Status, txHex)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return nil
}

func (this *Service) CheckExchangeRecords(ctx context.Context) error {
	defer func() {
		time.Sleep(12 * time.Second)
		this.checkExchangeRecordsCh <- struct{}{}
	}()
	db := this.service.Db
	rows, _, err := db.Query(`SELECT
        tx, status, device_id, tmm, points, direction, inserted_at
        FROM tmm.exchange_records
        WHERE status=2 ORDER BY inserted_at LIMIT 1000`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(rows) == 0 {
		log.Warn("no exchange records")
		return nil
	}
	pointsPerTs, err := common.GetPointsPerTs(this.service)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Checking %d exchange records", len(rows))
	for _, row := range rows {
		tmm, _ := decimal.NewFromString(row.Str(3))
		points, _ := decimal.NewFromString(row.Str(4))
		record := common.ExchangeRecord{
			Tx:         row.Str(0),
			Status:     common.ExchangeTxStatus(row.Uint(1)),
			DeviceId:   row.Str(2),
			Tmm:        tmm,
			Points:     points,
			Direction:  common.TMMExchangeDirection(row.Int(5)),
			InsertedAt: row.ForceLocaltime(6).Format(time.RFC3339),
		}
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, record.Tx)
		if err == nil {
			record.Status = common.ExchangeTxStatus(receipt.Status)
			if record.Status == common.ExchangeTxFailed && record.Direction == common.TMMExchangeIn {
				_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.exchange_records AS er SET d.points=d.points + er.points, d.total_ts = CEIL(d.total_ts + %s), er.status=0 WHERE d.id=er.device_id AND er.tx='%s'`, pointsPerTs.String(), db.Escape(record.Tx))
				if err != nil {
					log.Error(err.Error())
				}
				/*} else if record.Status == common.ExchangeTxSuccess && record.Direction == common.TMMExchangeOut {
				_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.exchange_records AS er SET d.points=d.points + er.points, d.total_ts = d.total_ts + %s, er.status=1 WHERE d.id=er.device_id AND er.tx='%s'`, pointsPerTs.String(), db.Escape(record.Tx))
				if err != nil {
					log.Error(err.Error())
				}*/
			} else {
				_, _, err := db.Query(`UPDATE tmm.exchange_records SET status=%d WHERE tx='%s'`, receipt.Status, db.Escape(record.Tx))
				if err != nil {
					log.Error(err.Error())
				}
			}
		} else {
			log.Error("%s, %s", err.Error(), record.Tx)
		}
	}
	return nil
}

func (this *Service) WechatPay() error {
	defer func() {
		time.Sleep(15 * time.Second)
		this.checkWechatPayCh <- struct{}{}
	}()
	db := this.service.Db
	rows, _, err := db.Query(`SELECT
        0, wt.tx, wt.cny, wt.client_ip, oi.open_id, wt.user_id
        FROM tmm.withdraw_txs AS wt
        INNER JOIN tmm.wx AS wx ON (wx.user_id=wt.user_id)
        INNER JOIN tmm.wx_openids AS oi ON (oi.union_id=wx.union_id AND oi.app_id='%s')
        WHERE wt.tx_status=1 AND wt.withdraw_status=2 AND wt.verified=1
        AND NOT EXISTS(SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=wx.user_id AND us.blocked=1 AND us.block_whitelist=0 LIMIT 1)
        UNION ALL
        SELECT pw.id, '', pw.cny, pw.client_ip, oi.open_id, pw.user_id
        FROM tmm.point_withdraws AS pw
        INNER JOIN tmm.wx AS wx ON (wx.user_id=pw.user_id)
        INNER JOIN tmm.wx_openids AS oi ON (oi.union_id=wx.union_id AND oi.app_id='%s')
        WHERE pw.trade_num='' AND pw.verified=1
        AND NOT EXISTS(SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=pw.user_id AND us.blocked=1 AND us.block_whitelist=0 LIMIT 1)
        LIMIT 1000`, db.Escape(this.config.Wechat.AppId), db.Escape(this.config.Wechat.AppId))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if len(rows) == 0 {
		log.Warn("Not payments")
		return nil
	}
	log.Info("WechatPay %d Accounts", len(rows))
	for _, row := range rows {
		txId := row.Uint64(0)
		txHex := row.Str(1)
		cny, err := decimal.NewFromString(row.Str(2))
		clientIp := row.Str(3)
		openId := row.Str(4)
		userId := row.Uint64(5)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		tradeNumToken, err := uuid.NewV4()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		paymentDesc := "UCoin提现"
		if txId > 0 {
			paymentDesc = "UCoin积分提现"
		}
		tradeNum := commonutils.Md5(tradeNumToken.String())
		payClient := wechatpay.NewClient(this.config.Wechat.AppId, this.config.Wechat.MchId, this.config.Wechat.Key, this.config.Wechat.CertCrt, this.config.Wechat.CertKey)
		payParams := &wechatpay.Request{
			TradeNum:    tradeNum,
			Amount:      cny.Mul(decimal.New(100, 0)).IntPart(),
			CallbackURL: fmt.Sprintf("%s/wechat/pay/callback", this.config.BaseUrl),
			OpenId:      openId,
			Ip:          clientIp,
			Desc:        paymentDesc,
		}
		payParams.Nonce = commonutils.Md5(payParams.TradeNum)
		payRes, err := payClient.Pay(payParams)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if payRes.ErrCode != "" {
			log.Error(payRes.ErrCodeDesc)
			continue
		}
		if txId == 0 {
			log.Info("Transferred %d CNY to Account:%d, UC WITHDRAW", payParams.Amount, userId)
			_, _, err = db.Query(`UPDATE tmm.withdraw_txs SET withdraw_status=1, trade_num='%s' WHERE tx='%s' AND withdraw_status=2`, db.Escape(tradeNum), db.Escape(txHex))
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			log.Info("Transferred %d CNY to Account:%d, POINT WITHDRAW", payParams.Amount, userId)
			_, _, err = db.Query(`UPDATE tmm.point_withdraws SET trade_num='%s' WHERE id=%d`, db.Escape(tradeNum), txId)
			if err != nil {
				log.Error(err.Error())
			}
		}
		this.PushMsg(userId, cny)
	}
	return nil
}

func (this *Service) PushMsg(userId uint64, cny decimal.Decimal) {
	db := this.service.Db
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
		title = "UCoin withdraw notify"
		content = fmt.Sprintf("You just received ¥%s from UCoin.", cny.String())
	case "zh":
		title = "UCoin 提现提醒"
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
		auther = auth.Auther{AppID: this.config.IOSXinge.AppId, SecretKey: this.config.IOSXinge.SecretKey}
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
		auther = auth.Auther{AppID: this.config.AndroidXinge.AppId, SecretKey: this.config.AndroidXinge.SecretKey}
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
	log.Info("Code:%d, %+v", rsp.StatusCode, r)
}
