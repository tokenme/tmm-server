package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/user"
)

func userRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		userGroup.GET("/info", user.InfoGetHandler)
		userGroup.GET("/balance", user.BalanceHandler)
		userGroup.POST("/update", handler.ApiSignFunc(), user.UpdateHandler)
		userGroup.GET("/invite/summary", user.InviteSummaryHandler)
		userGroup.GET("/invite/lastday-contribute", user.InviteLastdayContributeHandler)
		userGroup.GET("/invites", user.InviteListHandler)
		userGroup.GET("/credit/levels", user.CreditLevelsHandler)
        userGroup.GET("/withdraw", user.WithdrawHandler)
	}
	r.POST("/user/create", handler.ApiSignFunc(), user.CreateHandler)
	r.POST("/user/reset-password", handler.ApiSignFunc(), user.ResetPasswordHandler)
	r.GET("/user/avatar/:key", user.AvatarGetHandler)
}
