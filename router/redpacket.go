package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/redpacket"
)

func redpacketRouter(r *gin.Engine) {
	redpacketGroup := r.Group("/redpacket")
	redpacketGroup.GET("/show/:encryptedId/:encryptedUserId", redpacket.ShowHandler)
	redpacketGroup.POST("/submit", redpacket.SubmitHandler)
}
