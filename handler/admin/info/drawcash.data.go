package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
)

const drawDataKey = `info-data-draw`

func DrawCashDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	Context, err := redisConn.Do(`GET`, drawDataKey)
	if CheckErr(err, c) {
		return
	}
	if Context != nil {
		var data Data
		if CheckErr(json.Unmarshal(Context.([]byte), &data), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    data,
			})
			return
		}
	}
	query := `SELECT
    COUNT(*) AS users,
    l
FROM (
    SELECT 
 user_id, 
 IF(SUM(cny)=0, 0,FLOOR((SUM(cny)-1)/50) * 50 + 1) AS l
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
		Name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+50)
		indexName = append(indexName, Name)
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
