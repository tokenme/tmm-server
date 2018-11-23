package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/feedback"
)

func feedbackRouter(r *gin.Engine) {
	feedbackGroup := r.Group("/feedback")
	feedbackGroup.Use(AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc())
	{
		feedbackGroup.POST("/add", feedback.AddHandler)
		feedbackGroup.GET("/list", feedback.ListHandler)
		feedbackGroup.POST("/reply", feedback.ReplyHandler)
	}
}
