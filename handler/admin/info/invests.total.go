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

const totalInvestsKey = `info-total-invest`

func TotalInvestHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, totalInvestsKey))
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
	query := `
SELECT 
	SUM(inv.points),
	COUNT(1) AS total  
FROM tmm.goods AS g 
INNER JOIN tmm.good_invests AS inv ON (inv.good_id = g.id) 
WHERE redeem_status = 0`
	var total TotalInvests
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	totalPoint, err := decimal.NewFromString(rows[0].Str(0))
	if CheckErr(err, c) {
		return
	}

	total.TotalPoint = totalPoint
	total.TotalGoodsCount = rows[0].Int(1)

	bytes, err := json.Marshal(&total)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, totalInvestsKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    total,
	})
}
