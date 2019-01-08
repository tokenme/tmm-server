package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)

/*
3. 积分发放趋势图，默认显示一周，用户可选时间范围（近2周，近一个月，全部）
*/
func PointTrendHandler(c *gin.Context) {
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
	SUM(tmp.points),
	tmp.Date
FROM(
	SELECT
		DATE(inserted_at) AS date,
		SUM(points) AS points
	FROM 
		tmm.device_share_tasks
	WHERE 
		inserted_at > '%s' AND inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
	GROUP BY 
		date
	UNION ALL
	SELECT	
		DATE(app.inserted_at) AS date,
		SUM(app.points) AS points
	FROM
		tmm.device_app_tasks AS app
	WHERE
		app.inserted_at > '%s' AND app.status = 1 AND app.inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE) 
	GROUP BY 
		date
	UNION ALL
	SELECT
		DATE(inserted_at) AS date,	
		SUM(point) AS points
	FROM
		tmm.reading_logs
	WHERE 
		inserted_at > '%s' AND inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
	GROUP BY 
		date
	UNION ALL
	SELECT
		DATE(inserted_at) AS date,	
		SUM(bonus) AS points
	FROM
		tmm.invite_bonus
	WHERE 
		inserted_at > '%s' AND inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
	GROUP BY
		date
) AS tmp
GROUP BY tmp.date
ORDER BY tmp.date  
	`
	rows, _, err := db.Query(query,
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime),
		db.Escape(startTime), db.Escape(endTime))
	if CheckErr(err, c) {
		return
	}
	format := "%.0f"
	var indexName []string
	var valueList []string
	var data Data
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)

	dataMap := make(map[string]float64)
	for _, row := range rows {
		dataMap[row.Str(1)] = row.Float(0)
	}

	for {
		if tm.Equal(end) {
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, value))
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, 0.0))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf(format, value))
			tm = tm.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf(format, 0.0))
			tm = tm.AddDate(0, 0, 1)
		}
	}

	data.Title.Text = "积分发放"
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = "积分"
	data.Series = append(data.Series, Series{Data: valueList, Name: "积分", Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
