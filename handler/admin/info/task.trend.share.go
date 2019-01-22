package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
)

const (
	ShareArticleUser = iota
	ShareActicleNumber
)

func ShareTrendHandler(c *gin.Context) {
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

	var data TrendData
	var yaxisName string
	switch req.Type {
	case ShareArticleUser:
		data.Title = "分享文章人数"
		yaxisName = "人数"

	case ShareActicleNumber:
		data.Title = "分享文章量"
		yaxisName = "次数"
	default:
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    data,
		})
		return
	}
	query := `
    SELECT
        DATE(dst.inserted_at) AS record_on,
        COUNT(1) AS times,
        COUNT(DISTINCT d.user_id) AS users
    FROM tmm.device_share_tasks AS dst
    INNER JOIN tmm.devices AS d ON ( d.id = dst.device_id )
    WHERE dst.inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
    GROUP BY record_on
    `
	rows, _, err := db.Query(query, startTime, endTime)
	if CheckErr(err, c) {
		return
	}

	var indexName, valueList []string

	dataMap := make(map[string]int)
	for _, row := range rows {
		switch req.Type {
		case ShareArticleUser:
			dataMap[row.Str(0)] = row.Int(2)
		case ShareActicleNumber:
			dataMap[row.Str(0)] = row.Int(1)
		default:
			continue
		}
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
