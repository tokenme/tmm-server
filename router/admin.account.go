package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/account_controller"
)

func AccountRouter(r *gin.Engine) {
	AccountGroup := r.Group(`admin/account`)
	AccountGroup.Use(AuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		AccountGroup.GET(`/list`, account_controller.GetAccountList)
		AccountGroup.PUT(`/user`, account_controller.EditAccountHandler)
		AccountGroup.GET(`/user`, account_controller.UserInfoHandler)
		AccountGroup.GET(`/user/friend`, account_controller.FriendsHandler)
		AccountGroup.GET(`/user/point`, account_controller.MakePointHandler)
	}
	{
		AccountGroup.GET(`/user/exchange/uc`, account_controller.ExchangeByUcHandler)
		AccountGroup.GET(`/user/exchange/point`, account_controller.ExchangeByPointHandler)
	}
	{
		AccountGroup.GET(`/user/drawmoney/point`, account_controller.DrawMoneyByPointHandler)
		AccountGroup.GET(`/user/drawmoney/uc`, account_controller.DrawMoneyByUcHandler)
	}
}
