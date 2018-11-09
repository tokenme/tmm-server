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
	}
	r.GET("/good/list", good.ListHandler)
}
