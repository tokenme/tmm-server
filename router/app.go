package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/app"
)

func appRouter(r *gin.Engine) {
	appGroup := r.Group("/app")
	appGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		appGroup.GET("/sdks/:platform/:page/:pageSize", app.SdksHandler)
	}
}
