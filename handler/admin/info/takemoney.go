package info

import (
	. "github.com/tokenme/tmm/handler"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)

const tixianKey = `info-tixian`

func TixianInfoHandler(c *gin.Context) {
	redisConn := Service.Redis.Master.Get()

	context, err := redisConn.Do(`GET`, tixianKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := TixianInfo{}
		if json.Unmarshal(context.([]byte), &info) == nil {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    info,
			})
			return
		}
	}

	db := Service.Db

	query := `
	SELECT
          count(*) AS user,
		  SUM(total_cash) AS total,
          COUNT(IF(total_cash <= 10,total_cash,NULL)) AS a,
          COUNT(IF(10 < total_cash AND total_cash <= 100,total_cash,NULL)) AS b,
      	  COUNT(IF(100 < total_cash AND total_cash <= 1000,total_cash,NULL)) AS c,
  		  COUNT(IF(1000 < total_cash AND total_cash <= 10000,total_cash,NULL)) AS d,
		  COUNT(IF(10000 < total_cash ,total_cash,NULL)) AS e
 		  FROM tmm.top_withdraw_users
`
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	totalMoney, err := decimal.NewFromString(row.Str(res.Map(`total`)))
	if CheckErr(err, c) {
		return
	}
	info := TixianInfo{}
	info.TotalMoney = totalMoney
	info.TotalUser = row.Int(res.Map(`user`))
	info.AvgUserMoney = totalMoney.Div(decimal.NewFromFloat(float64(info.TotalUser)))
	info.LessTen = row.Int(res.Map(`a`))
	info.LessHundred = row.Int(res.Map(`b`))
	info.LessThousand = row.Int(res.Map(`c`))
	info.LessTenThousand = row.Int(res.Map(`d`))
	info.MoreThanTenThousand = row.Int(res.Map(`e`))

	query = `SELECT COUNT(1) FROM tmm.withdraw_txs `
	rows, _, err = db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}

	info.TotalTimes = rows[0].Int(0)

	query = `SELECT COUNT(1) FROM tmm.point_withdraws`
	rows, _, err = db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	info.TotalTimes = info.TotalTimes + rows[0].Int(0)
	info.AvgUserTimes = float64(info.TotalTimes) / float64(info.TotalUser)
	data, err := json.Marshal(&info)
	if CheckErr(err, c) {
		return
	}
	_, err = redisConn.Do(`SET`, tixianKey, data, "EX", KeyAlive)
	fmt.Println(err)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})

}
