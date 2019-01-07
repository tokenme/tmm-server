package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func UserTrendHandler(c *gin.Context) {
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
		DATE(created) AS date,
		COUNT(1),
		tmp.users 
	FROM 
		ucoin.users
	LEFT JOIN (
	SELECT
		COUNT(1) AS users
	FROM 
		ucoin.users 
	WHERE 
		created < DATE('%s')
	) AS tmp ON (1 = 1)
	WHERE created > DATE('%s')   AND created< DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
	GROUP BY date
	ORDER BY date
`

	rows, _, err := db.Query(query, db.Escape(startTime), db.Escape(startTime), db.Escape(endTime))
	if CheckErr(err, c) {
		return
	}
	var beforeUserNumber int
	var indexName, valueList []string
	if len(rows) != 0 {
		beforeUserNumber = rows[0].Int(2)
	}
	dataMap := make(map[string]int)
	for _, row := range rows {
		beforeUserNumber += row.Int(1)
		dataMap[row.Str(0)] = beforeUserNumber
	}
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)
	if len(rows) != 0 {
		beforeUserNumber = rows[0].Int(2)
	}

	for {
		if tm.Equal(end) {
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", value))
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", beforeUserNumber))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			beforeUserNumber = value
			valueList = append(valueList, fmt.Sprintf("%d", value))
			tm = tm.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", beforeUserNumber))
			tm = tm.AddDate(0, 0, 1)
		}
	}
	var data Data
	data.Title.Text = "用户增长"
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
