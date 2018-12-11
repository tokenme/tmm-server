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
	minuteTicker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-dailyTicker.C:
			this.activeAppRecycle()
			this.inviteSubmissionRecycle()
			this.expiredReadingLogs()
			this.expiredReadingLogKws()
		case <-minuteTicker.C:
			this.expiredMobileCode()
			this.blockBadUsers()
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

func (this *Handler) expiredMobileCode() error {
	db := this.Service.Db
	_, _, err := db.Query(`DELETE FROM tmm.mobile_codes WHERE updated_at<DATE_SUB(NOW(), INTERVAL 10 MINUTE)`)
	return err
}

func (this *Handler) blockBadUsers() error {
	db := this.Service.Db
	query := `INSERT INTO tmm.user_settings (user_id, blocked)
SELECT user_id, 1 FROM (
SELECT
  ib.user_id AS user_id,
    COUNT( DISTINCT ib.from_user_id ) AS invites,
    SUM(IF(da.total_ts>0, 1, 0)) AS apps,
    SUM(ib.bonus) AS bonus,
    SUM(IFNULL(da.total_ts, 0)) AS ts
    FROM
        tmm.invite_bonus AS ib
        LEFT JOIN tmm.devices AS d ON (d.user_id=ib.from_user_id)
        LEFT JOIN tmm.device_apps AS da ON ( da.device_id = d.id )
WHERE ib.task_type=0 AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=ib.user_id AND us.blocked=1 LIMIT 1)
GROUP BY ib.user_id
HAVING invites>=10 AND bonus > 20000 AND apps<invites/2) AS t
ON DUPLICATE KEY UPDATE blocked=VALUES(blocked)`
	_, _, err := db.Query(query)
	query = `INSERT INTO tmm.user_settings (user_id, blocked) SELECT user_id, 1 FROM tmm.wx AS ws
WHERE EXISTS (
    SELECT
        1
    FROM tmm.wx AS wx
    INNER JOIN tmm.user_settings AS us ON (us.user_id=wx.user_id)
    WHERE us.blocked=1 AND wx.open_id=ws.open_id LIMIT 1
) AND NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us WHERE us.user_id=ws.user_id AND us.blocked=1 AND us.block_whitelist=0 LIMIT 1)
ON DUPLICATE KEY UPDATE blocked=VALUES(blocked)`
	_, _, err = db.Query(query)
	return err
}
