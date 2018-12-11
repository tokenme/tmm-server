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

const taskDataKey = `info-data-draw`

func TaskDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, taskDataKey))
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
	query := `SELECT
    COUNT(*) AS users,
    l
FROM (   SELECT 
 tmp.device_id, 
 IF(SUM(tmp.total)=0, 0,FLOOR((SUM(tmp.total)-1)/100) * 100 + 1) AS l,
 SUM(tmp.total)
FROM(
	SELECT 
 		dev.device_id,
		COUNT(1)  AS total 
	FROM 
		tmm.device_share_tasks  AS dev
	GROUP BY 
		dev.device_id UNION ALL
	SELECT 
		app.device_id,
		COUNT(1)   AS total
	FROM 
		tmm.device_app_tasks AS app
	WHERE 
		app.status = 1
	GROUP BY 
		app.device_id
) AS tmp,tmm.user_devices AS ud
WHERE ud.device_id = tmp.device_id
GROUP BY ud.user_id )AS tmp
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
		Title:     "任务完成占比",
		IndexName: indexName,
		Value:     valueList,
	}
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, taskDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
