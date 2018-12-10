package info

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
)

const (
	KeyAlive = 60 * 60
)

type Stats struct {
	Top10       []*Users `json:"top_10,omitempty"`
	Numbers     int      `json:"numbers"`
	CurrentTime string   `json:"current_time"`
	Title       string   `json:"title"`
}

//提现
type DrawCashStats struct {
	Money decimal.Decimal `json:"money"`
	Stats
}

type TotalDrawCash struct {
	TotalCount int             `json:"total_count"`
	TotalUser  int             `json:"total_user"`
	TotalMoney decimal.Decimal `json:"total_money"`
}

//投资
type InvestsStats struct {
	InvestsPoint decimal.Decimal `json:"invests_point"`
	Stats
	Top10 []*Good `json:"top_10,omitempty"`
}

type TotalInvests struct {
	TotalPoint      decimal.Decimal `json:"total_point"`
	TotalGoodsCount int             `json:"total_goods_count"`
}

//积分
type PointStats struct {
	Point decimal.Decimal `json:"point"`
	Stats
}

//邀请
type InviteStats struct {
	InviteCount int `json:"invite_count"`
	Stats
}

type TotalInvite struct {
	TotalInviteCount int             `json:"total_invite_count"`
	TotalCost        decimal.Decimal `json:"total_cost"`
}

//交换
type ExchangeStats struct {
	ExchangeCount int `json:"exchange_count"`
	Stats
}

//任务
type TaskStats struct {
	TaskCount  int             `json:"task_count"`
	TotalPoint decimal.Decimal `json:"total_point,omitempty"`
	Stats
}

type TotalTask struct {
	TotaltaskCount int             `json:"totaltask_count"`
	TotalCost      decimal.Decimal `json:"total_cost"`
}

//其他类
type Users struct {
	Point              decimal.Decimal `json:"point,omitempty"`
	DrawCash           decimal.Decimal `json:"draw_cash,omitempty"`
	InviteCount        int             `json:"invite_count,omitempty"`
	Tmm                decimal.Decimal `json:"tmm,omitempty"`
	ExchangeCount      int             `json:"exchange_count,omitempty"`
	CompletedTaskCount int             `json:"completed_task_count,omitempty"`
	common.User
}

type StatsRequest struct {
	StartTime string `form:"start_time",json:"start_time"`
	EndTime   string `form:"end_time",json:"end_time" `
	Top10     bool   `form:"top_10",json:"top_10"`
}

type Good struct {
	Id    int             `json:"id"`
	Title string          `json:"title"`
	Point decimal.Decimal `json:"point"`
}
type Data struct {
	Title     string   `json:"title"`
	IndexName []string `json:"index_name"`
	Value     []int    `json:"value"`
}