package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
)

const pointDataKey = `info-Data-points`

func PointDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	Context, err := redisConn.Do(`GET`, pointDataKey)
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
    COUNT(*) AS users,
    l
FROM (
    SELECT
	 	d.user_id,
	 	IF(SUM(d.points)=0,0,FLOOR(((SUM(d.points)-1)/50))*50+1) AS l,
	 	d.points
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
		Name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+50)
		indexName = append(indexName, Name)
	}
	data := Data{
		Title:     "用户积分占比",
		IndexName: indexName,
		Value:     valueList,
	}
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
