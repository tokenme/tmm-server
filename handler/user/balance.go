package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"net/http"
)

func BalanceHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	query := `SELECT SUM(d.points) FROM devices AS d WHERE d.user_id=%d GROUP BY d.user_id`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}
	var userBalance common.UserBalance
	if len(rows) > 0 {
		userBalance.Points, _ = decimal.NewFromString(rows[0].Str(0))
	}
	tokenABI, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	_, _, decimals, _, _, _, _, _, balance, err := utils.TokenMeta(tokenABI, user.Wallet)
	balanceDecimal, err := decimal.NewFromString(balance.String())
	if CheckErr(err, c) {
		return
	}
	userBalance.TMM = balanceDecimal.Div(decimal.New(1, int32(decimals)))
	currency := c.Query("currency")
	if currency == "" {
		currency = "USD"
	}
	pointPrice := common.GetPointPrice(Service, Config)
	tmmPrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	if currency != "USD" {
		rate := forex.Rate(Service, "USD", currency)
		pointPrice = pointPrice.Mul(rate)
		tmmPrice = tmmPrice.Mul(rate)
	}
	userBalance.Cash = pointPrice.Mul(userBalance.Points).Add(tmmPrice.Mul(userBalance.TMM))
	c.JSON(http.StatusOK, userBalance)
}
