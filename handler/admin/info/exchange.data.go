package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
)

const exchangeDataKey = `info-data-exchange`

func ExchangeDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	Context, err := redisConn.Do(`GET`, exchangeDataKey)
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
		tmp.user_id, 
		FLOOR((tmp.times-1)/5) * 5 + 1 AS l
FROM(			
	 SELECT
		COUNT(1) AS times ,
		er.user_id
	 FROM
		tmm.exchange_records AS er 
	 WHERE 
		er.status = 1
	 GROUP BY 
		er.user_id 
) AS tmp
	 GROUP BY 
		user_id
)AS tmp 
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
		Name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+5)
		indexName = append(indexName, Name)
	}
	data := Data{
		Title:     "交换",
		IndexName: indexName,
		Value:     valueList,
	}
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, exchangeDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
