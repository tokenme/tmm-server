package token

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func TMMBalanceHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	tokenABI, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	name, symbol, decimals, _, _, _, _, _, _, err := utils.TokenMeta(tokenABI, user.Wallet)
	balance, err := utils.TokenBalanceOf(tokenABI, user.Wallet)
	if CheckErr(err, c) {
		return
	}
	balanceDecimal, err := decimal.NewFromString(balance.String())
	if CheckErr(err, c) {
		return
	}
	token := common.Token{
		Name:     name,
		Symbol:   symbol,
		Decimals: uint(decimals),
		Balance:  balanceDecimal.Div(decimal.New(1, int32(decimals))),
	}
	c.JSON(http.StatusOK, token)
}
