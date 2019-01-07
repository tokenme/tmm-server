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
	DATE(inv.inserted_at) AS date
FROM 
	tmm.invite_bonus AS inv
WHERE inv.inserted_at > '%s'  AND task_type = 0 AND inv.inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
GROUP BY date
	`
	var data Data
	if req.Type == 1 {
		data.Title.Text = "好友活跃数"
		query = fmt.Sprintf(query, `COUNT(distinct 
		IF(EXISTS(
		SELECT 
		10
		FROM tmm.devices AS dev 
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id )
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND app.status = 1)
		LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id )
		WHERE dev.user_id = inv.from_user_id AND ( 
		sha.inserted_at > inv.inserted_at OR 
		app.inserted_at > inv.inserted_at OR    
		reading.inserted_at > inv.inserted_at  )
		LIMIT 1
		)	,inv.from_user_id,NULL)
		)`, db.Escape(startTime), db.Escape(endTime))
	} else {
		data.Title.Text = "拉新好友数"
		query = fmt.Sprintf(query, "COUNT(distinct inv.from_user_id)", db.Escape(startTime), db.Escape(endTime))
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
