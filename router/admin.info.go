package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/admin/info"
)

func InfoRouter(r *gin.Engine) {
	InfoGroup := r.Group(`/admin/info`)
	InfoGroup.Use(AuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		InfoGroup.GET(`/drawcash/data`, info.DrawCashDataHandler)
		InfoGroup.GET(`/drawcash/info`, info.DrawCashStatsHandler)
		InfoGroup.GET(`/drawcash/total`, info.TotalDrawCashHandler)
	}
	{
		InfoGroup.GET(`/exchange/data`, info.ExchangeDataHandler)
		InfoGroup.GET(`/exchange/info`, info.ExchangeStatsHandler)
	}
	{
		InfoGroup.GET(`/invest/data`, info.InvestsDataHandler)
		InfoGroup.GET(`/invest/info`, info.InvestsStatsHandler)
		InfoGroup.GET(`/invest/total`, info.TotalInvestHandler)
	}
	{
		InfoGroup.GET(`/invite/data`, info.InviteDataHandler)
		InfoGroup.GET(`/invite/info`, info.InviteStatsHandler)
		InfoGroup.GET(`/invite/total`, info.TotalInviteHandler)

	}
	{
		InfoGroup.GET(`/point/data`, info.PointDataHandler)
		InfoGroup.GET(`/point/info`, info.PointStatsHandler)
	}
	{
		InfoGroup.GET(`/task/data`, info.TaskDataHandler)
		InfoGroup.POST(`/task/info`, info.TaskStatsHandler)
		InfoGroup.GET(`/task/total`, info.TotalTaskHandler)
	}
}
