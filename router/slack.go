package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/slack"
)

func slackRouter(r *gin.Engine) {
	slackGroup := r.Group("/slack")
	slackGroup.POST("/hook", slack.HookHandler)
}
