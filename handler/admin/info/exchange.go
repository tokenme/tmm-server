package info

import (
	"github.com/gin-gonic/gin"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	. "github.com/tokenme/tmm/handler"
)

const ExChangeKey = `info-exchange`

func ExchangeInfoHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	context, err := redisConn.Do(`Get`, ExChangeKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := ExChangeInfo{}
		if json.Unmarshal(context.([]byte), &info) == nil {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    info,
			})
			return
		}
	}
	query := `SELECT COUNT(IF(direction=1,tmm,NULL)) AS exchangeTmmTimes,
    COUNT(IF(direction=-1,tmm,NULL)) AS exchangePointTimes,
	COUNT(DISTINCT(user_id)) AS users,
    SUM(IF(direction=1, tmm, 0)) AS tmm,
    SUM(IF(direction=-1, points, 0)) AS point
	FROM tmm.exchange_records
 	WHERE status = 1 `
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]

	TotalTmm, err := decimal.NewFromString(row.Str(res.Map(`tmm`)))
	if CheckErr(err, c) {
		return
	}
	pointTotmm := Data{
		Total: TotalTmm,
		Times: row.Float(res.Map(`exchangeTmmTimes`)),
	}

	TotalPoint, err := decimal.NewFromString(row.Str(res.Map(`point`)))
	if CheckErr(err, c) {
		return
	}
	tmmToPoint := Data{
		Total: TotalPoint,
		Times: row.Float(res.Map(`exchangePointTimes`)),
	}

	Totaluser, err := decimal.NewFromString(row.Str(res.Map(`users`)))
	if CheckErr(err, c) {
		return
	}
	avgUser := Data{
		Total: TotalTmm.Div(Totaluser),
		Times: pointTotmm.Times / row.Float(res.Map(`users`)),
	}
	info := ExChangeInfo{}
	info.PointToTmm = pointTotmm
	info.TmmToPoint = tmmToPoint
	info.AvgUser = avgUser

	query = `SELECT COUNT(*),inserted_at FROM tmm.exchange_records  
	WHERE inserted_at < NOW() AND status = 1 GROUP BY date_format(inserted_at,'%Y-%m-%e')`
	rows, res, err = db.Query(query)

	if CheckErr(err, c) {
		return
	}

	sum := 0
	for _, row := range rows {
		sum += row.Int(0)
	}

	info.Daytimes = sum / int(len(rows))
	data, err := json.Marshal(info)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, ExChangeKey, data, "EX", KeyAlive)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})

}
