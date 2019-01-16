package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"strings"
	"fmt"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func ExchangeStatsHandler(c *gin.Context) {
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var when []string
	var startTime string
	db := Service.Db
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	exchangeKey := GetStatsKey(startTime, `exchange`)
	if !req.IsRefresh {
		var info PointStats
		if bytes, err := redis.Bytes(redisConn.Do(`GET`, exchangeKey)); err == nil && bytes != nil {
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
		redisConn.Do(`EXPIRE`, exchangeKey, 1)
	}

	query := `
SELECT 
	us.id AS id,
	wx.nick AS nickname, 
	tmp.tmm_add  AS tmm,
	us.mobile AS mobile
FROM(
	SELECT
		SUM(tmm) AS  tmm_add,
		er.user_id 
	FROM 
		tmm.exchange_records AS er
	WHERE 
		er.status = 1  AND direction=1 AND %s
	GROUP BY 
		user_id
) AS tmp ,ucoin.users AS us
LEFT JOIN tmm.wx AS wx ON (wx.user_id = us.id)
	WHERE tmp.user_id = us.id AND NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
	ORDER BY tmm DESC
LIMIT 10
`

	rows, res, err := db.Query(query, strings.Join(when, " AND "))
	if CheckErr(err, c) {
		return
	}

	var info ExchangeStats
	info.Title = `积分兑换UC数量排行榜`
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    info,
		})
		return
	}

	for _, row := range rows {
		user := &admin.User{}
		user.Tmm = fmt.Sprintf("%.2f", row.Float(res.Map(`tmm`)))
		user.Mobile = row.Str(res.Map(`mobile`))
		user.Id = row.Uint64(res.Map(`id`))
		user.Nick = row.Str(res.Map(`nickname`))
		info.Top10 = append(info.Top10, user)
	}

	if bytes, err := json.Marshal(&info); !CheckErr(err, c)  {
		redisConn.Do(`SET`, exchangeKey, bytes, `EX`, KeyAlive)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
