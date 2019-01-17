package info

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"strconv"
	"github.com/ziutek/mymysql/autorc"
	"strings"
	"fmt"
)

type FunnelData struct {
	LoginCount int                `json:"login_count"`
	UserList   []*admin.UserStats `json:"user_list,omitempty"`
}

func GetFunnelDataHandler(c *gin.Context) {

	date := c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`))
	days, err := strconv.Atoi(c.DefaultQuery(`days`, `1`))
	if CheckErr(err, c) {
		return
	}
	types, err := strconv.Atoi(c.DefaultQuery(`types`, `0`))
	if CheckErr(err, c) {
		return
	}
	tm, _ := time.Parse(`2006-01-02`, date)
	if Check(tm.AddDate(0, 0, days).After(time.Now()), admin.Not_Found, c) {
		return
	}

	query := `
SELECT 
	us.id,
	IFNULL(dev.id,0)
FROM 
	ucoin.users AS us
LEFT JOIN tmm.devices  AS dev ON (dev.user_id = us.id )
WHERE created > '%s' AND created < DATE_ADD('%s', INTERVAL 1 DAY)
`
	var UserList []string
	var LoginCount int
	db := Service.Db
	rows, _, err := db.Query(query, date, date)
	if CheckErr(err, c) {
		return
	}
	for _, row := range rows {
		UserList = append(UserList, row.Str(0))
		if row.Str(1) != "0" {
			LoginCount++
		}
	}

	stats := &FunnelData{}
	fieter := `AND (reading.user_id > 0 OR  sha.task_id > 0 OR app.task_id > 0 OR sign.days > 0)`

	for i := 0; i < days-1; i++ {
		UserList, stats, err = GetActiveData(db, UserList, tm.Format(`2006-01-02`), fieter)
		if CheckErr(err, c) {
			return
		}
		tm = tm.AddDate(0, 0, 1)
	}

	if types == 1 {
		fieter = `AND (IFNULL(reading.point,0)+IFNULL(SUM(sha.points),0) + IFNULL(SUM(app.task_id),0)+IFNULL(bonus.points,0)) = 1`
	}
	fmt.Println(UserList)
	_, stats, err = GetActiveData(db, UserList, tm.Format(`2006-01-02`), fieter)
	if CheckErr(err, c) {
		return
	}

	if days > 1 {
		stats.LoginCount = len(stats.UserList)
	} else {
		stats.LoginCount = LoginCount
	}
	c.JSON(http.StatusOK, admin.Response{
		Data:    stats,
		Message: admin.API_OK,
		Code:    0,
	})
}

func GetActiveData(db *autorc.Conn, idArray []string, date string, fieter string) ([]string, *FunnelData, error) {

	query := `
 SELECT 
	u.id AS id ,
	u.mobile AS mobile ,
	IFNULL(wx.nick,u.nickname) AS nick,
	u.created AS created ,
	(SELECT parent_id FROM tmm.invite_codes WHERE user_id = u.id ) AS parent_id,
	(SELECT wx.nick FROM tmm.invite_codes AS inv  LEFT JOIN tmm.wx AS wx ON wx.user_id = inv.parent_id WHERE inv.user_id = u.id ) AS parent_nick,
	IF(COUNT(sha.points)>0,1,0) AS _share,
	IF(COUNT(reading.point)>0,1,0) AS reading,
	IF(COUNT(app.task_id)>0,1,0) AS _app,
	IF(COUNT(sign.days) > 0,1,0) AS daysign,
	IFNULL(reading.point,0)+IFNULL(SUM(sha.points),0) + IFNULL(SUM(app.points),0)+IFNULL(bonus.points,0) AS total,
	IFNULL(SUM(wt.cny),0) + IFNULL(SUM(pt.cny),0) AS cny,
	sbl.second_count AS second_count,
	sbl.minute_count AS minute_count,
	IF(us.user_id > 0,IF(us.blocked = us.block_whitelist,0,1),0) AS blocked
FROM 
	ucoin.users AS u
LEFT JOIN tmm.wx AS wx ON (wx.nick = u.id)
LEFT JOIN tmm.devices AS dev ON  (dev.user_id = u.id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id = u.id)
LEFT JOIN tmm.share_blocked_logs AS sbl ON (sbl.user_id = u.id AND sbl.record_on = '%s')
LEFT JOIN tmm.withdraw_txs AS wt ON (wt.user_id = u.id AND wt.tx_status = 1 AND DATE(wt.inserted_at) = '%s' )
LEFT JOIN tmm.point_withdraws AS pt ON (pt.user_id = u.id AND DATE(pt.inserted_at) = '%s')
LEFT JOIN ( 
	SELECT 
		user_id  AS user_id,
		SUM(point) AS point 
	FROM 
		tmm.reading_logs 
	WHERE DATE(inserted_at) = '%s' OR DATE(updated_at) = '%s'
	GROUP BY user_id 
) AS reading ON reading.user_id = u.id
LEFT JOIN (
	SELECT 
		device_id , 
		SUM(points) AS points,
		task_id AS task_id 
	FROM
		tmm.device_share_tasks 
	WHERE DATE(inserted_at) = '%s' 
	GROUP BY device_id
	) AS sha ON (sha.device_id = dev.id )
LEFT JOIN (
	SELECT 
		device_id,
		SUM(IF(status = 1,points,0)) AS points,
		task_id AS task_id 	
	FROM
		tmm.device_app_tasks 
	WHERE DATE(inserted_at) = '%s'   AND  status = 1
	GROUP BY device_id 
	) AS app ON (app.device_id = dev.id )
LEFT JOIN (	 
	SELECT 
		user_id,
		days
	FROM 
		tmm.daily_bonus_logs  
	WHERE	  
		DATE_SUB(updated_on,INTERVAL days-1 DAY ) <= '%s' AND '%s' <= updated_on
	) AS sign ON (sign.user_id = u.id)
LEFT JOIN ( 
SELECT 
	SUM(bonus) AS points,
	user_id 
FROM 
	tmm.invite_bonus 
WHERE 	
	DATE(inserted_at) = '%s'  
GROUP BY user_id 
) AS bonus ON (bonus.user_id = u.id)
WHERE u.id IN (%s)  %s
GROUP BY u.id
`
	fmt.Printf(query, db.Escape(date), db.Escape(date),
		db.Escape(date), db.Escape(date), db.Escape(date), db.Escape(date),
		db.Escape(date), db.Escape(date), db.Escape(date), db.Escape(date),
		strings.Join(idArray, ","), fieter)
	rows, res, err := db.Query(query, db.Escape(date), db.Escape(date),
		db.Escape(date), db.Escape(date), db.Escape(date), db.Escape(date),
		db.Escape(date), db.Escape(date), db.Escape(date), db.Escape(date),
		strings.Join(idArray, ","), fieter)
	if err != nil {
		return nil, nil, err
	}

	data := &FunnelData{}
	newArray := []string{}
	for _, row := range rows {
		parent := &admin.User{}
		parent.Id = row.Uint64(res.Map(`parent_id`))
		parent.Nick = row.Str(res.Map(`parent_nick`))
		user := &admin.UserStats{}
		user.Id = row.Uint64(res.Map(`id`))
		user.Mobile = row.Str(res.Map(`mobile`))
		user.Nick = row.Str(res.Map(`nick`))
		user.Created = row.Str(res.Map(`created`))
		user.Parent = parent
		user.PointByReading = row.Int(res.Map(`reading`))
		user.PointByShare = row.Int(res.Map(`_share`))
		user.PointByDownLoadApp = row.Int(res.Map(`_app`))
		user.DaySign = row.Int(res.Map(`daysign`))
		user.Point = fmt.Sprintf("%.2f", row.Float(res.Map(`total`)))
		user.DrawCash = fmt.Sprintf("%.2f", row.Float(res.Map(`cny`)))
		user.TenMinuteCount = row.Int(res.Map(`second_count`))
		user.OneHourCount = row.Int(res.Map(`minute_count`))
		user.Blocked = row.Int(res.Map(`blocked`))
		data.UserList = append(data.UserList, user)
		newArray = append(newArray, row.Str(res.Map(`id`)))
	}

	return newArray, data, nil
}
