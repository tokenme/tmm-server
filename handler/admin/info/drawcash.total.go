package info

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

const totalDrawCashKey = `info-total-draw`

func TotalDrawCashHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, totalDrawCashKey))

	if context != nil && err == nil {
		var total TotalDrawCash
		if !CheckErr(json.Unmarshal(context, &total), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    total,
			})
			return
		}
		return
	}
	query := `
SELECT
    COUNT(*) AS users,
    SUM(tmp.cny) AS cny,
	SUM(tmp.total) AS total
FROM (
    SELECT 
 user_id, 
 SUM(cny) AS cny,
 SUM(total) AS total 
FROM(
		SELECT
            tx.user_id, 
			SUM( tx.cny ) AS cny,
			COUNT(1) AS total 
        FROM
            tmm.withdraw_txs AS tx
		WHERE
			tx.tx_status = 1
        GROUP BY
            tx.user_id UNION ALL
        SELECT
            pw.user_id, 
			SUM( pw.cny ) AS cny,
			COUNT(1) AS total 
        FROM
            tmm.point_withdraws AS pw
        GROUP BY pw.user_id
				) AS tmp
		GROUP BY user_id
) AS tmp
`
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	var total TotalDrawCash
	if CheckErr(err, c) {
		return
	}
	total.TotalCount = row.Int(res.Map(`total`))
	total.TotalUser = row.Int(res.Map(`users`))
	total.TotalMoney = fmt.Sprintf("%.2f",row.Float(res.Map(`cny`)))
	bytes, err := json.Marshal(&total)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, totalDrawCashKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    total,
	})

}
