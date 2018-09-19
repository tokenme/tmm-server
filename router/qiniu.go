package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/qiniu"
)

func qiniuRouter(r *gin.Engine) {
	qiniuGroup := r.Group("/qiniu")
	qiniuGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		qiniuGroup.GET("/uptoken", qiniu.UpTokenHandler)
	}
}
