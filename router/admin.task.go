package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/Verify"
	"github.com/tokenme/tmm/handler/admin/task"
	"github.com/tokenme/tmm/handler/qiniu"
)
func TaskRouter(r *gin.Engine){
	 taskR:=r.Group(`/admin/share`)
	 taskR.Use(AuthMiddleware.MiddlewareFunc(),Verify.VerifyAdminFunc())
	 {
	 	taskR.GET(`/getToken`,qiniu.UpTokenHandler)
	 	taskR.POST(`/add`,task.AddShareHandler)
	 	taskR.POST(`/exit`,task.ExitTaskHandler)
	 	taskR.GET(`/list`,task.GetShareListHandler)
	 	taskR.POST(`/online`,task.TaskUpdateHandler)
	 }
	 }

