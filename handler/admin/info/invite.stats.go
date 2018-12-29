package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
	"time"
	"github.com/shopspring/decimal"
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
		if req.Hours != 0 {
			startTime = time.Now().Add(-time.Hour * time.Duration(req.Hours)).Format("2006-01-02 15:04:05")
		} else {
			startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		}
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
	LEFT  JOIN tmm.wx AS wx ON  wx.user_id = ic.user_id 
	INNER JOIN ucoin.users AS us ON (us.id = ic.user_id)
	WHERE ic.task_type = 0 %s  
	AND NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
	GROUP BY us.id  
	ORDER BY total DESC 
	%s`
	rows, _, err := db.Query(query, strings.Join(when, " "), top10)
	if CheckErr(err, c) {
		return
	}
	var info InviteStats
	info.Numbers = len(rows)
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	if req.Hours != 0 {
		info.Title = "邀请排行榜(二小时)"
	}else{
		info.Title = "邀请排行榜"
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
			Data:    info,
		})
		return
	}
	for _, row := range rows {
		bouns,err:=decimal.NewFromString(row.Str(4))
		if CheckErr(err,c){
			return
		}
		inviteCount := row.Int(2)
		if req.Top10 {
			user := &admin.Users{
				InviteCount: inviteCount,
			}
			user.InviteBonus = bouns.Ceil()
			user.Mobile = row.Str(3)
			user.Id = row.Uint64(0)
			user.Nick = row.Str(1)

			info.Top10 = append(info.Top10, user)
		}
		info.InviteCount = info.InviteCount + inviteCount
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
