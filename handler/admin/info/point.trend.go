package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
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
	query := `SELECT
    SUM(points) AS points,
    record_on
FROM
(
SELECT dst.points AS points, DATE(dst.inserted_at) AS record_on
FROM tmm.device_share_tasks AS dst
WHERE dst.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
UNION ALL
SELECT dat.points AS points, DATE(dat.inserted_at) AS record_on
FROM tmm.device_app_tasks AS dat
WHERE dat.status=1 AND dat.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
UNION ALL
SELECT rl.point AS points, DATE(rl.inserted_at) AS record_on
FROM tmm.reading_logs AS rl
WHERE rl.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
UNION ALL
SELECT ib.bonus AS points, DATE(ib.inserted_at) AS record_on
FROM tmm.invite_bonus AS ib
WHERE ib.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
UNION ALL
SELECT dgt.points AS points, DATE(dgt.inserted_at) AS record_on
FROM tmm.device_general_tasks AS dgt 
WHERE dgt.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY) AND dgt.status = 1 
) AS tmp
GROUP BY record_on
	`
	rows, _, err := db.Query(query, startTime, endTime, startTime, endTime, startTime, endTime, startTime, endTime, startTime, endTime)
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
