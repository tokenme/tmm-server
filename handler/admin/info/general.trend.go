package info

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
)

const (
	GeneralTaskCount = iota
	GeneralUserCount
)

func GeneralTaskTrendHandler(c *gin.Context) {
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
	COUNT(DISTINCT  dev.user_id),
	COUNT(1) AS _count,
	DATE(dgt.inserted_at) AS record_on 
FROM 
	tmm.device_general_tasks  AS dgt  
LEFT JOIN tmm.devices AS dev ON (dev.id = dgt.device_id)
WHERE dgt.inserted_at BETWEEN '%s'  AND '%s'
GROUP BY record_on
	`
	db := Service.Db
	rows, _, err := db.Query(query, startTime, e.AddDate(0, 0, 1).Format(`2006-01-02`))
	if CheckErr(err, c) {
		return
	}

	var indexName, valueList []string
	var title string
	dataMap := make(map[string]int)
	format := "%d"
	for _, row := range rows {
		switch req.Type {
		case GeneralTaskCount:
			title = "普通任务数"
			dataMap[row.Str(2)] = row.Int(1)
		case GeneralUserCount:
			title = "普通任务用户数"
			dataMap[row.Str(2)] = row.Int(0)
		}
	}

	for {
		if s.Equal(e) {
			if value, ok := dataMap[s.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, s.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, value))
			} else {
				indexName = append(indexName, s.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, 0))
			}
			break
		}
		if value, ok := dataMap[s.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, s.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf(format, value))
			s = s.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, s.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf(format, 0))
			s = s.AddDate(0, 0, 1)
		}
	}

	var data Data
	data.Title.Text = title
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = title
	data.Series = append(data.Series, Series{Data: valueList, Name: title, Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})

}
