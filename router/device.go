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
		deviceGroup.POST("/points", device.PointsHandler)
	}
	r.POST("/device/bind", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.BindHandler)
	r.POST("/device/push-token", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.PushTokenHandler)
	r.POST("/device/unbind", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.UnbindHandler)
	r.GET("/device/list", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.ListHandler)
	r.GET("/device/get/:deviceId", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.GetHandler)
	r.GET("/device/apps/:deviceId", AuthMiddleware.MiddlewareFunc(), handler.ApiSignFunc(), device.AppsHandler)
}
