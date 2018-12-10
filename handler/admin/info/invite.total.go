package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"

	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

const totalInviteKey = `info-total-invite`

func TotalInviteHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, totalInviteKey))

	if context != nil && err ==nil{
		var data Data
		if !CheckErr(json.Unmarshal(context, &data), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    data,
			})
			return
		}
		return
	}
	query := `SELECT 
	SUM(bonus) AS cost, 
	COUNT(*)
	FROM tmm.invite_bonus 
	WHERE task_id = 0`
	var total TotalInvite
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	cost, err := decimal.NewFromString(rows[0].Str(0))
	if CheckErr(err, c) {
		return
	}
	total.TotalCost = cost
	total.TotalInviteCount = rows[0].Int(1)
	bytes, err := json.Marshal(&total)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, totalInviteKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    total,
	})
}
