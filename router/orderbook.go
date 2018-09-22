package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/orderbook"
)

func orderbookRouter(r *gin.Engine) {
	orderbookGroup := r.Group("/orderbook")
	orderbookGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		orderbookGroup.POST("/order/add", orderbook.OrderAddHandler)
		orderbookGroup.GET("/market/top/:side", orderbook.MarketTopHandler)
	}
}
