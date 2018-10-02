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
		orderbookGroup.POST("/order/cancel", orderbook.OrderCancelHandler)
		orderbookGroup.GET("/market/top/:side", orderbook.MarketTopHandler)
		orderbookGroup.GET("/rate", orderbook.RateHandler)
		orderbookGroup.GET("/orders/:page/:pageSize/:side", orderbook.OrdersHandler)
	}
	r.GET("/orderbook/market/graph/:hours", orderbook.MarketGraphHandler)
}
