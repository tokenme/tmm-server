package transferwatcher

import (
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"fmt"
	"github.com/FrontMage/xinge"
	"github.com/FrontMage/xinge/auth"
	xgreq "github.com/FrontMage/xinge/req"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	"io/ioutil"
	"net/http"
	"strings"
)

type Watcher struct {
	token     *eth.Token
	decimals  int
	service   *common.Service
	config    common.Config
	sink      chan *eth.TokenTransfer
	exitCh    chan struct{}
	canExitCh chan struct{}
}

func NewWatcher(tokenAddress string, service *common.Service, config common.Config) (*Watcher, error) {
	token, err := utils.NewToken(tokenAddress, service.GethWSS)
	if err != nil {
		return nil, err
	}
	decimals, err := utils.TokenDecimal(token)
	if err != nil {
		return nil, err
	}
	return &Watcher{
		token:     token,
		decimals:  decimals,
		service:   service,
		config:    config,
		sink:      make(chan *eth.TokenTransfer, 100),
		exitCh:    make(chan struct{}, 1),
		canExitCh: make(chan struct{}, 1),
	}, nil
}

func (this *Watcher) Start() error {
	go func() {
		select {
		case transfer := <-this.sink:
			this.handleTransfer(transfer)
		case <-this.exitCh:
			this.canExitCh <- struct{}{}
			return
		}
	}()
	_, err := this.token.WatchTransfer(nil, this.sink, nil, nil)
	if err == nil {
		log.Info("Transfer Watcher started!")
	}
	return err
}

func (this *Watcher) Stop() {
	this.exitCh <- struct{}{}
	<-this.canExitCh
	log.Info("Transfer Watcher stopped!")
	return
}

func (this *Watcher) handleTransfer(ev *eth.TokenTransfer) {
	tx := ev.Raw.TxHash.Hex()
	db := this.service.Db
	_, ret, err := db.Query(`UPDATE tmm.exchange_records AS er SET er.status=1 WHERE er.tx='%s'`, db.Escape(tx))
	if err != nil {
		log.Error(err.Error())
	}
	continueUpdate := true
	if ret != nil && ret.AffectedRows() > 0 {
		continueUpdate = false
	}
	if continueUpdate {
		_, ret, err := db.Query(`UPDATE tmm.withdraw_txs SET tx_status=1 WHERE tx='%s'`, db.Escape(tx))
		if ret != nil && ret.AffectedRows() > 0 {
			log.Error(err.Error())
			continueUpdate = false
		}
	}
	if continueUpdate {
		_, ret, err = db.Query(`UPDATE tmm.orderbooks SET deposit_tx_status=1 WHERE deposit_tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
		if ret != nil && ret.AffectedRows() > 0 {
			continueUpdate = false
		}
	}
	if continueUpdate {
		_, ret, err = db.Query(`UPDATE tmm.orderbooks SET withdraw_tx_status=1 WHERE withdraw_tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
		if ret != nil && ret.AffectedRows() > 0 {
			continueUpdate = false
		}
	}
	if continueUpdate {
		_, _, err = db.Query(`UPDATE tmm.orderbook_trades SET tx_status=1 WHERE tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
		if ret != nil && ret.AffectedRows() > 0 {
			continueUpdate = false
		}
	}
	this.push(ev)
}

func (this *Watcher) push(ev *eth.TokenTransfer) {
	value := decimal.NewFromBigInt(ev.Value, int32(-1*this.decimals))
	from := strings.ToLower(ev.From.Hex())
	to := strings.ToLower(ev.To.Hex())
	db := this.service.Db
	rows, _, err := db.Query(`SELECT d.user_id, d.push_token, u.wallet_addr, d.language, d.platform FROM tmm.devices AS d INNER JOIN ucoin.users AS u ON (u.id=d.user_id) WHERE u.wallet_addr IN ('%s', '%s')`, db.Escape(from), db.Escape(to))
	if err != nil {
		log.Error(err.Error())
		return
	}
	if len(rows) == 0 {
		return
	}
	for _, row := range rows {
		language := "en"
		if strings.Contains(row.Str(3), "zh") {
			language = "zh"
		}
		deviceToken := row.Str(1)
		platform := row.Str(4)
		var (
			title   string
			content string
		)
		if row.Str(2) == from {
			switch language {
			case "en":
				title = "UCoin transfer out notify"
				content = fmt.Sprintf("You have %s UC sent out.", value.String())
			case "zh":
				title = "UCoin 转出提醒"
				content = fmt.Sprintf("您的账户已转出 %s UC", value.String())
			}
		} else {
			switch language {
			case "en":
				title = "UCoin transfer in notify"
				content = fmt.Sprintf("You received %s UC.", value.String())
			case "zh":
				title = "UCoin 接收提醒"
				content = fmt.Sprintf("您的账户已接收 %s UC", value.String())
			}
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
		log.Info("%+v", r)
	}
}
