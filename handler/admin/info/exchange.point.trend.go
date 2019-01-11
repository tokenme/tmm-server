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
4.  积分兑换趋势图，默认显示一周，用户可选时间范围（近2周，近一个月，全部）
*/

const (
	PointToUc = iota
	UcToPoint
	UcToPointUserCount
)

func ExchangeTrendHandler(c *gin.Context) {
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
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)
	if Check(tm.Unix() > end.Unix(), `开始日期不能超过结束日期`, c) {
		return
	}

	query := `
SELECT
	%s,
	DATE(inserted_at) AS date
FROM
	tmm.exchange_records
WHERE
	inserted_at > '%s' AND inserted_at <  DATE_ADD('%s', INTERVAL 1 DAY) 
	AND status = 1 AND direction = %d
GROUP BY date
ORDER BY date 
`

	var direction int
	var field string
	var title string
	var yaxisName string
	if req.Type == PointToUc {
		direction = 1
		field = "SUM(tmm)"
		title = "积分兑UC数量"
		yaxisName = "数量"
	} else if req.Type == UcToPoint {
		title = "UC兑积分数量"
		direction = -1
		field = "SUM(points)"
	} else if req.Type == UcToPointUserCount {
		title = "UC兑积分人数"
		field = "COUNT(DISTINCT user_id)"
		direction = -1
		yaxisName = "人数"
	}

	db := Service.Db
	rows, _, err := db.Query(query, db.Escape(field), db.Escape(startTime), db.Escape(endTime), direction)
	if CheckErr(err, c) {
		return
	}

	var indexName []string
	var valueList []string
	var data TrendData

	dataMap := make(map[string]float64)
	for _, row := range rows {
		dataMap[row.Str(1)] = row.Float(0)
	}

	for {
		if tm.Equal(end) {
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%.0f", value))
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%.0f", 0.0))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%.0f", value))
			tm = tm.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%.0f", 0.0))
			tm = tm.AddDate(0, 0, 1)
		}
	}

	data.Title = title
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = yaxisName
	data.Series = append(data.Series, Series{Data: valueList, Name: data.Title, Type: "line"})

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    &data,
	})
}
