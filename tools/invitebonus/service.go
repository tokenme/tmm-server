package invitebonus

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
	checkBonusCh    chan struct{}
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
		checkBonusCh:    make(chan struct{}, 1),
		checkTxCh:       make(chan struct{}, 1),
		transferBonusCh: make(chan struct{}, 1),
		exitCh:          make(chan struct{}, 1),
		canStopCh:       make(chan struct{}, 1),
	}
}

func (this *Service) Start() {
	shouldStop := false
	ctx, cancel := context.WithCancel(context.Background())
	go this.CheckBonus(ctx)
	go this.CheckTransferBonus(ctx)
	go this.CheckTx(ctx)
	for !shouldStop {
		select {
		case <-this.checkBonusCh:
			go this.CheckBonus(ctx)
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
	rows, _, err := db.Query(`SELECT tmm_tx FROM tmm.invite_bonus WHERE tmm_tx!='' AND tx_status=2 ORDER BY id ASC LIMIT 1000`)
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
		_, _, err = db.Query(`UPDATE tmm.invite_bonus SET tx_status=%d WHERE tmm_tx='%s' AND tx_status=2`, receipt.Status, txHex)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
	return nil
}

func (this *Service) CheckBonus(ctx context.Context) error {
	defer func() {
		time.Sleep(6 * time.Hour)
		this.checkBonusCh <- struct{}{}
	}()
	if time.Now().Weekday() != time.Sunday {
		return nil
	}
	db := this.service.Db
	query := `INSERT INTO tmm.invite_bonus (user_id, from_user_id, tmm, task_type)
SELECT parent_id, parent_id, COUNT(DISTINCT user_id) * 50, 3
FROM
(
    SELECT parent_id, user_id, COUNT(DISTINCT record_on) AS days
    FROM (
        SELECT ic.parent_id AS parent_id, ic.user_id AS user_id, DATE(dst.inserted_at) AS record_on
        FROM tmm.invite_codes AS ic
        INNER JOIN tmm.devices AS d ON (d.user_id = ic.user_id)
        INNER JOIN tmm.device_share_tasks AS dst ON (dst.device_id=d.id AND dst.inserted_at >= DATE_SUB(NOW(), INTERVAL 7 DAY))
        WHERE ic.parent_id>0
        GROUP BY ic.parent_id, ic.user_id, record_on
        UNION ALL
        SELECT ic.parent_id AS parent_id, ic.user_id AS user_id, DATE(dat.inserted_at) AS record_on
        FROM tmm.invite_codes AS ic
        INNER JOIN tmm.devices AS d ON (d.user_id = ic.user_id)
        INNER JOIN tmm.device_app_tasks AS dat ON (dat.device_id=d.id AND dat.inserted_at >= DATE_SUB(NOW(), INTERVAL 7 DAY))
        WHERE ic.parent_id>0
        GROUP BY ic.parent_id, ic.user_id, record_on
    ) AS t GROUP BY parent_id, user_id
    HAVING days > 3
) AS t2 GROUP BY parent_id`
	_, _, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (this *Service) CheckTransferBonus(ctx context.Context) error {
	defer func() {
		time.Sleep(10 * time.Second)
		this.transferBonusCh <- struct{}{}
	}()
	db := this.service.Db
	rows, _, err := db.Query(`SELECT ib.id, u.wallet_addr, ib.tmm, u.id FROM tmm.invite_bonus AS ib INNER JOIN ucoin.users AS u ON (u.id=ib.user_id) WHERE ib.tmm_tx='' AND ib.task_type=3 ORDER BY ib.id ASC LIMIT 1000`)
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
		_, _, err = db.Query(`UPDATE tmm.invite_bonus SET tmm_tx='%s' WHERE id=%d`, receipt, id)
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

func (this *Service) FixBonus() {
	db := this.service.Db
	query := `SELECT DISTINCT ic.user_id, ic.parent_id, u.wallet_addr
FROM tmm.invite_codes AS ic
INNER JOIN ucoin.users AS u ON (u.id=ic.user_id)
WHERE
NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib WHERE ib.user_id=ic.user_id AND ib.user_id=ib.from_user_id AND ib.task_type=0 LIMIT 1)
AND NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib2 WHERE ib2.user_id=ic.parent_id AND ib2.from_user_id=ic.user_id AND ib2.task_type=0 LIMIT 1)
AND ic.parent_id>0
AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=ic.parent_id AND us.blocked=1 AND us.block_whitelist=0 LIMIT 1)
AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us2 WHERE us2.user_id=ic.user_id AND us2.blocked=1 AND us2.block_whitelist=0 LIMIT 1)
AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us3 WHERE us3.user_id=ic.grand_id AND us3.blocked=1 AND us3.block_whitelist=0 LIMIT 1)`
	ctx := context.Background()
	rows, _, err := db.Query(query)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, row := range rows {
		userId := row.Uint64(0)
		parentId := row.Uint64(1)
		userWallet := row.Str(2)
		log.Info("Transfer: %d", userId)
		{
			query := `SELECT d.id FROM tmm.devices AS d WHERE d.user_id=%d ORDER BY d.lastping_at DESC LIMIT 1`
			rows, _, err := db.Query(query, parentId)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			if len(rows) == 0 {
				continue
			}
			tokenAmount := decimal.New(543, 0)
			receipt, err := this.transferToken(userWallet, tokenAmount, ctx)
			if err != nil {
				log.Error("Bonus Transfer failed")
				continue
			}
			deviceId := rows[0].Str(0)
			_, _, err = db.Query(`UPDATE tmm.devices AS d SET d.points = d.points + 188 WHERE id='%s'`, db.Escape(deviceId))
			if err != nil {
				log.Error(err.Error())
				continue
			}
			_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, tmm, tmm_tx) VALUES (%d, %d, 0, %s, '%s'), (%d, %d, 188, 0, '')`, userId, userId, tokenAmount.String(), db.Escape(receipt), parentId, userId)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			log.Info("Transferred: %d, %s", userId, receipt)
		}
	}

	query = `SELECT ib.user_id, ib.tmm, u.wallet_addr
FROM tmm.invite_bonus AS ib
INNER JOIN ucoin.users AS u ON (u.id=ib.user_id)
WHERE
ib.user_id=ib.from_user_id
AND ib.tmm>0
AND ib.tmm_tx=''
AND ib.task_type=0
AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=ib.user_id AND us.blocked=1 AND us.block_whitelist=0 LIMIT 1)`
	rows, _, err = db.Query(query)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, row := range rows {
		userId := row.Uint64(0)
		tokenAmount, err := decimal.NewFromString(row.Str(1))
		if err != nil {
			log.Error(err.Error())
			continue
		}
		userWallet := row.Str(2)
		log.Info("Transfer: %d", userId)
		receipt, err := this.transferToken(userWallet, tokenAmount, ctx)
		if err != nil {
			log.Error("Bonus Transfer failed")
			continue
		}
		_, _, err = db.Query(`UPDATE tmm.invite_bonus SET tmm_tx='%s' WHERE user_id=from_user_id AND task_type=0 AND tmm>0 AND tmm_tx='' AND user_id=%d`, db.Escape(receipt), userId)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Info("Transferred: %d, %s", userId, receipt)
	}
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
		title = "UCoin active user bonus"
		content = fmt.Sprintf("You just received %sUC from UCoin because this user who is invited by you activated 3 days in 7 days", bonus.String())
	case "zh":
		title = "友币活跃用户奖励"
		content = fmt.Sprintf("由于您的下线7天内有3天活跃，您获得 %sUC 奖励", bonus.String())
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
