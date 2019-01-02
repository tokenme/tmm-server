package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"net/http"
)

func WithdrawHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	pointsQuery := `SELECT SUM(pw.cny) FROM point_withdraws AS pw WHERE pw.user_id=%d GROUP BY pw.user_id`
	pointsRows, _, err := db.Query(pointsQuery, user.Id)
	if CheckErr(err, c) {
		return
	}
	var userWithdraw common.UserWithdraw
	if len(pointsRows) > 0 {
		userWithdraw.Points, _ = decimal.NewFromString(pointsRows[0].Str(0))
	}
	tmmQuery := `SELECT SUM(wt.cny) FROM withdraw_txs AS wt WHERE wt.user_id=%d GROUP BY wt.user_id`
	tmmRows, _, err := db.Query(tmmQuery, user.Id)
	if CheckErr(err, c) {
		return
	}
	if len(tmmRows) > 0 {
		userWithdraw.TMM, _ = decimal.NewFromString(tmmRows[0].Str(0))
	}
	currency := c.Query("currency")
	if currency == "" {
		currency = "CNY"
	}
	if currency != "CNY" {
		rate := forex.Rate(Service, "CNY", currency)
		userWithdraw.Points = userWithdraw.Points.Mul(rate)
		userWithdraw.TMM = userWithdraw.TMM.Mul(rate)
	}
	c.JSON(http.StatusOK, userWithdraw)
}
