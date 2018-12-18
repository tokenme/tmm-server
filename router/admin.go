package router

import (
	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.Engine) {
	r.POST(`/admin/auth/login`,AuthMiddleware.LoginHandler)
	AricleRouter(r)
	TaskRouter(r)
	UserRouter(r)
	InfoRouter(r)
	AdRouter(r)
}
