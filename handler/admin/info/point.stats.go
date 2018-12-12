package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"time"
	"strings"
	"github.com/shopspring/decimal"
)

func PointStatsHandler(c *gin.Context) {

	db := Service.Db
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var shareWhen, appTaskWhen, invWhen []string
	var startTime, endTime string
	var top10 string
	endTime = time.Now().Format("2006-01-02 15:04:05")
	if req.StartTime != "" {
		startTime = req.StartTime
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
		invWhen = append(invWhen, fmt.Sprintf("AND inv.inserted_at >= '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
		invWhen = append(invWhen, fmt.Sprintf("AND inv.inserted_at >= '%s' ", db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at <= '%s' ", db.Escape(endTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at <= '%s' ", db.Escape(endTime)))
		invWhen = append(invWhen, fmt.Sprintf("AND inv.inserted_at <= '%s' ", db.Escape(endTime)))
	}

	if req.Top10 {
		top10 = " LIMIT 10"
	}
	query := `
SELECT
	wx.user_id AS id,
	wx.nick AS nick ,
	(tmp.points+ SUM(inv.bonus))AS points,
	us.mobile AS mobile
FROM(
	SELECT 
		 sha.device_id, 
		 SUM(sha.points) AS points
	FROM 
		 tmm.device_share_tasks AS sha	
	WHERE 
		 sha.points > 0 %s
	GROUP BY
       sha.device_id UNION ALL
	SELECT 
		 app.device_id, 
		 SUM(app.points) AS points
	FROM 
		 tmm.device_app_tasks AS app   
	WHERE
		 app.status = 0 %s
	GROUP BY
     	 app.device_id   
) AS tmp,
tmm.wx AS wx
INNER JOIN tmm.devices AS dev ON (dev.user_id = wx.user_id)
INNER JOIN ucoin.users AS us ON (us.id = wx.user_id)
INNER JOIN invite_bonus AS inv ON (inv.user_id = wx.user_id  %s)
WHERE 
		 tmp.device_id = dev.id 
GROUP BY 
		 wx.user_id
ORDER BY points DESC  %s`
	rows, res, err := db.Query(query, strings.Join(shareWhen, " "),
		strings.Join(appTaskWhen, " "), strings.Join(invWhen, " "),
	top10)
	if CheckErr(err, c) {
		return
	}
	var info PointStats
	for _, row := range rows {
		Point, err := decimal.NewFromString(row.Str(res.Map(`points`)))
		if CheckErr(err, c) {
			return
		}
		if req.Top10 {
			user := &Users{
				Point: Point.Ceil(),
			}
			user.Mobile = row.Str(res.Map(`mobile`))
			user.Id = row.Uint64(res.Map(`id`))
			user.Nick = row.Str(res.Map(`nick`))
			info.Top10 = append(info.Top10, user)
		}
		info.Point = info.Point.Add(Point)
	}
	info.Numbers = len(rows)
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Title = "积分排行榜"
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
