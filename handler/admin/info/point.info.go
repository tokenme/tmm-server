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

func PointInfoHandler(c *gin.Context) {

	db := Service.Db
	var req InfoRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var shareWhen []string
	var appTaskWhen []string
	var startTime, endTime string
	var top10 string
	endTime = time.Now().Format("2006-01-02 ")
	if req.StartTime != "" {
		startTime = req.StartTime
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
	}else{
		startTime = time.Now().AddDate(0,0,-7).Format("2006-01-02 ")
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at <= '%s' ", db.Escape(endTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at <= '%s' ", db.Escape(endTime)))
	}else{
		endTime = time.Now().String()
	}
	if req.Top10 {
		top10 = " LIMIT 10"
	}
	query := `
SELECT
	u.id AS id,
	u.nickname AS nick ,
	u.wx_nick AS wx_nick,
	tmp.points AS points 
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
ucoin.users AS u
INNER JOIN tmm.devices AS dev ON (dev.user_id = u.id )
WHERE 
		 tmp.device_id = dev.id 
GROUP BY 
		 u.id
ORDER BY points DESC %s`
	rows, res, err := db.Query(query, strings.Join(shareWhen, " "),
		strings.Join(appTaskWhen, " "),top10)
	if CheckErr(err, c) {
		return
	}
	var info PointInfo
	for _, row := range rows {
		Point, err := decimal.NewFromString(row.Str(res.Map(`points`)))
		if CheckErr(err, c) {
			return
		}
		if req.Top10 {
			user := &User{
				Id:     row.Int(res.Map(`id`)),
				Nick:   row.Str(res.Map(`nick`)),
				WxNick: row.Str(res.Map(`wx_nick`)),
				Point:  Point,
			}
			info.Top10 = append(info.Top10, user)
		}
		info.Point = info.Point.Add(Point)
	}

	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
