package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/redeem"
)

func redeemRouter(r *gin.Engine) {
	redeemGroup := r.Group("/redeem")
	redeemGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		redeemGroup.GET("/dycdp/list", redeem.DycdpListHandler)
		redeemGroup.POST("/dycdp/order/add", redeem.DycdpOrderAddHandler)
		redeemGroup.POST("/tmm/withdraw", redeem.TMMWithdrawHandler)
	}
	r.GET("/redeem/tmm/rate", redeem.TMMRateHandler)
}
