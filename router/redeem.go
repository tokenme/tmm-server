package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/redeem"
)

func redeemRouter(r *gin.Engine) {
	redeemGroup := r.Group("/redeem")
	redeemGroup.Use(AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc())
	{
		redeemGroup.GET("/dycdp/list", redeem.DycdpListHandler)
		redeemGroup.POST("/dycdp/order/add", redeem.DycdpOrderAddHandler)
		redeemGroup.POST("/tmm/withdraw", redeem.TMMWithdrawHandler)
		redeemGroup.GET("/tmm/withdraw/list", redeem.TMMWithdrawListHandler)
	}
	r.GET("/redeem/tmm/rate", handler.ApiSignFunc(), redeem.TMMRateHandler)
}
