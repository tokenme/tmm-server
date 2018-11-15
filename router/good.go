package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/good"
)

func goodRouter(r *gin.Engine) {
	goodGroup := r.Group("/good")
	goodGroup.Use(AuthMiddleware.MiddlewareFunc())
	{
		goodGroup.GET("/item/:id", good.ItemHandler)
		goodGroup.POST("/invest", good.InvestHandler)
		goodGroup.GET("/invests/item/:id/:page/:pageSize", good.ItemInvestsHandler)
		goodGroup.GET("/invests/my/:page/:pageSize", good.MyInvestsHandler)
		goodGroup.GET("/invest/withdraw/:id", good.InvestWithdrawHandler)
		goodGroup.GET("/invest/summary", good.InvestSummaryHandler)
		goodGroup.POST("/invest/redeem", good.InvestRedeemHandler)
	}
	r.GET("/good/list", good.ListHandler)
	r.POST("/good/txs", good.TxsHandler)
}
