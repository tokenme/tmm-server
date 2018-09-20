package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/device"
)

func deviceRouter(r *gin.Engine) {
	deviceGroup := r.Group("/device")
	deviceGroup.Use(handler.ApiCheckFunc())
	{
		deviceGroup.POST("/ping", device.PingHandler)
		deviceGroup.POST("/save", device.SaveHandler)
	}
	r.POST("/device/bind", AuthMiddleware.MiddlewareFunc(), device.BindHandler)
	r.POST("/device/unbind", AuthMiddleware.MiddlewareFunc(), device.UnbindHandler)
	r.GET("/device/list", AuthMiddleware.MiddlewareFunc(), device.ListHandler)
	r.GET("/device/get/:deviceId", AuthMiddleware.MiddlewareFunc(), device.GetHandler)
	r.GET("/device/apps/:deviceId", AuthMiddleware.MiddlewareFunc(), device.AppsHandler)
}
