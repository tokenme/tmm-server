package obs

import (
	"context"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/orderbook"
	commonutils "github.com/tokenme/tmm/utils"
	"math/big"
	"strings"
	"sync"
	"time"
)

var etherDecimals = decimal.New(params.Ether, 0)

type TxMap struct {
	Id uint64
	Tx string
}

type Server struct {
	config           common.Config
	service          *common.Service
	tokenDecimals    decimal.Decimal
	escrow           *eth.Escrow
	agentPrivKey     string
	agentPubKey      string
	globalLock       *sync.Mutex
	checkDepositCh   chan struct{}
	checkOrderbookCh chan struct{}
	checkDealCh      chan struct{}
	checkWithdrawCh  chan struct{}
	exitCh           chan struct{}
	canStopCh        chan struct{}
}

func NewServer(service *common.Service, config common.Config, globalLock *sync.Mutex) (*Server, error) {
	token, err := utils.NewToken(config.TMMTokenAddress, service.Geth)
	if err != nil {
		return nil, err
	}
	decimals, err := utils.TokenDecimal(token)
	if err != nil {
		return nil, err
	}
	escrow, err := utils.NewEscrow(config.TMMEscrowAddress, service.Geth)
	if err != nil {
		return nil, err
	}
	agentPrivKey, err := commonutils.AddressDecrypt(config.TMMAgentWallet.Data, config.TMMAgentWallet.Salt, config.TMMAgentWallet.Key)
	if err != nil {
		return nil, err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return nil, err
	}
	return &Server{
		config:           config,
		service:          service,
		tokenDecimals:    decimal.New(1, int32(decimals)),
		agentPrivKey:     agentPrivKey,
		agentPubKey:      agentPubKey,
		escrow:           escrow,
		globalLock:       globalLock,
		checkDepositCh:   make(chan struct{}, 1),
		checkOrderbookCh: make(chan struct{}, 1),
		checkDealCh:      make(chan struct{}, 1),
		checkWithdrawCh:  make(chan struct{}, 1),
		exitCh:           make(chan struct{}, 1),
		canStopCh:        make(chan struct{}, 1),
	}, nil
}

func (this *Server) Start() {
	shouldStop := false
	ctx, cancel := context.WithCancel(context.Background())
	go this.checkDeposits(ctx)
	go this.checkOrderbook(ctx)
	go this.checkDeals(ctx)
	go this.checkWithdraws(ctx)
	for !shouldStop {
		select {
		case <-this.checkDepositCh:
			go this.checkDeposits(ctx)
		case <-this.checkOrderbookCh:
			go this.checkOrderbook(ctx)
		case <-this.checkDealCh:
			go this.checkDeals(ctx)
		case <-this.checkWithdrawCh:
			go this.checkWithdraws(ctx)
		case <-this.exitCh:
			shouldStop = true
			cancel()
			this.canStopCh <- struct{}{}
			break
		}
	}
}

func (this *Server) Stop() {
	this.exitCh <- struct{}{}
	<-this.canStopCh
}

func (this *Server) checkOrderbook(ctx context.Context) {
	db := this.service.Db
	query := `SELECT ob.id, ob.side, ob.process_type, ob.quantity - ob.deal_quantity, ob.price, u.wallet_addr FROM tmm.orderbooks AS ob INNER JOIN ucoin.users AS u ON (u.id=ob.user_id) WHERE ob.deposit_tx_status=1 AND ob.online_status=0 AND ob.id>%d ORDER BY ob.id ASC LIMIT 1000`
	var (
		startId uint64
		endId   uint64
		orders  []*orderbook.Order
		engine  = orderbook.NewOrderBook()
	)
	for {
		rows, _, err := db.Query(query, startId)
		if err != nil {
			break
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			quantity, err := decimal.NewFromString(row.Str(3))
			if err != nil {
				log.Error(err.Error())
				continue
			}
			price, err := decimal.NewFromString(row.Str(4))
			if err != nil {
				log.Error(err.Error())
				continue
			}
			order := &orderbook.Order{
				TradeId:     row.Uint64(0),
				Side:        orderbook.Side(row.Uint(1)),
				ProcessType: orderbook.ProcessType(row.Uint(2)),
				Quantity:    quantity,
				Price:       price,
				Wallet:      row.Str(5),
			}
			orders = append(orders, order)
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	for _, order := range orders {
		trades, _ := engine.ProcessOrder(order, true)
		this.processTrades(ctx, trades)
	}
	time.Sleep(10 * time.Second)
	this.checkOrderbookCh <- struct{}{}
}

func (this *Server) processTrades(ctx context.Context, trades []*orderbook.Trade) (err error) {
	var (
		val          []string
		tradeVal     []string
		buyers       []ethcommon.Address
		sellers      []ethcommon.Address
		weiAmounts   []*big.Int
		tokenAmounts []*big.Int
		txHash       string
	)
	for _, trade := range trades {
		dealEth := trade.Quantity.Mul(trade.Price)
		if trade.Side == orderbook.Ask {
			sellers = append(sellers, ethcommon.HexToAddress(trade.Wallet))
			buyers = append(buyers, ethcommon.HexToAddress(trade.CounterPartyWallet))
		} else {
			buyers = append(buyers, ethcommon.HexToAddress(trade.Wallet))
			sellers = append(sellers, ethcommon.HexToAddress(trade.CounterPartyWallet))
		}
		wei := dealEth.Mul(etherDecimals)
		tokens := trade.Quantity.Mul(this.tokenDecimals)
		weiAmount, ok := new(big.Int).SetString(wei.Floor().String(), 10)
		if !ok {
			return fmt.Errorf("Internal Error: %s", wei.Floor().String())
		}
		tokenAmount, ok := new(big.Int).SetString(tokens.Floor().String(), 10)
		if !ok {
			return fmt.Errorf("Internal Error: %s", tokens.Floor().String())
		}
		weiAmounts = append(weiAmounts, weiAmount)
		tokenAmounts = append(tokenAmounts, tokenAmount)
		txHash, err = this.dealTransfer(ctx, buyers, sellers, weiAmounts, tokenAmounts)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		val = append(val, fmt.Sprintf("(%d, %s, %s)", trade.Id, trade.Quantity.String(), dealEth.String()))
		val = append(val, fmt.Sprintf("(%d, %s, %s)", trade.CounterParty, trade.Quantity.String(), dealEth.String()))
	}
	db := this.service.Db
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO tmm.orderbooks (id, deal_quantity, deal_eth) VALUES %s ON DUPLICATE KEY UPDATE deal_quantity=deal_quantity+VALUES(deal_quantity), deal_eth=deal_eth+VALUES(deal_eth), online_status=IF(quantity<=deal_quantity, 1, 0)`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
		}
	}
	for _, trade := range trades {
		tradeVal = append(tradeVal, fmt.Sprintf("(%d, %d, %d, %s, %s, '%s')", trade.Side, trade.Id, trade.CounterParty, trade.Quantity.String(), trade.Price.String(), db.Escape(txHash)))
		if len(tradeVal) > 0 {
			_, _, err := db.Query(`INSERT INTO tmm.orderbook_trades (side, trade_id, counter_party, quantity, price, tx) VALUES %s`, strings.Join(tradeVal, ","))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
	return nil
}

func (this *Server) dealTransfer(ctx context.Context, buyers []ethcommon.Address, sellers []ethcommon.Address, weiAmounts []*big.Int, tokenAmounts []*big.Int) (txHash string, err error) {
	gasPrice := new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Shannon))
	var gasLimit uint64 = 540000
	transactor := eth.TransactorAccount(this.agentPrivKey)
	nonce, err := eth.Nonce(ctx, this.service.Geth, this.service.Redis.Master, this.globalLock, this.agentPubKey, this.config.Geth)
	if err != nil {
		return "", err
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}
	eth.TransactorUpdate(transactor, transactorOpts, ctx)
	tx, err := this.escrow.BatchDeal(transactor, buyers, sellers, weiAmounts, tokenAmounts)
	if err != nil {
		return "", err
	}
	err = eth.NonceIncr(ctx, this.service.Geth, this.service.Redis.Master, this.globalLock, this.agentPubKey, this.config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	txHash = tx.Hash().Hex()
	return txHash, nil
}

func (this *Server) checkDeals(ctx context.Context) {
	db := this.service.Db
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(100, func(tx interface{}) error {
		defer wg.Done()
		txHex := tx.(string)
		log.Info("Checking Deal Receipt:%s", txHex)
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, txHex)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		if receipt == nil {
			return nil
		}
		db := this.service.Db
		_, _, err = db.Query(`UPDATE tmm.orderbook_trades SET tx_status=%d WHERE tx='%s'`, receipt.Status, txHex)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	})

	query := `SELECT id, tx FROM tmm.orderbook_trades WHERE tx_status=2 AND id>%d ORDER BY id ASC LIMIT 1000`
	var (
		startId uint64
		endId   uint64
	)
	for {
		rows, _, err := db.Query(query, startId)
		if err != nil {
			break
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			endId = row.Uint64(0)
			tx := row.Str(1)
			wg.Add(1)
			pool.Serve(tx)
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
	this.checkDealCh <- struct{}{}
}

func (this *Server) checkDeposits(ctx context.Context) {
	db := this.service.Db
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(100, func(tx interface{}) error {
		defer wg.Done()
		msg := tx.(TxMap)
		txHex := msg.Tx
		log.Info("Checking Deposit Receipt:%s", txHex)
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, txHex)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		if receipt == nil {
			return nil
		}
		db := this.service.Db
		_, _, err = db.Query(`UPDATE tmm.orderbooks SET deposit_tx_status=%d WHERE id=%d`, receipt.Status, msg.Id)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	})

	query := `SELECT id, deposit_tx FROM tmm.orderbooks WHERE deposit_tx_status=2 AND id>%d ORDER BY id ASC LIMIT 1000`
	var (
		startId uint64
		endId   uint64
	)
	for {
		rows, _, err := db.Query(query, startId)
		if err != nil {
			break
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			endId = row.Uint64(0)
			tx := row.Str(1)
			wg.Add(1)
			pool.Serve(TxMap{Id: endId, Tx: tx})
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
	this.checkDepositCh <- struct{}{}
}

func (this *Server) checkWithdraws(ctx context.Context) {
	db := this.service.Db
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(100, func(tx interface{}) error {
		defer wg.Done()
		msg := tx.(TxMap)
		txHex := msg.Tx
		log.Info("Checking Deposit Receipt:%s", txHex)
		receipt, err := utils.TransactionReceipt(this.service.Geth, ctx, txHex)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		if receipt == nil {
			return nil
		}
		db := this.service.Db
		_, _, err = db.Query(`UPDATE tmm.orderbooks SET withdraw_tx_status=%d WHERE id=%d`, receipt.Status, msg.Id)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	})

	query := `SELECT id, withdraw_tx FROM tmm.orderbooks WHERE withdraw_tx_status=2 AND id>%d ORDER BY id ASC LIMIT 1000`
	var (
		startId uint64
		endId   uint64
	)
	for {
		rows, _, err := db.Query(query, startId)
		if err != nil {
			break
		}
		if len(rows) == 0 {
			break
		}
		for _, row := range rows {
			endId = row.Uint64(0)
			tx := row.Str(1)
			wg.Add(1)
			pool.Serve(TxMap{Id: endId, Tx: tx})
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
	this.checkWithdrawCh <- struct{}{}
}
