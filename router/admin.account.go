package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/account_controller"
	"github.com/tokenme/tmm/handler/admin/verify"
)

func AccountRouter(r *gin.Engine) {
	AccountGroup := r.Group(`admin/account`)
	AccountGroup.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		AccountGroup.GET(`/list`, account_controller.GetAccountList)
		AccountGroup.POST(`/user`, account_controller.EditAccountHandler)
		AccountGroup.GET(`/user`, account_controller.UserInfoHandler)
		AccountGroup.GET(`/user/friend`, account_controller.FriendsHandler)
		AccountGroup.GET(`/user/point`, account_controller.MakePointHandler)
	}
	{
		AccountGroup.GET(`/user/exchange`, account_controller.ExchangeHandler)
	}
	{
		AccountGroup.GET(`/user/drawmoney/point`, account_controller.DrawMoneyByPointHandler)
		AccountGroup.GET(`/user/drawmoney/uc`, account_controller.DrawMoneyByUcHandler)
	}
}
