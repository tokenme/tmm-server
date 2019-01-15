package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"time"
	"strings"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func PointStatsHandler(c *gin.Context) {

	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var shareWhen, appTaskWhen, when []string
	var startTime string
	db := Service.Db
	if req.StartTime != "" {
		startTime = req.StartTime
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
		when = append(when, fmt.Sprintf(" inserted_at >= '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		shareWhen = append(shareWhen, fmt.Sprintf(" AND sha.inserted_at >= '%s' ", db.Escape(startTime)))
		appTaskWhen = append(appTaskWhen, fmt.Sprintf(" AND  app.inserted_at >= '%s' ", db.Escape(startTime)))
		when = append(when, fmt.Sprintf("  inserted_at >= '%s' ", db.Escape(startTime)))
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	pointKey := GetStatsKey(req.StartTime, `point`)
	if !req.IsRefresh {
		var info PointStats
		if bytes, err := redis.Bytes(redisConn.Do(`GET`, pointKey)); err == nil && bytes != nil {
			if !CheckErr(json.Unmarshal(bytes, &info), c) {
				c.JSON(http.StatusOK, admin.Response{
					Code:    0,
					Message: admin.API_OK,
					Data:    info,
				})
			}
			return
		}
	} else {
		redisConn.Do(`EXPIRE`, pointKey, 1)
	}

	query := `
SELECT
	us.id AS id,
	wx.nick AS nick ,
	IFNULL(tmp.points,0)+IFNULL(inv.bonus,0)+IFNULL(reading.points,0) AS points,
	us.mobile AS mobile
FROM ucoin.users AS us
LEFT JOIN (
	SELECT 
	SUM(tmp.points) AS points,
	dev.user_id   AS user_id
	FROM   tmm.devices AS dev
	LEFT JOIN (
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
		 app.status = 1 %s
	GROUP BY
     	 app.device_id  
	) AS tmp ON (tmp.device_id = dev.id)
	GROUP BY dev.user_id 
) AS tmp ON (tmp.user_id = us.id )
LEFT JOIN tmm.wx AS wx ON (wx.user_id = us.id)
LEFT JOIN (SELECT SUM(bonus) AS bonus,user_id AS user_id FROM tmm.invite_bonus  
 		  WHERE %s
   		  GROUP BY user_id)AS inv ON (inv.user_id = us.id )  
LEFT JOIN (SELECT 
		   SUM(point) AS points ,
			user_id AS user_id 
		FROM tmm.reading_logs 
		WHERE %s 
		GROUP BY user_id) AS reading ON(reading.user_id = us.id)
WHERE IFNULL(tmp.points,0)+IFNULL(inv.bonus,0)+IFNULL(reading.points,0)  > 0 
AND NOT EXISTS  
		(SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
GROUP BY 
		 us.id
ORDER BY points DESC 
LIMIT 10`

	rows, res, err := db.Query(query, strings.Join(shareWhen, " "),
		strings.Join(appTaskWhen, " "), strings.Join(when, " "), strings.Join(when, " "))
	if CheckErr(err, c) {
		return
	}
	var info PointStats
	info.Title = "积分排行榜"
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
			Data:    info,
		})
		return
	}

	for _, row := range rows {
		user := &admin.User{}
		user.Point = fmt.Sprintf("%.0f", row.Float(res.Map(`points`)))
		user.Mobile = row.Str(res.Map(`mobile`))
		user.Id = row.Uint64(res.Map(`id`))
		user.Nick = row.Str(res.Map(`nick`))
		info.Top10 = append(info.Top10, user)
	}

	if bytes, err := json.Marshal(&info); !CheckErr(err, c) {
		redisConn.Do(`SET`, pointKey, bytes, `EX`, KeyAlive)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
