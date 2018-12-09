package info

import (
	. "github.com/tokenme/tmm/handler"
	"github.com/gin-gonic/gin"
	"fmt"
	"strings"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"time"
)

func TaskInfoHandler(c *gin.Context) {
	db := Service.Db
	var req InfoRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var shaTaskwhen []string
	var appTaskwhen []string
	var startTime, endTime string
	var top10 string
	endTime = time.Now().Format("2006-01-02 ")
	if req.StartTime != "" {
		startTime = req.StartTime
		shaTaskwhen = append(shaTaskwhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskwhen = append(appTaskwhen, fmt.Sprintf(" AND app.inserted_at >= '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0,0,-7).Format("2006-01-02 ")
		shaTaskwhen = append(shaTaskwhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskwhen = append(appTaskwhen, fmt.Sprintf(" AND app.inserted_at >= '%s' ", db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		shaTaskwhen = append(shaTaskwhen, fmt.Sprintf(" AND sha.inserted_at <= '%s' ", db.Escape(endTime)))
		appTaskwhen = append(appTaskwhen, fmt.Sprintf(" AND app.inserted_at <= '%s' ", db.Escape(endTime)))
	}

	if req.Top10 {
		top10 = "LIMIT 10"
	}
	query := `
SELECT 
	u.id AS id,
	u.mobile AS mobile,
	u.nickname AS nickname , 
	u.wx_nick AS wx_nick,
	SUM(tmp.point) AS point ,
	SUM(tmp.count_) AS _count
FROM (
	SELECT 
 		ud.user_id, 
 		SUM(tmp.total) AS count_,
 		SUM(tmp.point) AS point
	FROM(
	SELECT 
		sha.device_id,
		SUM(sha.points) AS point,
		COUNT(1)  AS total 
	FROM 
		tmm.device_share_tasks  AS sha
	WHERE
		sha.points > 0 %s
	GROUP BY 
		sha.device_id UNION ALL
	SELECT 
		app.device_id,
		SUM(app.points) AS point,
		COUNT(1)   AS total
	FROM 
		tmm.device_app_tasks AS app
	WHERE 
		app.status = 1 %s
	GROUP BY 
		app.device_id
) AS tmp
INNER JOIN tmm.user_devices AS ud ON (ud.device_id = tmp.device_id)
GROUP BY ud.user_id
	) AS tmp,ucoin.users AS u 
	WHERE tmp.user_id = u.id
	GROUP BY u.id 
	ORDER BY point DESC
%s`
	rows, res, err := db.Query(query, strings.Join(shaTaskwhen, ""),
		strings.Join(appTaskwhen, ""), top10)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var info TaskInfo
	for _, row := range rows {
		point, err := decimal.NewFromString(row.Str(res.Map(`point`)))
		if CheckErr(err, c) {
			return
		}
		count := row.Int(res.Map(`_count`))
		if req.Top10 {
			User := &User{
				Id:                 row.Int(res.Map(`id`)),
				Mobile:             row.Str(res.Map(`mobile`)),
				Nick:               row.Str(res.Map(`nickname`)),
				WxNick:             row.Str(res.Map(`wx_nick`)),
				Point:              point,
				CompletedTaskCount: count,
			}
			info.Top10 = append(info.Top10, User)
		}
		info.TaskCount = info.TaskCount + count
		info.TotalPoint = info.TotalPoint.Add(point)
	}
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
