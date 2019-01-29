package info

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
)

func PointStatsHandler(c *gin.Context) {

	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var startTime string
	db := Service.Db
	if req.StartTime != "" {
		startTime = req.StartTime
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	pointKey := GetStatsKey(startTime, `point`)
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
	var query string
	if startTime != "1970-1-1" {
		query = fmt.Sprintf(`SELECT tmp.user_id AS id, wx.nick AS nick, SUM(tmp.points) AS points, u.mobile AS mobile
FROM
(
    SELECT d.user_id AS user_id, SUM(dst.points) AS points
    FROM tmm.device_share_tasks AS dst
    INNER JOIN tmm.devices AS d ON (d.id=dst.device_id)
    WHERE dst.inserted_at>='%s'
	GROUP BY user_id
    UNION ALL
    SELECT d.user_id AS user_id, SUM(app.points) AS points
    FROM tmm.device_app_tasks AS app
    INNER JOIN tmm.devices AS d ON (d.id=app.device_id)
    WHERE app.status = 1 AND app.inserted_at>='%s'
	GROUP BY user_id 
    UNION ALL
    SELECT ib.user_id AS user_id, SUM(ib.bonus) AS points
    FROM tmm.invite_bonus AS ib
    WHERE ib.inserted_at>='%s'
	GROUP BY user_id 
    UNION ALL
    SELECT rl.user_id AS user_id, SUM(rl.point) AS points
    FROM tmm.reading_logs AS rl
    WHERE rl.inserted_at>='%s'
	GROUP BY user_id 
	UNION ALL  
	SELECT d.user_id AS user_id, SUM(dgt.points) AS points
	FROM tmm.device_general_tasks AS dgt
	INNER JOIN tmm.devices AS d ON (d.id = dgt.device_id)
	WHERE dgt.inserted_at>='%s' AND dgt.status = 1 
	GROUP BY user_id 
) AS tmp
INNER JOIN ucoin.users AS u ON (u.id=tmp.user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
WHERE NOT EXISTS (SELECT 1 FROM user_settings AS us WHERE us.blocked=1 AND us.user_id=tmp.user_id AND us.block_whitelist=0 LIMIT 1)
GROUP BY tmp.user_id ORDER BY points DESC LIMIT 10`, startTime, startTime, startTime, startTime, startTime)
	} else {
		query = `SELECT d.user_id AS id, wx.nick AS nick, SUM(d.points) AS points, u.mobile AS mobile
FROM tmm.devices AS d
INNER JOIN ucoin.users AS u ON (u.id=d.user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
WHERE NOT EXISTS (SELECT 1 FROM user_settings AS us WHERE us.blocked=1 AND us.user_id=d.user_id AND us.block_whitelist=0  LIMIT 1) AND u.id > 1
GROUP BY d.user_id ORDER BY points DESC LIMIT 10`
	}
	rows, res, err := db.Query(query)
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
