package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/Verify"
	"github.com/tokenme/tmm/handler/admin/article"

)

func AricleRouter(r *gin.Engine) {
	aricle:=r.Group(`admin/article`)
	aricle.Use(AuthMiddleware.MiddlewareFunc(),Verify.VerifyAdminFunc())
	{
		aricle.POST(`/add`,article.AddArticleHandler)
		aricle.GET(`/type`,article.CategoryListHandler)
		aricle.DELETE(`/delete/:id`,article.DeleteArticleHandler)
		aricle.POST(`/edit`,article.EditArticleHandler)
		aricle.GET(`/list`,article.GetArticleHandler)
	}
	}
