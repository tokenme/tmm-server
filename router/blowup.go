package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/blowup"
)

func blowupRouter(r *gin.Engine) {
	blowupGroup := r.Group("/blowup")
	blowupGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		blowupGroup.GET("/notify", blowup.NotifyHandler)
		blowupGroup.POST("/bid", blowup.BidHandler)
		blowupGroup.POST("/escape", blowup.EscapeHandler)
	}
}
