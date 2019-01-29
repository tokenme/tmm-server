package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/article"
)

func articleRouter(r *gin.Engine) {
	articleGroup := r.Group("/article")
	articleGroup.GET("/show/:id", article.ShowHandler)
	articleGroup.GET("/rand", article.RandHandler)
}
