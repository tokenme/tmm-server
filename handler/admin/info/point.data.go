package info

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"math"
	"net/http"
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
        IF(points  >= 10000,
            IF(points < 100000,
                FLOOR(points/10000)*10000+1,
                IF(points < 1000000,
                    FLOOR(points/100000)*100000+1,
                    FLOOR(points/1000000)*1000000+1
                )
            ),
            FLOOR(IF(points=0, 1, points) /1000)*1000+1
        ) AS l,
        COUNT(1) AS users,
        points
    FROM
    (
        SELECT
            d.user_id,
            SUM(d.points) AS points
        FROM tmm.devices AS d
        WHERE NOT EXISTS (SELECT 1 FROM tmm.user_settings AS us WHERE us.blocked=1 AND us.user_id=d.user_id AND us.block_whitelist=0 LIMIT 1)
        GROUP BY d.user_id
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
	var valueList []string
	var data Data
	var pointList []string
	for _, row := range rows {
		valueList = append(valueList, row.Str(1))
		pointList = append(pointList, fmt.Sprintf("%.0f", row.Float(2)))
		startPoint := row.Int(0)
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
	data.Series = append(data.Series, Series{Data: GetPercentList(valueList), Name: "人数占比"})
	data.Series = append(data.Series, Series{Data: GetPercentList(pointList), Name: "积分占比"})
	data.Title.Text = "用户积分"
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "积分区间"
	data.Yaxis.Name = "人数"
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
