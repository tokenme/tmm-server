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

const  inviteDataKey = `info-data-invite`

func InviteDataHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, inviteDataKey))
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
FROM (
    SELECT
        ic.parent_id, FLOOR((COUNT(1)-1)/5)*5+1 AS l
    FROM tmm.invite_codes AS ic
    WHERE ic.parent_id>0
    GROUP BY ic.parent_id
) AS tmp
GROUP BY l ORDER BY l`
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
		name := fmt.Sprintf(`%d-%d`, row.Int(1), row.Int(1)+5)
		indexName = append(indexName, name)
	}
	data := Data{
		Title:     "邀请人数占比",
		IndexName: indexName,
		Value:     valueList,
	}
	bytes, err := json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, inviteDataKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: data,
	})
}
