package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
)

const InvestsKey = `info-invests`

func InvestsInfoHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	context, err := redisConn.Do(`GET`, InvestsKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := InvestsInfo{}
		if json.Unmarshal(context.([]byte), &info) == nil {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    info,
			})
			return
		}
	}
	query := `SELECT 
	SUM(points) AS total_points,
	COUNT(*) as total_goods 
	FROM tmm.good_invests AS i 
	INNER JOIN goods AS g ON g.id = i.good_id    
	WHERE i.redeem_status != 2 `
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	totalPoint, err := decimal.NewFromString(row.Str(0))
	if CheckErr(err, c) {
		return
	}
	info := InvestsInfo{}
	info.TotalPoint = totalPoint
	info.TotalGoods = row.Int(1)
	info.AvgGoodsInvestsPoint = totalPoint.Div(decimal.NewFromFloat(row.Float(1)))

	data, err := json.Marshal(info)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do("SET", InvestsKey, data, "EX", KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
