package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/bonus"
)

func bonusRouter(r *gin.Engine) {
	bonusGroup := r.Group("/bonus")
	bonusGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		bonusGroup.GET("/daily/status", handler.ApiSignPassFunc(), bonus.DailyStatusHandler)
		bonusGroup.POST("/daily/commit", handler.ApiSignFunc(), bonus.DailyCommitHandler)
		bonusGroup.POST("/reading", handler.ApiSignFunc(), bonus.ReadingHandler)
	}
}
