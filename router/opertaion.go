package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/operation"
)

func operationRouter(r *gin.Engine) {
	operationGroup := r.Group("/operation")
    operationGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
        operationGroup.GET("/headlines", operation.HeadlinesHandler)
        operationGroup.GET("/activities", operation.ActivitiesHandler)
    }
}
