package info

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
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

	query := `
SELECT
	tmp.date,
	tmp.value
FROM(
%s
	) AS tmp`
	var data TrendData
	var yaxisName string
	if req.Type == ShareActicleNumber {
		yaxisName = "次数"
	} else {
		yaxisName = "人数"
	}
	switch req.Type {
	case ShareArticleUser:
		data.Title = "分享文章人数"
		shareArticleUser := fmt.Sprintf(`
	SELECT 
		COUNT(DISTINCT dev.user_id) AS value, 
		DATE(sha.inserted_at) AS date 
	FROM 
		tmm.device_share_tasks AS sha
	INNER JOIN 
		tmm.devices AS dev ON ( dev.id = sha.device_id )
	WHERE 
		sha.inserted_at > '%s'  AND sha.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY)
	GROUP BY date
	`, startTime, endTime)
		query = fmt.Sprintf(query, shareArticleUser)
	case ShareActicleNumber:
		data.Title = "分享文章量"
		shareArticleNumber := fmt.Sprintf(`
	SELECT 
		COUNT(1) AS value, 
		DATE(sha.inserted_at) AS date
	FROM 
		tmm.device_share_tasks AS sha
	WHERE 
		sha.inserted_at > '%s'  AND sha.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY)
	GROUP BY date 
	`, startTime, endTime)
		query = fmt.Sprintf(query, shareArticleNumber)
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
