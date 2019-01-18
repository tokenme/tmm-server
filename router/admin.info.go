package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/info"
	"github.com/tokenme/tmm/handler/admin/verify"
)

func InfoRouter(r *gin.Engine) {
	InfoGroup := r.Group(`/admin/info`)
	InfoGroup.Use(AdminAuthMiddleware.MiddlewareFunc(), verify.VerifyAdminFunc())
	{
		InfoGroup.GET(`/drawcash/data`, info.DrawCashDataHandler)
		InfoGroup.GET(`/drawcash/info`, info.DrawCashStatsHandler)
		InfoGroup.GET(`/drawcash/total`, info.TotalDrawCashHandler)
	}
	{
		InfoGroup.GET(`/exchange/data`, info.ExchangeDataHandler)
		InfoGroup.GET(`/exchange/info`, info.ExchangeStatsHandler)
		InfoGroup.GET(`/exchange/rate`, info.ExchangeRateHandler)
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
		InfoGroup.GET(`/task/info`, info.TaskStatsHandler)
		InfoGroup.GET(`/task/total`, info.TotalTaskHandler)
	}
	{
		InfoGroup.GET(`/user/total`, info.UserStatsHandler)
	}
	{
		InfoGroup.GET(`/stats/stats`, info.StatsHandler)
	}
	{
		InfoGroup.GET(`/trend/drawcash`, info.DrawCashTrendHandler)
		InfoGroup.GET(`/trend/exchange`, info.ExchangeTrendHandler)
		InfoGroup.GET(`/trend/point`, info.PointTrendHandler)
		InfoGroup.GET(`/trend/invite`, info.InviteTrendHandler)
		InfoGroup.GET(`/trend/user`, info.UserTrendHandler)
		InfoGroup.GET(`/trend/task`, info.TaskTrendHandler)
		InfoGroup.GET(`/trend/app`, info.AppTrendHandler)
		InfoGroup.GET(`/trend/share`, info.ShareTrendHandler)
		InfoGroup.GET(`/trend/uc`, info.UcTrendHandler)
		InfoGroup.GET(`/trend/stats`, info.UserFunnelStatsHandler)
	}
	{
		InfoGroup.GET(`/funnel/stats`, info.UserFunnelStatsHandler)
		InfoGroup.GET(`/funnel/data`, info.GetFunnelDataHandler)
	}
	{
		InfoGroup.GET(`/current/drawcash/data`, info.GetWithDrawDataHandler)
		InfoGroup.GET(`/current/drawcash/stats`, info.GetTodayWithDrawStatsHandler)
	}
}
