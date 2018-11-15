package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ykt"
	"net/http"
	"strconv"
)

func ItemHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	itemId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	yktReq := ykt.GoodInfoRequest{
		Id:  itemId,
		Uid: user.Id,
	}
	res, err := yktReq.Run()
	if CheckErr(err, c) {
		return
	}
	good := res.Data.Data
	good.CommissionPoints = decimal.New(Config.GoodCommissionPoints, 0)
	db := Service.Db
	rows, _, err := db.Query(`SELECT points FROM tmm.good_invests WHERE good_id=%d AND user_id=%d AND gi.redeem_status=0 LIMIT 1`, good.Id, user.Id)
	if err == nil && len(rows) > 0 {
		row := rows[0]
		points, _ := decimal.NewFromString(row.Str(0))
		good.InvestPoints = points
	} else if err != nil {
		log.Error(err.Error())
	}
	rows, _, err = db.Query(`SELECT SUM(points), COUNT(*), SUM(income) FROM
            (SELECT
                gi.points AS points,
                gi.bonus AS bonus,
                SUM(IFNULL(tx.income, 0)) AS income
            FROM tmm.good_invests AS gi
            LEFT JOIN tmm.good_txs AS tx ON (tx.good_id=gi.good_id AND tx.created_at>=gi.inserted_at)
            WHERE
                gi.good_id=%d
            AND gi.redeem_status = 0
            GROUP BY gi.good_id) AS tmp`, good.Id)
	if err == nil && len(rows) > 0 {
		row := rows[0]
		points, _ := decimal.NewFromString(row.Str(0))
		income, _ := decimal.NewFromString(row.Str(2))
		good.TotalInvest = points
		good.TotalInvestors = row.Uint(1)
		good.InvestIncome = income
	} else if err != nil {
		log.Error(err.Error())
	}
	c.JSON(http.StatusOK, good)
}
