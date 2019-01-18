package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
)

func AdminRouter(r *gin.Engine) {
	r.POST(`/admin/auth/login`,AuthMiddleware.LoginHandler,verify.RecordLoginHandler)
	AricleRouter(r)
	TaskRouter(r)
	UserRouter(r)
	InfoRouter(r)
	AdRouter(r)
	AccountRouter(r)
	WithdrawRouter(r)
	AuditRouter(r)
}
