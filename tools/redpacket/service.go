package redpacket

import (
	"context"
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
	commonutils "github.com/tokenme/tmm/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Service struct {
	service         *common.Service
	config          common.Config
	globalLock      *sync.Mutex
	checkTxCh       chan struct{}
	transferBonusCh chan struct{}
	exitCh          chan struct{}
	canStopCh       chan struct{}
}

func NewService(service *common.Service, config common.Config, globalLock *sync.Mutex) *Service {
	return &Service{
		service:         service,
		config:          config,
		globalLock:      globalLock,
		checkTxCh:       make(chan struct{}, 1),
		transferBonusCh: make(chan struct{}, 1),
		exitCh:          make(chan struct{}, 1),
		canStopCh:       make(chan struct{}, 1),
	}
}

func (this *Service) Start() {
	shouldStop := false
	ctx, cancel := context.WithCancel(context.Background())
	go this.CheckTransferBonus(ctx)
	go this.CheckTx(ctx)
	for !shouldStop {
		select {
		case <-this.checkTxCh:
			go this.CheckTx(ctx)
		case <-this.transferBonusCh:
			go this.CheckTransferBonus(ctx)
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
	rows, _, err := db.Query(`SELECT tx FROM tmm.redpacket_recipients WHERE tx!='' AND tx_status=2 ORDER BY id ASC LIMIT 1000`)
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
		_, _, err = db.Query(`UPDATE tmm.redpacket_recipients SET tx_status=%d WHERE tx='%s' AND tx_status=2`, receipt.Status, txHex)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return nil
}

func (this *Service) CheckTransferBonus(ctx context.Context) error {
	defer func() {
		time.Sleep(10 * time.Second)
		this.transferBonusCh <- struct{}{}
	}()
	db := this.service.Db
	_, _, err := db.Query(`UPDATE tmm.redpacket_recipients AS rr,
(SELECT rr.id, wx.user_id, wx.union_id
FROM tmm.redpacket_recipients AS rr
INNER JOIN tmm.wx ON (wx.union_id=rr.union_id)
WHERE rr.user_id IS NULL) AS t
SET rr.user_id=t.user_id
WHERE rr.id=t.id`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	rows, _, err := db.Query(`SELECT rb.id, u.wallet_addr, rb.tmm, u.id FROM tmm.redpacket_recipients AS rb INNER JOIN ucoin.users AS u ON (u.id=rb.user_id) WHERE rb.tx='' ORDER BY rb.id ASC LIMIT 1000`)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	for _, row := range rows {
		id := row.Uint64(0)
		walletAddr := row.Str(1)
		tokenAmount, _ := decimal.NewFromString(row.Str(2))
		userId := row.Uint64(3)
		receipt, err := this.transferToken(walletAddr, tokenAmount, ctx)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		this.PushMsg(userId, tokenAmount)
		_, _, err = db.Query(`UPDATE tmm.redpacket_recipients SET tx='%s' WHERE id=%d`, receipt, id)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return nil
}

func (this *Service) transferToken(userWallet string, tokenAmount decimal.Decimal, c context.Context) (receipt string, err error) {
	token, err := utils.NewToken(this.config.TMMTokenAddress, this.service.Geth)
	if err != nil {
		return receipt, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return receipt, err
	}
	tmmInt := tokenAmount.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if !ok {
		return receipt, nil
	}

	agentPrivKey, err := commonutils.AddressDecrypt(this.config.TMMAgentWallet.Data, this.config.TMMAgentWallet.Salt, this.config.TMMAgentWallet.Key)
	if err != nil {
		return receipt, err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return receipt, err
	}

	tokenBalance, err := utils.TokenBalanceOf(token, agentPubKey)
	if err != nil {
		return receipt, err
	}
	if amount.Cmp(tokenBalance) == 1 {
		return receipt, nil
	}

	transactor := eth.TransactorAccount(agentPrivKey)
	this.globalLock.Lock()
	defer this.globalLock.Unlock()
	nonce, err := eth.Nonce(c, this.service.Geth, this.service.Redis.Master, agentPubKey, this.config.Geth)
	if err != nil {
		return receipt, err
	}
	var gasPrice *big.Int
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.Transfer(token, transactor, userWallet, amount)
	if err != nil {
		return receipt, err
	}
	err = eth.NonceIncr(c, this.service.Geth, this.service.Redis.Master, agentPubKey, this.config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	receipt = tx.Hash().Hex()
	return receipt, nil
}

func (this *Service) PushMsg(userId uint64, bonus decimal.Decimal) {
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
		title = "UCoin Redpacket"
		content = fmt.Sprintf("You just received %sUC from UCoin Redpacket", bonus.String())
	case "zh":
		title = "友币红包"
		content = fmt.Sprintf("您获得 %sUC 红包", bonus.String())
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
