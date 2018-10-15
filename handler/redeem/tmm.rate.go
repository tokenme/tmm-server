package redeem

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"net/http"
)

func TMMRateHandler(c *gin.Context) {
	currency := c.Query("currency")
	if currency == "" {
		currency = "USD"
	}
	price := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	if currency != "USD" {
		rate := forex.Rate(Service, "USD", currency)
		price = price.Mul(rate)
	}
	c.JSON(http.StatusOK, common.ExchangeRate{
		Rate:      price,
		MinPoints: decimal.New(int64(Config.MinTMMRedeem), 0),
	})
}
