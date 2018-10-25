package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/bonus"
)

func bonusRouter(r *gin.Engine) {
	bonusGroup := r.Group("/bonus")
	bonusGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		bonusGroup.GET("/daily/status", bonus.DailyStatusHandler)
		bonusGroup.POST("/daily/commit", bonus.DailyCommitHandler)
	}
}