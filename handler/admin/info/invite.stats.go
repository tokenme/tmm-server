package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
	"time"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func InviteStatsHandler(c *gin.Context) {
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var when []string
	var startTime string
	db := Service.Db
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(` AND ib.inserted_at >= '%s'`, db.Escape(startTime)))
	} else {
		if req.Hours != 0 {
			startTime = time.Now().Add(-time.Hour * time.Duration(req.Hours)).Format("2006-01-02 15:04:05")
		} else {
			startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		}
		when = append(when, fmt.Sprintf(` AND ib.inserted_at >= '%s'`, db.Escape(startTime)))
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	inviteKey := GetStatsKey(req.StartTime, `invite`)
	if !req.IsRefresh {
		var info PointStats
		if bytes, err := redis.Bytes(redisConn.Do(`GET`, inviteKey)); err != nil && bytes != nil {
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
		redisConn.Do(`EXPIRE`, inviteKey, 1)
	}

	query := `
	SELECT
		u.id AS id ,
		wx.nick AS nickname,
		COUNT(DISTINCT ib.from_user_id) AS total,
		u.mobile AS mobile
	FROM tmm.invite_bonus AS ib
	INNER JOIN ucoin.users AS u ON (u.id = ib.user_id)
	LEFT  JOIN tmm.wx AS wx ON  wx.user_id = ib.user_id
	WHERE ib.task_type = 0 AND ib.user_id != ib.from_user_id %s  
	AND NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=u.id AND us.block_whitelist=0  LIMIT 1)
	GROUP BY u.id  
	ORDER BY total DESC 
	LIMIT 10`

	rows, _, err := db.Query(query, strings.Join(when, " "))
	if CheckErr(err, c) {
		return
	}

	var info InviteStats
	if req.Hours != 0 {
		info.Title = "邀请排行榜(二小时)"
	} else {
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
		user := &admin.User{}
		user.InviteCount = row.Int(2)
		user.Mobile = row.Str(3)
		user.Id = row.Uint64(0)
		user.Nick = row.Str(1)
		info.Top10 = append(info.Top10, user)
	}

	if bytes, err := json.Marshal(&info); !CheckErr(err, c) && req.Hours == 0{
		redisConn.Do(`SET`, inviteKey, bytes, `EX`, KeyAlive)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
