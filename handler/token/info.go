package token

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func InfoHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	tokenAddress := c.Param("address")
	var token common.Token
	if tokenAddress == "0x" {
		ethBalance, _ := eth.BalanceOf(Service.Geth, c, user.Wallet)
		token = common.Token{
			Name:     "Ethereum",
			Symbol:   "ETH",
			Decimals: 18,
			Icon:     "https://www.ethereum.org/images/logos/ETHEREUM-ICON_Black_small.png",
			Balance:  decimal.NewFromBigInt(ethBalance, -18),
		}
	} else {
		tokenABI, err := utils.NewToken(tokenAddress, Service.Geth)
		if CheckErr(err, c) {
			return
		}
		name, err := tokenABI.Name(nil)
		if CheckErr(err, c) {
			return
		}
		symbol, err := tokenABI.Symbol(nil)
		if CheckErr(err, c) {
			return
		}
		decimals, err := tokenABI.Decimals(nil)
		if CheckErr(err, c) {
			return
		}
		balance, err := utils.TokenBalanceOf(tokenABI, user.Wallet)
		if CheckErr(err, c) {
			return
		}
		balanceDecimal, err := decimal.NewFromString(balance.String())
		if CheckErr(err, c) {
			return
		}
		token = common.Token{
			Name:     name,
			Symbol:   symbol,
			Decimals: uint(decimals),
		}

		if token.Decimals > 0 {
			token.Balance = balanceDecimal.Div(decimal.New(1, int32(token.Decimals)))
		} else {
			token.Balance = balanceDecimal
		}

		if token.Decimals > 0 {
			token.Balance = balanceDecimal.Div(decimal.New(1, int32(token.Decimals)))
		} else {
			token.Balance = balanceDecimal
		}
	}

	c.JSON(http.StatusOK, token)
}
