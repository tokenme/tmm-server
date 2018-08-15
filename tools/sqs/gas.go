package sqs

import (
	"context"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/go-msgqueue/msgqueue"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/utils"
	"math/big"
	"sync"
)

type GasMessage struct {
	Type   MsgType `codec:"type"`
	To     string  `codec:"to"`
	Amount uint64  `codec:"n"`
}

func (this GasMessage) String() string {
	return fmt.Sprintf("%d-%s-%d", this.Type, this.To, this.Amount)
}

type GasQueue struct {
	service   *common.Service
	config    common.Config
	queue     msgqueue.Queue
	processor *msgqueue.Processor
	txCh      chan GasMessage
	exitCh    chan struct{}
	canStopCh chan struct{}
}

func NewGasQueue(m msgqueue.Manager, service *common.Service, config common.Config) *GasQueue {
	queue := &GasQueue{
		service:   service,
		config:    config,
		txCh:      make(chan GasMessage, 1000),
		exitCh:    make(chan struct{}, 1),
		canStopCh: make(chan struct{}, 1),
	}
	opt := &msgqueue.Options{
		Name:    config.SQS.GasQueue,
		Handler: queue.Handler,
	}
	q := m.NewQueue(opt)
	queue.queue = q
	queue.processor = q.Processor()
	return queue
}

func (this *GasQueue) Start() {
	this.processor.Start()
	ctx := context.Background()
	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(10000, func(tx interface{}) error {
		defer wg.Done()
		msg := tx.(GasMessage)
		log.Info("Deposing Gas to:%s, amount:%d", msg.To, msg.Amount)
		outputWalletPrivateKey, err := utils.AddressDecrypt(this.config.OutputWallet.Data, this.config.OutputWallet.Salt, this.config.TokenSalt)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		transactor := eth.TransactorAccount(outputWalletPrivateKey)
		transactor.Value = new(big.Int).SetUint64(msg.Amount)
		_, err = eth.Transfer(transactor, this.service.Geth, ctx, msg.To)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		log.Info("Deposed Gas to:%s, amount:%d", msg.To, msg.Amount)
		return nil
	})

	shouldStop := false
	for !shouldStop {
		select {
		case tx := <-this.txCh:
			wg.Add(1)
			pool.Serve(tx)
		case <-this.exitCh:
			shouldStop = true
			break
		}
	}
	wg.Wait()
	this.canStopCh <- struct{}{}
}

func (this *GasQueue) Stop() {
	this.processor.Stop()
	this.exitCh <- struct{}{}
	<-this.canStopCh
}

func (this *GasQueue) Deposit(to string, amount uint64) error {
	return this.queue.Call(GasMessage{Type: GasMsg, To: to, Amount: amount})
}

func (this *GasQueue) Handler(msg GasMessage) error {
	switch msg.Type {
	case GasMsg:
		this.txCh <- msg
	}
	return nil
}
