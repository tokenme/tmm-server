package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/ad"
)

func adRouter(r *gin.Engine) {
	adGroup := r.Group("/ad")
	adGroup.GET("/imp/:code", ad.ImpHandler)
	adGroup.GET("/clk/:code", ad.ClkHandler)
}
