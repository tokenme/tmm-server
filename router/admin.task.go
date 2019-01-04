package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/task"
	"github.com/tokenme/tmm/handler/qiniu"
	"github.com/tokenme/tmm/handler/admin/app"
)
func TaskRouter(r *gin.Engine){
	 taskR:=r.Group(`/admin/share`)
	 taskR.Use(AuthMiddleware.MiddlewareFunc(),verify.VerifyAdminFunc())
	 {
	 	taskR.GET(`/getToken`,qiniu.UpTokenHandler)
	 	taskR.POST(`/add`,task.AddShareHandler)
	 	taskR.POST(`/modify`,task.ModifyTaskHandler)
	 	taskR.GET(`/list`,task.GetTaskListHandler)
	 	taskR.GET(`/edit`,task.GetTaskHandler)
		taskR.POST(`/add-app`,app.AddShareAppHandler)
		taskR.GET(`/list-app`,app.GetShareAppHandler)
		taskR.POST(`/modify-app`,app.ModifyShareAppHandler)
		taskR.GET(`/get-app`, app.GetAppTaskHandler)
	 }
}

