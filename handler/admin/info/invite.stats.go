package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
	"time"
)

func InviteStatsHandler(c *gin.Context) {
	db := Service.Db
	var req StatsRequest
	var startTime, endTime string
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var when []string
	endTime = time.Now().Format("2006-01-02 15:04:05")
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(` AND ic.inserted_at >= '%s'`, db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		when = append(when, fmt.Sprintf(` AND ic.inserted_at >= '%s'`, db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		when = append(when, fmt.Sprintf(` AND ic.inserted_at <= '%s'`, db.Escape(endTime)))
	}
	var top10 string
	if req.Top10 {
		top10 = " LIMIT 10"
	}
	query := `
	SELECT
		wx.user_id AS id ,
		wx.nick AS nickname,
		COUNT(*) AS total,
		us.mobile AS mobile,
		SUM(ic.bonus) AS bouns
	FROM tmm.invite_bonus AS ic
	INNER JOIN tmm.wx AS wx ON  wx.user_id = ic.user_id 
	INNER JOIN ucoin.users AS us ON (us.id = wx.user_id)
	WHERE ic.task_id = 0 %s
	GROUP BY wx.user_id   
	ORDER BY total DESC 
	%s`
	rows, _, err := db.Query(query, strings.Join(when, " "), top10)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var info InviteStats
	for _, row := range rows {
		inviteCount := row.Int(2)
		if req.Top10 {
			user := &Users{
				InviteCount: inviteCount,
			}
			user.InviteBonus = row.Int(row.Int(4))
			user.Mobile = row.Str(3)
			user.Id = row.Uint64(0)
			user.Nick = row.Str(1)

			info.Top10 = append(info.Top10, user)
		}
		info.InviteCount = info.InviteCount + inviteCount
	}
	info.Numbers = len(rows)
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Title = "邀请排行榜"
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
