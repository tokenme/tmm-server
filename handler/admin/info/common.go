package info

import (
	"github.com/shopspring/decimal"
)

const (
	KeyAlive = 60 * 60
)

//提现
type DrawCashInfo struct {
	Top10       []*User         `json:"top_10,omitempty"`
	Money       decimal.Decimal `json:"money"`
	CurrentTime string          `json:"current_time"`
}

type TotalDrawCash struct {
	TotalCount int             `json:"total_count"`
	TotalUser  int             `json:"total_user"`
	TotalMoney decimal.Decimal `json:"total_money"`
}

//投资
type InvestsInfo struct {
	Top10        []*Good         `json:"top_10,omitempty"`
	InvestsPoint decimal.Decimal `json:"invests_point"`
	CurrentTime  string          `json:"current_time"`
}

type TotalInvests struct {
	TotalPoint      decimal.Decimal `json:"total_point"`
	TotalGoodsCount int             `json:"total_goods_count"`
}

//积分
type PointInfo struct {
	CurrentTime string          `json:"current_time"`
	Point       decimal.Decimal `json:"point"`
	Top10       []*User         `json:"top_10,omitempty"`
}

//邀请
type InviteInfo struct {
	InviteCount int     `json:"invite_count"`
	Top10       []*User `json:"top_10,omitempty"`
	CurrentTime string  `json:"current_time"`
}

type TotalInvite struct {
	TotalInviteCount int             `json:"total_invite_count"`
	TotalCost        decimal.Decimal `json:"total_cost"`
}

//交换
type ExchangeInfo struct {
	ExchangeCount int     `json:"exchange_count"`
	Top10         []*User `json:"top_10,omitempty"`
	CurrentTime   string  `json:"current_time,omitempty"`
}

//任务
type TaskInfo struct {
	TaskCount   int             `json:"task_count"`
	Top10       []*User         `json:"top_10,omitempty"`
	CurrentTime string          `json:"current_time,omitempty"`
	TotalPoint  decimal.Decimal `json:"total_point,omitempty"`
}

type TotalTask struct {
	TotaltaskCount int             `json:"totaltask_count"`
	TotalCost      decimal.Decimal `json:"total_cost"`
}

//其他类
type User struct {
	Id                 int             `json:"id,omitempty"`
	CountryCode        int             `json:"country_code,omitempty"`
	Mobile             string          `json:"mobile,omitempty"`
	Nick               string          `json:"nick,omitempty"`
	WxNick             string          `json:"wx_nick,omitempty"`
	Point              decimal.Decimal `json:"point,omitempty"`
	DrawCash           decimal.Decimal `json:"draw_cash,omitempty"`
	InvestsCount       int             `json:"invests_count,omitempty"`
	Tmm                decimal.Decimal `json:"tmm,omitempty"`
	ExchangeCount      int             `json:"exchange_count,omitempty"`
	CompletedTaskCount int             `json:"completed_task_count,omitempty"`
}

type InfoRequest struct {
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
