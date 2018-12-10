package info

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/garyburd/redigo/redis"
)

const totalTaskCashKey = `info-total-task`

func TotalTaskHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, totalTaskCashKey))

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
	SUM(points) AS cost,
	SUM(total) AS total 
FROM(
SELECT 
	SUM(dev.points) AS points ,
	COUNT(1)   AS total 
FROM 
	tmm.device_share_tasks  AS dev
	UNION ALL
SELECT 

	SUM(app.points) AS points ,
	COUNT(1)   AS total
FROM 
	tmm.device_app_tasks AS app
WHERE 
	app.status = 1
) AS tmp
`
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	var total TotalTask
	totalCost, err := decimal.NewFromString(row.Str(res.Map(`cost`)))
	if CheckErr(err, c) {
		return
	}
	total.TotalCost = totalCost
	total.TotaltaskCount = row.Int(res.Map(`total`))
	bytes, err := json.Marshal(&total)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, totalTaskCashKey, bytes, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    total,
	})

}
