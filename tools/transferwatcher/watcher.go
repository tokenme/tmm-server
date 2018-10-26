package transferwatcher

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/mkideal/log"
	//"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
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
	if ret.AffectedRows() > 0 {
		continueUpdate = false
	}
	if continueUpdate {
		_, ret, err = db.Query(`UPDATE tmm.orderbooks SET deposit_status=1 WHERE deposit_tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
		if ret.AffectedRows() > 0 {
			continueUpdate = false
		}
	}
	if continueUpdate {
		_, ret, err = db.Query(`UPDATE tmm.orderbooks SET withdraw_status=1 WHERE withdraw_tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
		if ret.AffectedRows() > 0 {
			continueUpdate = false
		}
	}
	if continueUpdate {
		_, _, err = db.Query(`UPDATE tmm.orderbook_trades SET tx_status=1 WHERE tx='%s'`, db.Escape(tx))
		if err != nil {
			log.Error(err.Error())
		}
	}
}
