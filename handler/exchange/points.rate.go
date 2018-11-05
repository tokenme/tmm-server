package exchange

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	//"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	"net/http"
)

func PointsRateHandler(c *gin.Context) {
	points, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"rate": points})
}
