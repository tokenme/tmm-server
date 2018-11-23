package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

func AdminRouter(r *gin.Engine) {
	AdminGroup := r.Group("/admin")
	AdminAuth := AuthMiddleware
	AdminAuth.Unauthorized = func(c *gin.Context, code int, message string) {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": message,
		})
	}
	r.POST(`/admin/login`, AdminAuth.LoginHandler)
	AdminGroup.Use(AuthMiddleware.MiddlewareFunc())
	Article := AdminGroup.Group(`/article`)
	Article.Use(AuthMiddleware.MiddlewareFunc())
	{
		Article.POST(`/add`, admin.AddArticleHandler)
		Article.GET(`/info`, admin.GetInfoHandler)
		Article.GET(`/list`, admin.GetArticleHandler)
		Article.POST(`/modifiy`, admin.ModifiyArticleHandler)
		Article.DELETE(`/delete/:id`, admin.DeleteArticleHandler)
		Article.POST(`/online`, admin.OnlineAndOfflineHandler)
	}

}
