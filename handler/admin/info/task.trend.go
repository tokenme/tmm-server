package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const (
	FinishTaskUser = iota
	FinishTaskNumber
)

type TrendData struct {
	Title  string   `json:"title"`
	Yaxis  Axis     `json:"yAxis"`
	Xaxis  Axis     `json:"xAxis"`
	Series []Series `json:"series"`
}

func TaskTrendHandler(c *gin.Context) {
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
	query := `
SELECT
	tmp.date,
	tmp.value
FROM(
%s
	) AS tmp
`
	taskquery := `
	SELECT 
		SUM(tmp.times) AS '%s',
		tmp.date AS date,
		COUNT(DISTINCT tmp.user_id) AS '%s'
	FROM 
		ucoin.users  AS u 
	LEFT JOIN (
		SELECT 		
			DATE(sha.inserted_at) AS date ,
			dev.user_id AS user_id ,
			COUNT(1) AS times
		FROM 
			tmm.device_share_tasks  AS sha 
		INNER JOIN 
			tmm.devices AS dev ON  (dev.id = sha.device_id)
		WHERE 
			sha.inserted_at > '%s' AND sha.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
		GROUP BY 
			date,dev.user_id 
	UNION ALL 
		SELECT 
			DATE(app.inserted_at) AS date ,
			dev.user_id AS user_id,
			COUNT(1) AS times
		FROM 
			tmm.device_app_tasks  AS app 
		INNER JOIN 
			tmm.devices AS dev ON  (dev.id = app.device_id)
		WHERE 
			app.inserted_at > '%s' AND app.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
		GROUP BY date,dev.user_id 
	UNION ALL 
		SELECT 
			DATE(inserted_at) AS date,
			user_id AS user_id ,
			COUNT(1) AS times
		FROM 
			reading_logs 
		WHERE 
			inserted_at > '%s' AND inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
		GROUP BY 
			date,user_id 
	) AS tmp ON (tmp.user_id = u.id)
	GROUP BY 
		date `
	var data TrendData
	var yaxisName string
	if req.Type == FinishTaskNumber {
		yaxisName = "次数"
	} else {
		yaxisName = "人数"
	}
	switch req.Type {
	case FinishTaskUser:
		data.Title = "任务完成人数"
		finishTaskUser := fmt.Sprintf(taskquery, "not_use_total_times", "value",
			db.Escape(startTime), db.Escape(endTime),
			db.Escape(startTime), db.Escape(endTime),
			db.Escape(startTime), db.Escape(endTime), )
		query = fmt.Sprintf(query, finishTaskUser)
	case FinishTaskNumber:
		data.Title = "任务完成量"
		finishTaskNumber := fmt.Sprintf(taskquery, "value", "not_use_total_users_number",
			db.Escape(startTime), db.Escape(endTime),
			db.Escape(startTime), db.Escape(endTime),
			db.Escape(startTime), db.Escape(endTime), )
		query = fmt.Sprintf(query, finishTaskNumber)
	default:
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    data,
		})
		return
	}
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	var indexName, valueList []string

	dataMap := make(map[string]int)
	for _, row := range rows {
		dataMap[row.Str(0)] = row.Int(1)
	}
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)

	for {
		if tm.Equal(end) {
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", value))
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", 0))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", value))
			tm = tm.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", 0))
			tm = tm.AddDate(0, 0, 1)
		}
	}

	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = yaxisName
	data.Series = append(data.Series, Series{Data: valueList, Name: data.Title, Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})

}
