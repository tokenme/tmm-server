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

type UserActiveStats struct {
	NewUser     uint   `json:"new_user"`
	LoginNumber uint   `json:"login_number"`
	TodayActive Active `json:"today_active"`
	TwoActive   Active `json:"two_active"`
	ThreeActive Active `json:"three_active"`
}

type Active struct {
	TotalActive uint `json:"total_active"`
	DaySign     uint `json:"day_sign"`
	Share       uint `json:"share"`
	DownLoadApp uint `json:"down_load_app"`
	Read        uint `json:"read"`
}

func UserFunnelStatsHandler(c *gin.Context) {
	startDate := c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`))
	tm, _ := time.Parse(`2006-01-02`, startDate)
	nextDate := tm.AddDate(0, 0, 1).Format("2006-01-02")
	thirdDate := tm.AddDate(0, 0, 2).Format("2006-01-02")
	endDate := tm.AddDate(0, 0, 3).Format("2006-01-02")
	query := `
SELECT
COUNT(*) AS total_users,
SUM(IF(logined, 1, 0)) AS logined_users,
SUM(IF(today_reading OR today_st OR today_at OR today_checkin, 1, 0)) today_active_users,
SUM(IF(today_reading, 1, 0)) today_reading_users,
SUM(IF(today_st, 1, 0)) today_st_users,
SUM(IF(today_at, 1, 0)) today_at_users,
SUM(IF(today_checkin, 1, 0)) today_checkin_users,

SUM(IF(nextday_reading OR nextday_st OR nextday_at OR nextday_checkin, 1, 0)) nextday_active_users,
SUM(IF(nextday_reading, 1, 0)) nextday_reading_users,
SUM(IF(nextday_st, 1, 0)) nextday_st_users,
SUM(IF(nextday_at, 1, 0)) nextday_at_users,
SUM(IF(nextday_checkin, 1, 0)) nextday_checkin_users,

SUM(IF(thirdday_reading OR thirdday_st OR thirdday_at OR thirdday_checkin, 1, 0)) thirdday_active_users,
SUM(IF(thirdday_reading, 1, 0)) thirdday_reading_users,
SUM(IF(thirdday_st, 1, 0)) thirdday_st_users,
SUM(IF(thirdday_at, 1, 0)) thirdday_at_users,
SUM(IF(thirdday_checkin, 1, 0)) thirdday_checkin_users
FROM (
SELECT
    us.id,
    EXISTS (SELECT 1 FROM tmm.devices AS d WHERE d.user_id=us.id LIMIT 1) AS logined,
    %s,
    %s,
    %s
FROM
    ucoin.users AS us
WHERE created BETWEEN '%s' AND '%s') AS t`
	db := Service.Db
	rows, _, err := db.Query(query, userFunnelStatsSubQuery(startDate, nextDate, "today"), userFunnelStatsSubQuery(nextDate, thirdDate, "nextday"), userFunnelStatsSubQuery(thirdDate, endDate, "thirdday"), startDate, nextDate)
	if CheckErr(err, c) {
		return
	}
	var stats UserActiveStats
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Data:    stats,
			Message: admin.API_OK,
			Code:    0,
		})
	}
	row := rows[0]
	stats = UserActiveStats{
		NewUser:     row.Uint(0),
		LoginNumber: row.Uint(1),
		TodayActive: Active{
			TotalActive: row.Uint(2),
			Read:        row.Uint(3),
			Share:       row.Uint(4),
			DownLoadApp: row.Uint(5),
			DaySign:     row.Uint(6),
		},
		TwoActive: Active{
			TotalActive: row.Uint(7),
			Read:        row.Uint(8),
			Share:       row.Uint(9),
			DownLoadApp: row.Uint(10),
			DaySign:     row.Uint(11),
		},
		ThreeActive: Active{
			TotalActive: row.Uint(12),
			Read:        row.Uint(13),
			Share:       row.Uint(14),
			DownLoadApp: row.Uint(15),
			DaySign:     row.Uint(16),
		},
	}

	c.JSON(http.StatusOK, admin.Response{
		Data:    stats,
		Message: admin.API_OK,
		Code:    0,
	})
}

func userFunnelStatsSubQuery(startDate string, endDate string, prefix string) string {
	return fmt.Sprintf(`EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id=us.id AND (rl.inserted_at BETWEEN '%s' AND '%s' OR rl.updated_at BETWEEN '%s' AND '%s') LIMIT 1) AS %s_reading,
    EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst INNER JOIN tmm.devices AS d ON (d.id=dst.device_id) WHERE d.user_id=us.id AND dst.inserted_at BETWEEN '%s' AND '%s' LIMIT 1) AS %s_st,
    EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat INNER JOIN tmm.devices AS d ON (d.id=dat.device_id) WHERE d.user_id=us.id AND dat.inserted_at BETWEEN '%s' AND '%s' LIMIT 1) AS %s_at,
    EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id=us.id AND dbl.updated_on>='%s' AND dbl.updated_on<=DATE_ADD('%s', INTERVAL dbl.days-1 DAY) LIMIT 1) AS %s_checkin`, startDate, endDate, startDate, endDate, prefix, startDate, endDate, prefix, startDate, endDate, prefix, startDate, startDate, prefix)
}
