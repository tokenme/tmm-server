package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin"
)

func AdminRouter(r *gin.Engine) {
	AdminGroup := r.Group("/admin")
	r.POST(`/admin/login`, AuthMiddleware.LoginHandler)
	AdminGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		AdminGroup.POST(`/add`, admin.AddArticleHandler)
		AdminGroup.GET(`/info`, admin.GetInfoHandler)
		AdminGroup.GET(`/list/:sortid`, admin.GetArticleHandler)
		AdminGroup.POST(`/modifiy`, admin.ModifiyArticleHandler)
		AdminGroup.DELETE(`/delete/:id`, admin.DeleteArticleHandler)
	}

}
