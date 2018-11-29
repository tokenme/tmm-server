package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/article"

)

func AricleRouter(r *gin.Engine) {
	aricle:=r.Group(`admin/article`)
	aricle.Use(AuthMiddleware.MiddlewareFunc(),verify.VerifyAdminFunc())
	{
		//aricle.POST(`/add`,article.AddArticleHandler)
		aricle.GET(`/type`,article.CategoryListHandler)
		//aricle.DELETE(`/delete/:id`,article.DeleteArticleHandler)
		//aricle.GET(`/edit`,article.GetArticleHandler)
		//aricle.POST(`/modify`,article.ModifyArticleHandler)
		//aricle.GET(`/list`,article.GetArticleListHandler)
	}
	}
