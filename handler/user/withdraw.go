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
	query := `
        SELECT SUM(pw.cny) FROM point_withdraws AS pw WHERE pw.user_id=%d GROUP BY pw.user_id
        UNION ALL
        SELECT SUM(wt.cny) FROM withdraw_txs AS wt WHERE wt.user_id=%d GROUP BY wt.user_id
    `
	rows, _, err := db.Query(query, user.Id, user.Id)
	if CheckErr(err, c) {
		return
	}
	var userWithdraw common.UserWithdraw
	if len(rows) > 0 {
		userWithdraw.Points, _ = decimal.NewFromString(rows[0].Str(0))
		if len(rows) > 1 {
			userWithdraw.TMM, _ = decimal.NewFromString(rows[1].Str(0))
		}
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
