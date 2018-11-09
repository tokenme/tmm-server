package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/good"
)

func goodRouter(r *gin.Engine) {
	goodGroup := r.Group("/good")
	goodGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		goodGroup.GET("/item/:id", good.ItemHandler)
		goodGroup.POST("/invest", good.InvestHandler)
		goodGroup.GET("/invests/item/:id/:page/:pageSize", good.ItemInvestsHandler)
	}
	r.GET("/good/list", good.ListHandler)
}
