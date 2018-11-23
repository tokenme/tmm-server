package tmmwithdraw

import (
	"context"
	"fmt"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/wechatpay"
	commonutils "github.com/tokenme/tmm/utils"
	"time"
)

type Service struct {
	service          *common.Service
	config           common.Config
	checkTxCh        chan struct{}
	checkWechatPayCh chan struct{}
	exitCh           chan struct{}
	canStopCh        chan struct{}
}

func NewService(service *common.Service, config common.Config) *Service {
	return &Service{
		service:          service,
		config:           config,
		checkTxCh:        make(chan struct{}, 1),
		checkWechatPayCh: make(chan struct{}, 1),
		exitCh:           make(chan struct{}, 1),
		canStopCh:        make(chan struct{}, 1),
	}
}

func (this *Service) Start() {
	shouldStop := false
	ctx, cancel := context.WithCancel(context.Background())
	go this.CheckTx(ctx)
	go this.WechatPay()
	for !shouldStop {
		select {
		case <-this.checkTxCh:
			go this.CheckTx(ctx)
		case <-this.checkWechatPayCh:
			go this.WechatPay()
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
	db := this.service.Db
	rows, _, err := db.Query(`SELECT tx FROM tmm.withdraw_txs WHERE tx_status=2 ORDER BY inserted_at ASC LIMIT 1000`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	for _, row := range rows {
		txHex := row.Str(0)
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, txHex)
		if err != nil {
			log.Error(err.Error())
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
	time.Sleep(10 * time.Second)
	this.checkTxCh <- struct{}{}
	return nil
}

func (this *Service) WechatPay() error {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT wt.tx, wt.cny, wt.client_ip, oi.open_id FROM tmm.withdraw_txs AS wt INNER JOIN tmm.wx AS wx ON (wx.user_id=wt.user_id) INNER JOIN tmm.wx_openids AS oi ON (oi.union_id=wx.union_id AND oi.app_id='%s') WHERE wt.tx_status=1 AND wt.withdraw_status=2 ORDER BY wt.inserted_at ASC LIMIT 1000`, db.Escape(this.config.Wechat.AppId))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	for _, row := range rows {
		txHex := row.Str(0)
		cny, err := decimal.NewFromString(row.Str(1))
		clientIp := row.Str(2)
		openId := row.Str(3)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		tradeNumToken, err := uuid.NewV4()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		tradeNum := commonutils.Md5(tradeNumToken.String())
		payClient := wechatpay.NewClient(this.config.Wechat.AppId, this.config.Wechat.MchId, this.config.Wechat.Key, this.config.Wechat.CertCrt, this.config.Wechat.CertKey)
		payParams := &wechatpay.Request{
			TradeNum:    tradeNum,
			Amount:      cny.Mul(decimal.New(100, 0)).IntPart(),
			CallbackURL: fmt.Sprintf("%s/wechat/pay/callback", this.config.BaseUrl),
			OpenId:      openId,
			Ip:          clientIp,
			Desc:        "UCoin提现",
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
		_, _, err = db.Query(`UPDATE tmm.withdraw_txs SET withdraw_status=1, trade_num='%s' WHERE tx='%s' AND withdraw_status=2`, db.Escape(tradeNum), db.Escape(txHex))
		if err != nil {
			log.Error(err.Error())
		}
	}
	time.Sleep(10 * time.Second)
	this.checkWechatPayCh <- struct{}{}
	return nil
}
