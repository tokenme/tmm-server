package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)

//1. 拉新好友数趋势图，默认显示一周，用户可选时间范围（近2周，近一个月，全部）
//2. 好友活跃数量趋势图，默认显示一周，用户可选时间范围（近2周，近一个月，全部）
func InviteTrendHandler(c *gin.Context) {
	db := Service.Db
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	startTime := time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
	endTime := time.Now().Format(`2006-01-02`)
	if req.StartTime != "" {
		startTime = req.StartTime
	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}
	s, _ := time.Parse(`2006-01-02`, startTime)
	e, _ := time.Parse(`2006-01-02`, endTime)
	if s.Unix() > e.Unix() {
		c.JSON(http.StatusOK, admin.Response{
			Code:    1,
			Message: "起始日期不能超过结束日期",
			Data:    nil,
		})
		return
	}
	query := `
SELECT
	%s,
	DATE(u.created) AS date
FROM 
	ucoin.users AS u
	INNER JOIN tmm.invite_codes AS inv_code ON (inv_code.user_id = u.id AND inv_code.parent_id > 0)
WHERE u.created > '%s'   AND u.created <  DATE_ADD('%s', INTERVAL 1 DAY)
GROUP BY date
	`
	var data Data
	if req.Type == 1 {
		data.Title.Text = "好友活跃数"
		query = fmt.Sprintf(query, fmt.Sprintf(`COUNT(distinct 
		IF(EXISTS(
		SELECT 
		1
		FROM tmm.devices AS dev 
		LEFT JOIN tmm.device_share_tasks AS sha 
	ON (sha.device_id = dev.id AND sha.inserted_at > '%s' AND sha.inserted_at <DATE_ADD('%s', INTERVAL 1 DAY))
		LEFT JOIN tmm.device_app_tasks AS app 
	ON (app.device_id = dev.id  AND app.inserted_at > '%s' AND app.inserted_at <DATE_ADD('%s', INTERVAL 1 DAY))
		LEFT JOIN tmm.reading_logs AS reading 
	ON (reading.user_id = dev.user_id AND (
		(reading.inserted_at > '%s' AND reading.inserted_at <DATE_ADD('%s', INTERVAL 1 DAY))
	OR  (reading.inserted_at > '%s' AND reading.inserted_at <DATE_ADD('%s', INTERVAL 1 DAY))))
		LEFT JOIN tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id AND updated_on >= '%s'  AND updated_on < DATE_ADD('%s', INTERVAL 1 DAY)) 
		WHERE dev.user_id = u.id AND (
		sha.task_id > 0 OR 
		app.task_id > 0 OR 
		reading.point > 0 OR
		daily.user_id > 0
		)
		LIMIT 1
		),u.id,NULL)
		)`, db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime)),
		db.Escape(startTime), db.Escape(endTime),)
	} else if req.Type == 0 {
		data.Title.Text = "拉新好友数"
		query = fmt.Sprintf(query, "COUNT(1)",db.Escape(startTime),db.Escape(endTime))
	}else {
		data.Title.Text = "用户活跃"
		query =fmt.Sprintf(`
	SELECT 
		tmp.value AS value,
		tmp.date AS date
	FROM (
		SELECT 
		COUNT(DISTINCT tmp.user_id) AS value,
		tmp.date AS date
	FROM 
	(
		SELECT  
	dev.user_id AS user_id ,
	DATE(sha.inserted_at) AS date 
	FROM 
	  tmm.device_share_tasks  AS sha 
	INNER JOIN tmm.devices AS dev ON dev.id = sha.device_id
	WHERE sha.inserted_at > '%s' AND sha.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY)
	GROUP BY user_id,date
	UNION ALL
	SELECT 
		dev.user_id AS user_id ,
		DATE(app.inserted_at) AS date 
	FROM 
		tmm.device_app_tasks  AS app 
	INNER JOIN tmm.devices AS dev ON dev.id = app.device_id
	WHERE app.inserted_at > '%s' AND app.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY)  
	GROUP BY user_id,date
	UNION ALL 
	SELECT 
		user_id  AS user_id ,
		DATE(reading.inserted_at) AS date 
	FROM 
		tmm.reading_logs AS reading 
	WHERE reading.inserted_at > '%s' AND reading.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY)
	GROUP BY user_id,date
	UNION ALL 
	SELECT
		user_id   AS user_id ,
		DATE(reading.updated_at) AS date 
	FROM 
		tmm.reading_logs AS reading 
	WHERE reading.updated_at > '%s' AND reading.updated_at < DATE_ADD('%s', INTERVAL 1 DAY)
	UNION ALL 
	SELECT 
		user_id  AS user_id ,
		DATE(updated_on) AS date 
	FROM 
		daily_bonus_logs
	WHERE updated_on >= '%s' AND  updated_on <= DATE_ADD('%s', INTERVAL 1 DAY)
	GROUP BY user_id,date
	) AS tmp 
	GROUP BY tmp.date 
	) AS tmp
	GROUP BY date 
`,
db.Escape(startTime),db.Escape(endTime),
db.Escape(startTime),db.Escape(endTime),
db.Escape(startTime),db.Escape(endTime),
db.Escape(startTime),db.Escape(endTime),
db.Escape(startTime),db.Escape(endTime),)

	}
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName []string
	var valueList []string
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)

	dataMap := make(map[string]int)
	for _, row := range rows {
		dataMap[row.Str(1)] = row.Int(0)
	}

	for {
		if tm.Equal(end){
			if value,ok:=dataMap[tm.Format(`2006-01-02`)];ok{
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", value))
			}else{
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", 0))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", value))
			tm = tm.AddDate(0,0,1)
		}else{
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", 0))
			tm = tm.AddDate(0,0,1)
		}
	}

	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = "人数"
	data.Series = append(data.Series, Series{Data: valueList, Name: "人数", Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
