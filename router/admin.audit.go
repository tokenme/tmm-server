package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/audit"
)

func AuditRouter(r *gin.Engine) {
	auditGroup := r.Group(`/admin/audit`)
	auditGroup.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		auditGroup.GET(`/app/list`, audit.AuditAppTaskListHandler)
		auditGroup.POST(`/app/edit`, audit.EditAppTaskHandler)
	}
}
