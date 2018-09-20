package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/user"
)

func userRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		userGroup.GET("/info", user.InfoGetHandler)
		userGroup.POST("/update", user.UpdateHandler)
		userGroup.GET("/invite/summary", user.InviteSummaryHandler)
	}
	r.POST("/user/create", user.CreateHandler)
	r.POST("/user/reset-password", user.ResetPasswordHandler)
	r.GET("/user/avatar/:key", user.AvatarGetHandler)
}
