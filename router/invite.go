package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/invite"
)

func inviteRouter(r *gin.Engine) {
	inviteGroup := r.Group("/invite")
	inviteGroup.GET("/:code", invite.ShowHandler)
	inviteGroup.POST("/:code", invite.SubmitHandler)
}
