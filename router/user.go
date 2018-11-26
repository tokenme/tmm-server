package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/user"
)

func userRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc())
	{
		userGroup.GET("/info", user.InfoGetHandler)
		userGroup.POST("/update", user.UpdateHandler)
		userGroup.GET("/invite/summary", user.InviteSummaryHandler)
		userGroup.GET("/credit/levels", user.CreditLevelsHandler)
	}
	r.POST("/user/create", handler.ApiSignFunc(), user.CreateHandler)
	r.POST("/user/reset-password", handler.ApiSignFunc(), user.ResetPasswordHandler)
	r.GET("/user/avatar/:key", user.AvatarGetHandler)
}
