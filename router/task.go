package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/task"
)

func taskRouter(r *gin.Engine) {
	taskGroup := r.Group("/task")
	taskGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		taskGroup.GET("/shares", task.SharesHandler)
		taskGroup.GET("/apps", task.AppsHandler)
		taskGroup.POST("/app/install", task.AppInstallHandler)
		taskGroup.GET("/apps/check", task.AppsCheckHandler)
		taskGroup.GET("/records", task.RecordsHandler)
	}
	r.GET("/share/:encryptedTaskId/:encryptedDeviceId", task.ShareHandler)
}
