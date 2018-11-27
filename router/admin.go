package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"github.com/tokenme/tmm/handler/qiniu"
)

func AdminRouter(r *gin.Engine) {
	AdminGroup := r.Group("/admin")
	AdminAuth := AuthMiddleware
	AdminAuth.Unauthorized = func(c *gin.Context, code int, message string) {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": message,})
	}
	r.POST(`/admin/login`, AdminAuth.LoginHandler)
	AdminGroup.Use(AuthMiddleware.MiddlewareFunc())

	/* @Router: admin/article 						*/
	Article := AdminGroup.Group(`/article`)
	Article.Use(AuthMiddleware.MiddlewareFunc(), admin.VerfiyAdminFunc())
	{
		//Article.POST(`/add`, admin.AddArticleHandler)
		//Article.POST(`/modifiy`, admin.ModifiyArticleHandler)
		//Article.DELETE(`/delete/:id`, admin.DeleteArticleHandler)
		//Article.GET(`/list`, admin.GetArticleHandler)
		Article.GET(`/info`, admin.GetInfoHandler)
		Article.POST(`/online`, admin.OnlineAndOfflineHandler)
		Article.GET(`/type`, admin.GetTypeHander)
	}

	/* @Router: admin/share 						*/
	share := AdminGroup.Group(`/share`)
	share.Use(AuthMiddleware.MiddlewareFunc(), admin.VerfiyAdminFunc())
	{
		share.GET(`/getToken`, qiniu.UpTokenHandler)
		share.POST(`/add`, admin.AddShareHandler)
		share.GET(`/list`, admin.GetShareListHandler)
	}

}
