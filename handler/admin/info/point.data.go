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
const AnotherKey = `info-data-points-another`

func PointDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	types := c.DefaultQuery(`type`, "-1")
	defer redisConn.Close()
	var context []byte
	var err error
	if types != "-1" {
		context, err = redis.Bytes(redisConn.Do(`GET`, AnotherKey))
	} else {
		context, err = redis.Bytes(redisConn.Do(`GET`, pointDataKey))
	}
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
    l,
	SUM(tmp.total_points)
FROM (
    SELECT
	 	d.id,
	 	IF (d.total_points=0,0,
		IF(d.total_points  >= 10000,
	FLOOR(d.total_points/10000)*10000+1,
	FLOOR(d.total_points/1000)*1000+1) ) AS l,
	d.total_points
  FROM tmm.top_points_users AS d
WHERE NOT EXISTS
	(SELECT 1 FROM tmm.user_settings AS us  
	WHERE us.blocked= 1 AND us.user_id=d.id AND us.block_whitelist=0  LIMIT 1)
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
	var valueList []string
	var data Data
	if types == "-1" {
		for _, row := range rows {
			valueList = append(valueList, row.Str(0))
			startPoint := row.Int(1)
			var endPoints int
			if startPoint == 1 {
				startPoint = 0
				endPoints = 1001
			} else {
				log10 := math.Log10(float64(startPoint))
				endPoints = startPoint + int(math.Pow10(int(log10)))
			}
			name := fmt.Sprintf(`%d-%d`, startPoint, endPoints)
			indexName = append(indexName, name)
		}
		data.Series = append(data.Series, Series{Data: valueList, Name: "用户人数"})
		data.Series = append(data.Series, Series{Data: GetPercentList(valueList), Name: "占比"})
		data.Title.Text = "用户积分"
	} else {
		var pointList []string
		for _, row := range rows[1:] {
			valueList = append(valueList, row.Str(0))
			startPoint := row.Int(1)
			pointList = append(pointList, fmt.Sprintf("%.0f", row.Float(2)))
			log10 := math.Log10(float64(startPoint))
			endPoints := startPoint + int(math.Pow10(int(log10)))
			name := fmt.Sprintf(`%d-%d`, startPoint, endPoints)
			indexName = append(indexName, name)
		}
		data.Series = append(data.Series, Series{Data: valueList, Name: "用户人数"})
		data.Series = append(data.Series, Series{Data: GetPercentList(pointList), Name: "积分占比"})
		data.Title.Text = "用户积分(折线是积分占比)"
	}
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "积分区间"
	data.Yaxis.Name = "人数"
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	if types != "-1" {
		redisConn.Do(`SET`, AnotherKey, bytes, `EX`, KeyAlive)
	} else {
		redisConn.Do(`SET`, pointDataKey, bytes, `EX`, KeyAlive)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
