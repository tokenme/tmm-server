package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/task"
)

func taskRouter(r *gin.Engine) {
	taskGroup := r.Group("/task")
	taskGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		taskGroup.GET("/shares", handler.ApiSignPassFunc(), task.SharesHandler)
		taskGroup.GET("/apps", task.AppsHandler)
		taskGroup.POST("/app/install", handler.ApiSignFunc(), task.AppInstallHandler)
        taskGroup.POST("/app/certificates/upload", handler.ApiSignFunc(), task.AppCertificatesUploadHandler)
		taskGroup.GET("/apps/check", task.AppsCheckHandler)
		taskGroup.GET("/records", task.RecordsHandler)
		taskGroup.POST("/share/add", handler.ApiSignFunc(), task.ShareAddHandler)
		taskGroup.POST("/share/update", handler.ApiSignFunc(), task.ShareUpdateHandler)
		taskGroup.POST("/app/add", handler.ApiSignFunc(), task.AppAddHandler)
	}
	r.GET("/share/:encryptedTaskId/:encryptedDeviceId", task.ShareHandler)
	r.GET("/s/track/:encryptedTaskId", task.ShareTrackHandler)
	r.GET("/s/imp/:encryptedTaskId/:encryptedDeviceId", task.ShareImpHandler)
}
