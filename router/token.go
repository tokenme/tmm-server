package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/token"
)

func tokenRouter(r *gin.Engine) {
	tokenGroup := r.Group("/token")
	tokenGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		tokenGroup.GET("/tmm/balance", token.TMMBalanceHandler)
		tokenGroup.GET("/assets", token.AssetsHandler)
		tokenGroup.GET("/transactions/:address/:page/:pageSize", token.TransactionsHandler)
		tokenGroup.POST("/transfer", handler.ApiSignFunc(), token.TransferHandler)
		tokenGroup.GET("/info/:address", token.InfoHandler)
	}
}
