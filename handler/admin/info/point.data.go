package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"math"
)

const pointDataKey = `info-Data-points`

func PointDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, pointDataKey))
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
		return
	}
	query := `
SELECT
    COUNT(*) AS users,
    l
FROM (
    SELECT
	 	d.user_id,
	 	IF(SUM(d.points)=0,0,
		IF(SUM(d.points) >= 10000,
	FLOOR(((SUM(d.points)-1)/10000))*10000+1,
	FLOOR(((SUM(d.points)-1)/1000))*1000+1) ) AS l
    FROM tmm.devices AS d
	GROUP BY 
		d.user_id
) AS tmp
GROUP BY l 
ORDER BY l
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
		startPoint := row.Int(1)
		var endPoints int
		if startPoint < 0 {
			startPoint = 0
			endPoints = 1
		} else {
			log10 := math.Log10(float64(startPoint))
			endPoints = startPoint + int(math.Pow10(int(log10))) - 1
		}
		name := fmt.Sprintf(`%d-%d`, startPoint, endPoints)
		indexName = append(indexName, name)
	}
	var data Data
	data.Title.Text = "用户积分"
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "积分区间"
	data.Yaxis.Name = "人数"
	data.Series.Data = valueList
	data.Series.Name = "用户人数"
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, pointDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
