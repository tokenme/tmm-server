package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/user"
	"github.com/tokenme/tmm/handler/admin/verify"
)

func UserRouter(r *gin.Engine) {
	userR := r.Group(`/admin/user`)
	userR.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		userR.GET(`/info`, user.GetUserInfoHandler)
	}
}
