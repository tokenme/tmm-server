package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
)

const investsDataKey = `info-data-invest`

func InvestsDataHandler(c *gin.Context) {
	db := Service.Db

	redisConn := Service.Redis.Master.Get()
	Context, err := redisConn.Do(`GET`, investsDataKey)
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
	query := `
SELECT 
	COUNT(*) AS total,
	l
FROM (
	SELECT 
	    g.user_id,IF(SUM(g.points)=0,0,FLOOR((SUM(g.points)-1)/50)*50+1) AS l
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
		indexName = append(indexName, fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+50))
	}
	data := Data{
		Title:     "商品投资",
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
