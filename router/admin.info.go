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
		InfoGroup.POST(`/drawcash/info`, info.DrawCashInfoHandler)
		InfoGroup.GET(`/drawcash/total`, info.TotalDrawCashHandler)
	}
	{
		InfoGroup.GET(`/exchange/data`, info.ExchangeDataHandler)
		InfoGroup.POST(`/exchange/info`, info.ExchangeInfoHandler)
	}
	{
		InfoGroup.GET(`/invest/data`, info.InvestsDataHandler)
		InfoGroup.POST(`/invest/info`, info.InvestsInfoHandler)
		InfoGroup.GET(`/invest/total`, info.TotalInvestHandler)
	}
	{
		InfoGroup.GET(`/invite/data`, info.InviteDataHandler)
		InfoGroup.POST(`/invite/info`, info.InviteInfoHandler)
		InfoGroup.GET(`/invite/total`, info.TotalInviteHandler)

	}
	{
		InfoGroup.GET(`/point/data`, info.PointDataHandler)
		InfoGroup.POST(`/point/info`, info.PointInfoHandler)
	}
	{
		InfoGroup.GET(`/task/data`, info.TaskDataHandler)
		InfoGroup.POST(`/task/info`, info.TaskInfoHandler)
		InfoGroup.GET(`/task/total`, info.TotalTaskHandler)
	}
}
