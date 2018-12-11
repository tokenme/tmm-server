package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

const investsDataKey = `info-data-invest`

func InvestsDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, investsDataKey))
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
	COUNT(*) AS total,
	l
FROM (
	SELECT 
	    g.user_id,IF(SUM(g.points)=0,0,FLOOR((SUM(g.points)-1)/1000)*1000+1) AS l
	FROM 
		tmm.good_invests AS g
	WHERE
		g.redeem_status = 0
	GROUP BY 
		g.user_id 
) AS tmp
GROUP BY l ORDER BY l
`

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var indexName []string
	var valueList []int
	for _, row := range rows {
		valueList = append(valueList, row.Int(0))
		indexName = append(indexName, fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+1000))
	}
	data := Data{
		Title:     "商品投资 - X轴:积分 - Y轴:商品数量 ",
		IndexName: indexName,
		Value:     valueList,
	}
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, investsDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
