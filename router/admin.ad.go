package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/ad"
)

func AdRouter(r *gin.Engine) {
	adGroup := r.Group(`admin/ad`)
	adGroup.Use(AuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		adGroup.GET(`/token`, ad.GetUpdateToken)
		adGroup.POST(`/creative/add`, ad.AddCreativeHandler)
		adGroup.GET(`/creative/list`, ad.CreativeListHandler)
		adGroup.GET(`/creative/edit`, ad.EditCreativeHanlder)
	}
	{
		adGroup.POST(`/group/add`, ad.AddGroupsHandler)
		adGroup.GET(`/group/list`, ad.GroupListHanlder)
		adGroup.GET(`/group/edit`, ad.EditGroupHanler)
	}
	{
		adGroup.GET(`/adzone/list`, ad.AdzoneListHandler)
	}
}
