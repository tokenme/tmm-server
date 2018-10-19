package blowup

import (
	//"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	"sync"
	"time"
)

type Server struct {
	service         *common.Service
	config          common.Config
	eventsCh        chan Event
	channels        map[string]chan Event
	checkSessionCh  chan struct{}
	createSessionCh chan struct{}
	stopCh          chan struct{}
	exitCh          chan struct{}
	wg              *sync.WaitGroup
	sync.RWMutex
}

type Event struct {
	SessionId uint64          `json:"session_id"`
	Rate      decimal.Decimal `json:"rate"`
	Value     decimal.Decimal `json:"value,omitempty"`
}

type Bid struct {
	SessionId uint64
	Value     decimal.Decimal
}

func NewServer(service *common.Service, config common.Config) *Server {
	return &Server{
		service:         service,
		config:          config,
		eventsCh:        make(chan Event, 100),
		channels:        make(map[string]chan Event),
		createSessionCh: make(chan struct{}, 1),
		checkSessionCh:  make(chan struct{}, 1),
		stopCh:          make(chan struct{}, 1),
		exitCh:          make(chan struct{}, 1),
		wg:              &sync.WaitGroup{},
	}
}

func (this *Server) Start() {
	go func() {
		for {
			select {
			case <-this.checkSessionCh:
				go this.checkSession()
			case ev := <-this.eventsCh:
				this.RLock()
				channels := this.channels
				this.RUnlock()
				for _, ch := range channels {
					ch <- ev
				}
			case <-this.stopCh:
				close(this.checkSessionCh)
				log.Warn("Blowup Stopping")
				this.Lock()
				for uid, channel := range this.channels {
					close(channel)
					delete(this.channels, uid)
				}
				this.Unlock()
				log.Warn("Blowup Channel Closed")
				this.exitCh <- struct{}{}
				return
			}
		}
	}()
	this.checkSessionCh <- struct{}{}
}

func (this *Server) Stop() {
	this.wg.Wait()
	this.stopCh <- struct{}{}
	<-this.exitCh
	return
}

func (this *Server) NewChannel(uid string) {
	this.Lock()
	defer this.Unlock()
	if _, found := this.channels[uid]; !found {
		this.channels[uid] = make(chan Event, 100)
	}
}

func (this *Server) CloseChannel(uid string) {
	this.Lock()
	defer this.Unlock()
	if _, found := this.channels[uid]; found {
		close(this.channels[uid])
		delete(this.channels, uid)
	}
}

func (this *Server) Channel(uid string) chan Event {
	this.RLock()
	if channel, found := this.channels[uid]; found {
		this.RUnlock()
		return channel
	}
	this.RUnlock()
	return nil
}

func (this *Server) checkSession() error {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT id, NOW()-created_at FROM tmm.blowup_sessions WHERE finish_at>NOW() LIMIT 1`)
	if err != nil {
		log.Error(err.Error())
		go func() {
			time.Sleep(1 * time.Second)
			this.checkSessionCh <- struct{}{}
		}()
		return err
	}
	if len(rows) == 0 {
		_, _, err := db.Query(`DELETE FROM tmm.blowup_sessions`)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		this.eventsCh <- Event{}
		go func() {
			time.Sleep(20 * time.Second)
			this.createSession()
			this.checkSessionCh <- struct{}{}
		}()
		return nil
	}
	row := rows[0]
	rate := row.Int64(1)
	this.eventsCh <- Event{
		SessionId: row.Uint64(0),
		Rate:      decimal.New(rate, -2),
	}
	go func() {
		time.Sleep(1 * time.Second)
		this.checkSessionCh <- struct{}{}
	}()
	return nil
}

func (this *Server) createSession() error {
	db := this.service.Db
	_, ret, err := db.Query(`INSERT INTO tmm.blowup_sessions (created_at, finish_at) VALUES (NOW(), FROM_UNIXTIME((UNIX_TIMESTAMP(NOW())+CEIL(RAND()*100) + 30)))`)
	if err != nil {
		log.Error(err.Error())
		go func() {
			time.Sleep(1 * time.Second)
			this.checkSessionCh <- struct{}{}
		}()
		return err
	}
	this.eventsCh <- Event{
		SessionId: ret.InsertId(),
	}
	return nil
}
