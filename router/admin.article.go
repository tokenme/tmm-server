package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/article"
	"github.com/tokenme/tmm/handler/admin/verify"
)

func AricleRouter(r *gin.Engine) {
	aricle := r.Group(`admin/article`)
	aricle.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		//aricle.POST(`/add`,article.AddArticleHandler)
		aricle.GET(`/type`, article.CategoryListHandler)
		//aricle.DELETE(`/delete/:id`,article.DeleteArticleHandler)
		//aricle.GET(`/edit`,article.GetArticleHandler)
		//aricle.POST(`/modify`,article.ModifyArticleHandler)
		//aricle.GET(`/list`,article.GetArticleListHandler)
	}
}
