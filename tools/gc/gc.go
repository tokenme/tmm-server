package gc

import (
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"time"
)

const (
	ActiveAppGCDays int = 2
)

type Handler struct {
	Service *common.Service
	Config  common.Config
	exitCh  chan struct{}
}

func New(service *common.Service, config common.Config) *Handler {
	return &Handler{
		Service: service,
		Config:  config,
		exitCh:  make(chan struct{}, 1),
	}
}

func (this *Handler) Start() {
	log.Info("GC Start")
	dailyTicker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-dailyTicker.C:
			this.activeAppRecycle()
			this.inviteSubmissionRecycle()
			this.expiredReadingLogs()
			this.expiredReadingLogKws()
		case <-this.exitCh:
			dailyTicker.Stop()
			return
		}
	}
}

func (this *Handler) Stop() {
	close(this.exitCh)
	log.Info("GC Stopped")
}

func (this *Handler) activeAppRecycle() error {
	db := this.Service.Db
	_, _, err := db.Query(`UPDATE tmm.apps SET is_active=0 WHERE lastping_at<DATE_SUB(NOW(), INTERVAL %d DAY)`, ActiveAppGCDays)
	return err
}

func (this *Handler) inviteSubmissionRecycle() error {
	db := this.Service.Db
	_, _, err := db.Query(`DELETE FROM tmm.invite_submissions WHERE inserted_at<DATE_SUB(NOW(), INTERVAL 1 DAY)`)
	return err
}

func (this *Handler) expiredReadingLogs() error {
	db := this.Service.Db
	_, _, err := db.Query(`DELETE FROM tmm.reading_logs WHERE inserted_at<DATE_SUB(NOW(), INTERVAL 30 DAY)`)
	return err
}

func (this *Handler) expiredReadingLogKws() error {
	db := this.Service.Db
	_, _, err := db.Query(`DELETE FROM tmm.user_reading_kws WHERE updated_at<DATE_SUB(NOW(), INTERVAL 7 DAY)`)
	return err
}
