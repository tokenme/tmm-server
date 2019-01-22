package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
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
	var query string
	var data Data
	if req.Type == 1 {
		data.Title.Text = "好友活跃数"
		query = fmt.Sprintf(`SELECT
    COUNT(DISTINCT user_id) AS users,
    record_on
FROM
(
SELECT d.user_id AS user_id, DATE(dst.inserted_at) AS record_on
FROM tmm.device_share_tasks AS dst
INNER JOIN tmm.devices AS d ON (d.id=dst.device_id)
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=d.user_id AND ic.parent_id>0)
WHERE dst.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT d.user_id AS user_id, DATE(dat.inserted_at) AS record_on
FROM tmm.device_app_tasks AS dat
INNER JOIN tmm.devices AS d ON (d.id=dat.device_id)
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=d.user_id AND ic.parent_id>0)
WHERE dat.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT rl.user_id AS user_id, DATE(rl.inserted_at) AS record_on
FROM tmm.reading_logs AS rl
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=rl.user_id AND ic.parent_id>0)
WHERE rl.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT dbl.user_id AS user_id, DATE(dbl.updated_on) AS record_on
FROM tmm.daily_bonus_logs AS dbl
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=dbl.user_id AND ic.parent_id>0)
WHERE dbl.updated_on BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
) AS tmp
GROUP BY record_on`, startTime, endTime, startTime, endTime, startTime, endTime, startTime, endTime)
	} else if req.Type == 0 {
		data.Title.Text = "拉新好友数"
		query = fmt.Sprintf(`SELECT
    COUNT(u.id),
    DATE(u.created) AS record_on
FROM
    ucoin.users AS u
    INNER JOIN tmm.invite_codes AS ic ON (ic.user_id = u.id AND ic.parent_id > 0)
WHERE u.created BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
GROUP BY record_on`, startTime, endTime)
	} else {
		data.Title.Text = "用户活跃"
		query = fmt.Sprintf(`SELECT
    COUNT(DISTINCT user_id) AS users,
    record_on
FROM
(
SELECT d.user_id AS user_id, DATE(dst.inserted_at) AS record_on
FROM tmm.device_share_tasks AS dst
INNER JOIN tmm.devices AS d ON (d.id=dst.device_id)
WHERE dst.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT d.user_id AS user_id, DATE(dat.inserted_at) AS record_on
FROM tmm.device_app_tasks AS dat
INNER JOIN tmm.devices AS d ON (d.id=dat.device_id)
WHERE dat.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT rl.user_id AS user_id, DATE(rl.inserted_at) AS record_on
FROM tmm.reading_logs AS rl
WHERE rl.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT dbl.user_id AS user_id, DATE(dbl.updated_on) AS record_on
FROM tmm.daily_bonus_logs AS dbl
WHERE dbl.updated_on BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
) AS tmp
GROUP BY record_on`, startTime, endTime, startTime, endTime, startTime, endTime, startTime, endTime)

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
	data.Yaxis.Name = "人数"
	data.Series = append(data.Series, Series{Data: valueList, Name: "人数", Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
