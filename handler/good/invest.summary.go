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
        SUM(points), SUM(bonus), SUM(income)
FROM
        (SELECT
                gi.points AS points,
                gi.bonus AS bonus,
                gi.points/t.points * t.income AS income
        FROM tmm.good_invests AS gi
        INNER JOIN (SELECT good_id, SUM(points) AS points, SUM(income) AS income FROM
(SELECT
        gi.good_id AS good_id,
        gi.points AS points,
        SUM(IFNULL(tx.income, 0)) AS income
FROM tmm.good_invests AS gi
LEFT JOIN tmm.good_txs AS tx ON (tx.good_id=gi.good_id AND tx.created_at>=gi.inserted_at)
WHERE
        EXISTS (SELECT 1 FROM tmm.good_invests AS gi2 WHERE gi2.good_id=gi.good_id AND gi2.user_id=%d AND gi2.redeem_status=0 AND gi2.inserted_at<=tx.created_at)
AND gi.redeem_status = 0
GROUP BY gi.good_id) AS tmp GROUP BY good_id) AS t ON (t.good_id = gi.good_id)
        WHERE
                gi.user_id=%d
        AND gi.redeem_status = 0
        GROUP BY gi.good_id) AS tmp`, user.Id, user.Id)
	if CheckErr(err, c) {
		return
	}
	var summary common.GoodInvestSummary
	if len(rows) > 0 {
		row := rows[0]
		invest, _ := decimal.NewFromString(row.Str(0))
		bonus, _ := decimal.NewFromString(row.Str(1))
		income, _ := decimal.NewFromString(row.Str(2))
		summary = common.GoodInvestSummary{
			Invest: invest,
			Bonus:  bonus,
			Income: income,
		}
	}
	c.JSON(http.StatusOK, summary)
}
