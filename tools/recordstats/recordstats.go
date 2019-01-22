package recordstats

import (
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"time"
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
	log.Info("record-stats Start")
	dailyTicker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-dailyTicker.C:
			this.RecordExchange()
			this.RecordMakePoints()
			this.RecordWithDrawCash()
		case <-this.exitCh:
			dailyTicker.Stop()
			return
		}
	}
}

func (this *Handler) Stop() {
	close(this.exitCh)
	log.Info("record-stats Stopped")
}

func (this *Handler) RecordWithDrawCash() error {
	db := this.Service.Db

	_, _, err := db.Query(`INSERT INTO tmm.user_withdarw_logs(user_id,points,points_cny,record_on) 
SELECT 
	user_id AS user_id ,
	SUM(points) AS points ,
	SUM(cny) AS cny ,
	DATE(inserted_at) AS record_on
FROM 	
	tmm.point_withdraws 
WHERE
	verified = 1  AND inserted_at > '%s' AND inserted_at < '%s'
GROUP BY 
	user_id,record_on 
ON DUPLICATE KEY UPDATE points =VALUES(points),points_cny=VALUES(points_cny)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))

	_, _, err = db.Query(`INSERT INTO tmm.user_withdarw_logs(user_id,tmm,tmm_cny,record_on) 
SELECT 
	user_id AS user_id ,
	SUM(tmm) AS tmm ,
	SUM(cny) AS cny ,
	DATE(inserted_at) AS record_on
FROM 	
	tmm.withdraw_txs  
WHERE 
	tx_status = 1 AND verified != -1 AND inserted_at > '%s' AND inserted_at < '%s'
GROUP BY 
	user_id,record_on 
ON DUPLICATE KEY UPDATE tmm =VALUES(tmm),tmm_cny=VALUES(tmm_cny)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))

	return err
}

func (this *Handler) RecordMakePoints() error {
	db := this.Service.Db

	_, _, err := db.Query(`INSERT INTO tmm.user_points_logs(user_id,reading,record_on)
SELECT 
	user_id AS user_id,
	SUM(point) AS reading,
	DATE(inserted_at) record_on
FROM 
	tmm.reading_logs
WHERE inserted_at > '%s' AND inserted_at < '%s'
GROUP BY user_id,record_on
ON DUPLICATE KEY UPDATE reading =VALUES(reading)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))

	if err != nil {
		return err
	}

	_, _, err = db.Query(`INSERT INTO tmm.user_points_logs(user_id,_share,share_viewers,share_count,record_on)
SELECT
	dev.user_id,
	SUM(dst.points),
	SUM(dst.viewers),
	COUNT(DISTINCT dst.task_id),
	DATE(dst.inserted_at) record_on
FROM 
	tmm.device_share_tasks AS dst 
INNER JOIN tmm.devices AS dev ON(dev.id = dst.device_id)
WHERE dst.inserted_at > '%s' AND inserted_at < '%s'
GROUP BY user_id,record_on
ON DUPLICATE KEY UPDATE _share =VALUES(_share), share_viewers = VALUES(share_viewers), share_count = VALUES(share_count)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))

	if err != nil {
		return err
	}

	_, _, err = db.Query(`INSERT INTO tmm.user_points_logs(user_id,app,record_on)
SELECT
	dev.user_id,
	SUM(dat.points),
	DATE(dat.inserted_at) record_on
FROM 
	tmm.device_app_tasks AS dat 
INNER JOIN tmm.devices AS dev ON(dev.id = dat.device_id)
WHERE inserted_at > '%s' AND inserted_at < '%s' AND status = 1 
GROUP BY user_id,record_on
ON DUPLICATE KEY UPDATE app =VALUES(app)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))


	_,_,err = db.Query(`INSERT INTO tmm.user_points_logs(user_id,invite,bonus,invite_count,tmm_bonus,record_on)
SELECT
	ib.user_id AS user_id,
	SUM(IF(ib.task_type = 0 ,ib.bonus,0)) AS invite,
	SUM(IF(ib.task_type != 0 ,ib.bonus,0)) AS bouns,
	COUNT(DISTINCT IF(ib.user_id != ib.from_user_id AND ib.task_type = 0,ib.from_user_id,NULL)) AS invite_count,
	SUM(IF(ib.tx_status = 1,tmm,0)) AS tmm_bouns,
	DATE(ib.inserted_at) AS record_on
FROM 
	tmm.invite_bonus AS ib 
WHERE inserted_at > '%s' AND inserted_at < '%s' AND tx_status != 0 
GROUP BY user_id,record_on
ON DUPLICATE KEY UPDATE invite =VALUES(invite), bonus = VALUES(bonus) , invite_count = VALUES(invite_count) , 
tmm_bonus = VALUES(tmm_bonus)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))
	return err
}


func (this *Handler) RecordExchange() error {
	db := this.Service.Db

	_, _, err := db.Query(`INSERT INTO tmm.user_exchange_logs(user_id,record_on,get_uc,pay_points,pay_uc,get_points)
SELECT 
	user_id AS user_id ,
	DATE(inserted_at) AS record_on,
	SUM(IF(direction = 1, tmm,0)) AS get_uc,
	SUM(IF(direction = 1, points,0)) AS pay_points,
	SUM(IF(direction = -1, tmm,0)) AS pay_uc,
	SUM(IF(direction = -1, points,0)) AS get_points
FROM 
	tmm.exchange_records 
WHERE 
	inserted_at > '%s' AND inserted_at < '%s' AND status = 1
GROUP BY user_id,record_on
ON DUPLICATE KEY UPDATE 
pay_points = VALUES(pay_points),get_uc = VALUES(get_uc),
pay_uc = VALUES(pay_uc), get_points = VALUES(get_points)
`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`),time.Now().Format(`2006-01-02`))
	return err
}



