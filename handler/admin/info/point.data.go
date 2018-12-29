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
    l,
	SUM(tmp.total_points)
FROM (
	SELECT
		d.user_id,
	 	IF (d.total_points=0,0,
		IF(d.total_points  >= 10000,
		FLOOR(d.total_points/10000)*10000+1,
		FLOOR(d.total_points/1000)*1000+1) ) AS l,
		d.total_points
	FROM (
	SELECT 
		dev.user_id AS user_id ,
		IFNULL(reading.points,0)+IFNULL(inv.bonus,0) + IFNULL(task.points,0) AS total_points
	FROM 
		tmm.devices AS dev
	LEFT JOIN (
		SELECT 
			SUM(inv.bonus) AS bonus,						
			inv.user_id
		FROM 
			tmm.invite_bonus AS inv 
		GROUP BY 
			inv.user_id 
		) AS inv  ON (inv.user_id = dev.user_id)
	LEFT JOIN (
		SELECT 
			SUM(tmp.points) AS points,
			dev.user_id AS user_id 
		FROM (
			SELECT 
				sha.device_id, 
				SUM(sha.points) AS points
			FROM 
				tmm.device_share_tasks AS sha	
			WHERE 
				sha.points > 0 
			GROUP BY
				sha.device_id UNION ALL
			SELECT 
				app.device_id, 
				SUM(app.points) AS points
			FROM 
				tmm.device_app_tasks AS app   
			WHERE
				app.status = 1
			GROUP BY
				app.device_id   
		) AS tmp 
		INNER JOIN tmm.devices  AS dev ON (dev.id = tmp.device_id)
		GROUP BY 
			dev.user_id 
	) AS task ON (task.user_id = dev.user_id)
	LEFT JOIN (
		SELECT 
			SUM(point) AS points,
			user_id 
		FROM 
			tmm.reading_logs 
		GROUP BY 
			user_id 
	) AS reading  ON (reading.user_id =dev.user_id)
	GROUP BY 
		dev.user_id 
	)AS d 
WHERE NOT EXISTS
	(SELECT 1 FROM tmm.user_settings AS us  
	WHERE us.blocked= 1 AND us.user_id=d.user_id AND us.block_whitelist=0  LIMIT 1)
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
	var pointList []string
	for _, row := range rows {
		valueList = append(valueList, row.Str(0))
		pointList = append(pointList, fmt.Sprintf("%.0f", row.Float(2)))
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
