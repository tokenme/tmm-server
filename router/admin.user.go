package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/user"
)

func UserRouter(r *gin.Engine){
	userR:=r.Group(`/admin/user`)
	userR.Use(AuthMiddleware.MiddlewareFunc(),verify.VerifyAdminFunc())
	{
		userR.GET(`/info`,user.GetUserInfoHandler)
	}
}