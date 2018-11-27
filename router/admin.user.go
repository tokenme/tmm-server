package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/Verify"
	"github.com/tokenme/tmm/handler/admin/user"
)

func UserRouter(r *gin.Engine){
	userR:=r.Group(`/admin/user`)
	userR.Use(AuthMiddleware.MiddlewareFunc(),Verify.VerifyAdminFunc())
	{
		userR.GET(`/info`,user.GetUserInfoHandler)
	}
}