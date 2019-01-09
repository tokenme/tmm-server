package info

import (
	"github.com/gin-gonic/gin"
	"time"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const (
	InstallAppUser = iota
	InstallAppNumber
)

func AppTrendHandler(c *gin.Context) {
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
	tmp.value,
	tmp.before_users,
	tmp.before_times
FROM(
%s
	) AS tmp`
	appquery := `
	SELECT 
		DATE(app.inserted_at) AS date,
		COUNT(1) AS '%s'  ,
		COUNT(DISTINCT IF(NOT EXISTS (
			SELECT 
				1 
			FROM 
				tmm.device_app_tasks AS app 
			INNER JOIN tmm.devices AS _dev ON (_dev.id = app.device_id)
			WHERE 
				app.inserted_at < '%s' AND _dev.user_id = dev.user_id AND status = 1
		),dev.user_id,NULL)) AS '%s' ,
		beforeData.times AS before_times,
		beforeData.users AS before_users
	FROM 
		device_app_tasks AS app 
	INNER JOIN 
		tmm.devices AS dev ON( dev.id = app.device_id )
	LEFT JOIN (
		SELECT 
			COUNT(DISTINCT dev.user_id) AS users ,
			COUNT(1) AS times
		FROM 
			device_app_tasks  AS app 
		INNER JOIN 
			tmm.devices AS dev ON (dev.id = app.device_id) 
		WHERE 
			app.inserted_at < '%s' AND status = 1
	) AS beforeData ON 1 = 1
	WHERE 
		app.status = 1  AND app.inserted_at > '%s' 
		AND app.inserted_at < DATE_ADD('%s', INTERVAL 1 DAY) 
	GROUP BY 
		date 
	`
	var yaxisName string
	var data TrendData
	if req.Type == InstallAppUser {
		yaxisName = "人数"
	}else{
		yaxisName = "次数"
	}
	switch req.Type {
	case InstallAppUser:
		data.Title = "安装应用总人数"
		installAppUser := fmt.Sprintf(appquery,
			"not_use_total_install_times",
			db.Escape(startTime), "value", db.Escape(startTime),
			db.Escape(startTime), db.Escape(endTime))
		query = fmt.Sprintf(query, installAppUser)
	case InstallAppNumber:
		data.Title = "安装应用总次数"
		installAppNumber := fmt.Sprintf(appquery,
			"value",
			db.Escape(startTime), "not_use__total_users_number", db.Escape(startTime),
			db.Escape(startTime), db.Escape(endTime))
		query = fmt.Sprintf(query, installAppNumber)
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
	var beforeValue int
	var indexName, valueList []string
	if len(rows) != 0 {
		if req.Type == InstallAppUser {
			beforeValue = rows[0].Int(2)
		} else if req.Type == InstallAppNumber {
			beforeValue = rows[0].Int(3)
		}
	}
	dataMap := make(map[string]int)
	for _, row := range rows {
		beforeValue += row.Int(1)
		dataMap[row.Str(0)] = beforeValue
	}
	if len(rows) != 0 {
		if req.Type == InstallAppUser {
			beforeValue = rows[0].Int(2)
		} else if req.Type == InstallAppNumber {
			beforeValue = rows[0].Int(3)
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
				valueList = append(valueList, fmt.Sprintf("%d", beforeValue))
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			beforeValue = value
			valueList = append(valueList, fmt.Sprintf("%d", value))
			tm = tm.AddDate(0, 0, 1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf("%d", beforeValue))
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
