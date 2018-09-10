package exchange

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func TMMRateHandler(c *gin.Context) {
	exchangeRate, err := common.GetExchangeRate(Config, Service)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, exchangeRate)
}
