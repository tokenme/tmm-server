package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/exchange"
)

func exchangeRouter(r *gin.Engine) {
	exchangeGroup := r.Group("/exchange")
	exchangeGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		exchangeGroup.POST("/tmm/change", handler.ApiSignPassFunc(), exchange.TMMChangeHandler)
		exchangeGroup.GET("/records", handler.ApiSignFunc(), exchange.RecordsHandler)
	}
	r.GET("/exchange/tmm/rate", exchange.TMMRateHandler)
	r.GET("/exchange/points/rate", exchange.PointsRateHandler)
}
