package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type UserActiveStats struct {
	NewUser        int           `json:"new_user"`
	LoginNumber    uint          `json:"login_number"`
	TodayActive    Active        `json:"today_active"`
	TwoActive      Active        `json:"two_active"`
	ThreeActive    Active        `json:"three_active"`
}

type Active struct {
	TotalActive int `json:"total_active"`
	DaySign     int `json:"day_sign"`
	Share       int `json:"share"`
	DownLoadApp int `json:"down_load_app"`
	Read        int `json:"read"`
}

func UserFunnelStatsHandler(c *gin.Context) {
	date := c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`))

	query := `
SELECT 
	us.id,
	IFNULL(dev.id,0)
FROM 
	ucoin.users AS us
LEFT JOIN tmm.devices  AS dev ON (dev.user_id = us.id )
WHERE created > '%s' AND created < DATE_ADD('%s', INTERVAL 1 DAY)
ORDER BY id DESC 
`
	var UserList []string
	var LoginList uint
	db := Service.Db
	rows, _, err := db.Query(query, date, date)
	if CheckErr(err, c) {
		return
	}

	for _, row := range rows {
		UserList = append(UserList, row.Str(0))
		if row.Str(1) != "0" {
			LoginList++
		}
	}

	var Stats UserActiveStats

	Stats.NewUser = len(UserList)
	Stats.LoginNumber = LoginList

	var todayActive, TwoActive, threeActive Active
	tm, _ := time.Parse(`2006-01-02`, date)
	if tm.Before(time.Now()) && len(UserList) > 0 {
		UserList, todayActive,err = GetActive(UserList, tm.Format(`2006-01-02`))
		if CheckErr(err,c){
			return
		}
	}
	if tm.AddDate(0, 0, 1).Before(time.Now()) && len(UserList) > 0 {
		tm = tm.AddDate(0, 0, 1)
		UserList, TwoActive,err = GetActive(UserList, tm.Format(`2006-01-02`))
		if CheckErr(err,c){
			return
		}
	}
	if tm.AddDate(0, 0, 1).Before(time.Now()) && len(UserList) > 0 {
		tm = tm.AddDate(0, 0, 1)
		_, threeActive,err = GetActive(UserList, tm.Format(`2006-01-02`))
		if CheckErr(err,c){
			return
		}
	}

	Stats.TodayActive = todayActive
	Stats.TwoActive = TwoActive
	Stats.ThreeActive = threeActive

	c.JSON(http.StatusOK, admin.Response{
		Data:    Stats,
		Message: admin.API_OK,
		Code:    0,
	})
}

func GetActive(idArray []string, when string) ([]string, Active,error) {
	db := Service.Db
	var active Active
	query := `
SELECT 
	u.id,
	IF(COUNT(reading.user_id) > 0,1,0) AS reading,	
	IF(COUNT(sha.task_id) > 0,1,0) AS _share,	
	IF(COUNT(app.task_id) > 0,1,0) AS _app,
	IF(COUNT(sign.days) > 0,1,0) AS daysign
FROM 
	ucoin.users AS u
LEFT JOIN tmm.devices AS dev ON  (dev.user_id = u.id)
LEFT JOIN (
	SELECT 
		user_id  AS user_id
	FROM 
		tmm.reading_logs 
	WHERE DATE(inserted_at) = '%s' OR DATE(updated_at) = '%s'
) AS reading ON reading.user_id = u.id
LEFT JOIN (
	SELECT 
		device_id , 
		task_id AS task_id
	FROM
		tmm.device_share_tasks 
	WHERE DATE(inserted_at) = '%s' 
	) AS sha ON (sha.device_id = dev.id )
LEFT JOIN (
	SELECT 
		device_id,
		task_id AS task_id
	FROM
		tmm.device_app_tasks 
	WHERE DATE(inserted_at) = '%s' 
	) AS app ON (app.device_id = dev.id )
LEFT JOIN (	
	SELECT 
		user_id AS user_id ,
		days AS days 
	FROM 
		tmm.daily_bonus_logs  
	WHERE	  
		DATE_SUB(updated_on,INTERVAL days-1 DAY ) <= '%s' AND '%s' <= updated_on
	) AS sign ON (sign.user_id = u.id)
WHERE u.id IN (%s) 
AND (reading.user_id > 0 OR  sha.task_id > 0 OR app.task_id > 0 OR sign.days > 0)
GROUP BY u.id
`
	rows, _, err := db.Query(query, db.Escape(when), db.Escape(when),
		db.Escape(when), db.Escape(when), db.Escape(when), db.Escape(when),
		strings.Join(idArray, ","))
	if err != nil {
		return nil, active,err
	}
	var newArray []string
	for _, row := range rows {
		active.TotalActive++
		newArray = append(newArray, row.Str(0))
		active.Read += row.Int(1)
		active.Share += row.Int(2)
		active.DownLoadApp += row.Int(3)
		active.DaySign += row.Int(4)
	}

	return newArray, active,nil
}
