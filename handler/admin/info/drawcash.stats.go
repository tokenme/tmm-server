package info

import (
	. "github.com/tokenme/tmm/handler"
	"github.com/gin-gonic/gin"
	"fmt"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

func DrawCashStatsHandler(c *gin.Context) {
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var txwhen,ptwhen []string
	var startTime string
	db := Service.Db
	if req.StartTime != "" {
		startTime = req.StartTime
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  > '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  > '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  > '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  > '%s' ", db.Escape(startTime)))
	}

	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	drawcashKey := GetStatsKey(req.StartTime, `drawcash`)
	if !req.IsRefresh {
		var info DrawCashStats
		if bytes, err := redis.Bytes(redisConn.Do(`GET`, drawcashKey)); err != nil && bytes != nil {
			if !CheckErr(json.Unmarshal(bytes, &info), c) {
				c.JSON(http.StatusOK, admin.Response{
					Code:    0,
					Message: admin.API_OK,
					Data:    info,
				})
			}
			return
		}
	}else{
		redisConn.Do(`EXPIRE`,drawcashKey,1)
	}

	query := `SELECT 
	us.id AS id ,
	wx.nick AS nickname , 
	IFNULL(tmp.cny,0) AS cny,
	us.mobile AS mobile
FROM (
 SELECT 
 user_id, 
 SUM(cny) AS cny
FROM(
	SELECT
            tx.user_id, 
			SUM( tx.cny ) AS cny
        FROM
            tmm.withdraw_txs AS tx
		WHERE
			tx.tx_status = 1 %s
        GROUP BY
            tx.user_id UNION ALL
        SELECT
            pw.user_id, 
			SUM( pw.cny ) AS cny
        FROM
            tmm.point_withdraws AS pw
        WHERE 
			pw.cny > 0 %s
		GROUP BY pw.user_id
				) AS tmp
		GROUP BY user_id
) AS tmp
INNER JOIN ucoin.users AS us ON (us.id = tmp.user_id)
LEFT JOIN tmm.wx AS wx  ON (wx.user_id = us.id)
WHERE  NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
GROUP BY us.id 
ORDER BY cny DESC 
LIMIT 10`

	rows, res, err := db.Query(query, strings.Join(txwhen, ""), strings.Join(ptwhen, ""))
	if CheckErr(err, c) {
		return
	}

	var info DrawCashStats
	info.Title = `提现排行榜`
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    info,
		})
		return
	}
	for _, row := range rows {
			user := &admin.User{
				DrawCash: fmt.Sprintf("%.2f", row.Float(res.Map(`cny`))),
			}
			user.Mobile = row.Str(res.Map(`mobile`))
			user.Id = row.Uint64(res.Map(`id`))
			user.Nick = row.Str(res.Map(`nickname`))
			info.Top10 = append(info.Top10, user)
	}

	if bytes, err := json.Marshal(&info); !CheckErr(err,c){
		redisConn.Do(`SET`, drawcashKey, bytes,`EX`,KeyAlive)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
