package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/exchange"
)

func exchangeRouter(r *gin.Engine) {
	exchangeGroup := r.Group("/exchange")
	exchangeGroup.Use(AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc())
	{
		exchangeGroup.POST("/tmm/change", exchange.TMMChangeHandler)
		exchangeGroup.GET("/records", exchange.RecordsHandler)
	}
	r.GET("/exchange/tmm/rate", handler.ApiSignFunc(), exchange.TMMRateHandler)
	r.GET("/exchange/points/rate", handler.ApiSignFunc(), exchange.PointsRateHandler)
}
