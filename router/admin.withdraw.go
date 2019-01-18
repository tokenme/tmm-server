package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/withdraw"
)

func WithdrawRouter(r *gin.Engine) {
	AccountGroup := r.Group(`admin/withdraw`)
	AccountGroup.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		AccountGroup.GET(`/list`, withdraw.GetWithDrawList)
		AccountGroup.POST(`/edit`, withdraw.EditWithDrawHandler)
	}
}
