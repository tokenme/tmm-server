package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func InvestSummaryHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT
        SUM(points), SUM(income)
FROM
    (SELECT good_id, SUM(points) AS points, SUM(income) AS income
    FROM (
        SELECT gi.good_id AS good_id, gi.points AS points, IFNULL(tx.income, 0) * gi.points/SUM(gi2.points) AS income
        FROM tmm.good_invests AS gi
        INNER JOIN tmm.good_txs AS tx ON (tx.good_id=gi.good_id AND tx.created_at>=gi.inserted_at)
        INNER JOIN tmm.good_invests AS gi2 ON (gi2.good_id=gi.good_id AND gi2.inserted_at<=tx.created_at)
        WHERE gi.user_id=%d
        GROUP BY tx.oid) AS tmp
    GROUP BY good_id
) AS tmp`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var summary common.GoodInvestSummary
	if len(rows) > 0 {
		row := rows[0]
		invest, _ := decimal.NewFromString(row.Str(0))
		income, _ := decimal.NewFromString(row.Str(1))
		summary = common.GoodInvestSummary{
			Invest: invest,
			Income: income,
		}
	}
	c.JSON(http.StatusOK, summary)
}
