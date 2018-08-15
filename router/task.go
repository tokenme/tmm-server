package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/task"
)

func taskRouter(r *gin.Engine) {
	taskGroup := r.Group("/task")
	taskGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		taskGroup.GET("/shares/:page/:pageSize", task.SharesHandler)
	}
	r.GET("/share/:encryptedTaskId/:encryptedUserId", task.ShareHandler)
}
