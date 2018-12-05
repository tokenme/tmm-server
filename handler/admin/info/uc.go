package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
	"time"
)

const UcInfoKey = `Info-Ucoin`

func UcInfoHandler(c *gin.Context) {
	redisConn := Service.Redis.Master.Get()
	info := UcInfo{}
	Context, err := redisConn.Do(`Get`, UcInfoKey)
	if err == nil {
		json.Unmarshal(Context.([]byte), &info)
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Data:    info,
			Message: admin.API_OK,
		})
		return
	}
	db := Service.Db
	query := `SELECT 
 SUM(IF(direction=1, tmm, 0)) AS supply,
 SUM(IF(direction=-1, tmm, 0)) AS burn,
 SUM(IF(DATE_ADD(inserted_at,INTERVAL 1 DAY) > NOW(),tmm,0)) AS d,
 COUNT(DISTINCT(user_id))
 FROM tmm.exchange_records
 WHERE status in (2,1)`

	row, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(row) == 0, `Not Find`, c) {
		return
	}
	totalSupply, err := decimal.NewFromString(row[0].Str(0))
	if CheckErr(err, c) {
		return
	}
	totalBurn, err := decimal.NewFromString(row[0].Str(1))
	if CheckErr(err, c) {
		return
	}

	daySupply, err := decimal.NewFromString(row[0].Str(2))
	if CheckErr(err, c) {
		return
	}
	user, err := decimal.NewFromString(row[0].Str(3))
	if CheckErr(err, c) {
		return
	}
	info.TotalSupply = totalSupply
	info.Totalburn = totalBurn
	info.CurrentSupply = info.Totalburn.Sub(info.Totalburn)
	info.DaySupply = daySupply
	info.AvgPersonSupply = info.TotalSupply.Div(user)
	redisConn.Do(`SETEX`,UcInfoKey,60*60*24,time.Now().Format("2006-01-02 15:04:05"))
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Data:    info,
		Message: admin.API_OK,
	})
}
