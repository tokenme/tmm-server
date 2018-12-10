package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

const drawDataKey = `info-data-draw`

func DrawCashDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, drawDataKey))
	if context != nil  && err ==nil{
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
    COUNT(*) AS users,
    l
FROM (
    SELECT 
 user_id, 
 IF(SUM(cny)=0, 0,FLOOR((SUM(cny)-1)/100) * 100 + 1) AS l
FROM(
	SELECT
            tx.user_id, SUM( tx.cny ) AS cny
        FROM
            tmm.withdraw_txs AS tx
		WHERE 
			tx.tx_status = 1
        GROUP BY
            tx.user_id UNION ALL
        SELECT
            pw.user_id, SUM( pw.cny ) AS cny
        FROM
            tmm.point_withdraws AS pw
        GROUP BY pw.user_id
				) AS tmp
		GROUP BY user_id
) AS tmp
GROUP BY l ORDER BY l

`

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName []string
	var valueList []int
	for _, row := range rows {
		valueList = append(valueList, row.Int(0))
		name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+100)
		indexName = append(indexName, name)
	}
	data := Data{
		Title:     "提现金额占比",
		IndexName: indexName,
		Value:     valueList,
	}
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, drawDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}