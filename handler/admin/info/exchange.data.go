package info

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

const exchangeDataKey = `info-data-exchange`

func ExchangeDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, exchangeDataKey))

	if context != nil && err == nil {
		var data Data
		if !CheckErr(json.Unmarshal(context, &data), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    data,
			})
			return
		}
	}
	query := `SELECT
    COUNT(1) AS users,
    l
FROM(
	 SELECT
        FLOOR((COUNT(1)-1)/10) * 10 + 1 AS l, er.user_id
     FROM tmm.exchange_records AS er
     WHERE er.status=1
     GROUP BY er.user_id
) AS tmp
GROUP BY l ORDER BY l
`
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName []string
	var valueList []string
	for _, row := range rows {
		valueList = append(valueList, row.Str(0))
		name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+10)
		indexName = append(indexName, name)
	}
	var data Data
	data.Title.Text = "兑换次数"
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "兑换次数区间"
	data.Yaxis.Name = "人数"
	data.Series = append(data.Series, Series{Data: valueList, Name: "人数"})
	data.Series = append(data.Series, Series{Data: GetPercentList(valueList), Name: "占比"})
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
