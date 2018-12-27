package info

import (
	"github.com/shopspring/decimal"
	"fmt"
	"strconv"
	"github.com/tokenme/tmm/handler/admin"
)

const (
	KeyAlive = 60 * 60
)

type Stats struct {
	Top10       []*admin.Users `json:"top_10,omitempty"`
	Numbers     int            `json:"numbers"`
	CurrentTime string         `json:"current_time"`
	Title       string         `json:"title"`
}

//提现
type DrawCashStats struct {
	Money decimal.Decimal `json:"money"`
	Stats
}

type TotalDrawCash struct {
	TotalCount int             `json:"total_count"`
	TotalUser  int             `json:"total_user"`
	TotalMoney string          `json:"total_money"`
	Uc         decimal.Decimal `json:"uc"`
	Point      decimal.Decimal `json:"point"`
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
type StatsRequest struct {
	StartTime string `form:"start_date",json:"start_date"`
	EndTime   string `form:"end_date",json:"end_date" `
	Top10     bool   `form:"top_10",json:"top_10"`
	Hours     int    `form:"hours" ,json:"hours"`
}

type Good struct {
	Id    int             `json:"id"`
	Title string          `json:"title"`
	Point decimal.Decimal `json:"point"`
}

type Data struct {
	Title  Title    `json:"title"`
	Yaxis  Axis     `json:"yAxis"`
	Xaxis  Axis     `json:"xAxis"`
	Series []Series `json:"series"`
}

type Axis struct {
	Name string   `json:"name"`
	Data []string `json:"data,omitempty"`
}

type Title struct {
	Text string `json:"text"`
}
type Series struct {
	Data []string `json:"data"`
	Name string   `json:"name"`
}

type StatsData struct {
	PointExchangeNumber int    `json:"point_exchange_number"`
	UcoinExchangeNumber int    `json:"ucoin_exchange_number"`
	Cash                string `json:"cash"`
	PointSupply         string `json:"point_supply"`
	UcSupply            string `json:"uc_supply"`
	TotalTaskUser       int    `json:"total_task_user"`
	TotalFinishTask     int    `json:"total_finish_task"`
	InviteNumber        int    `json:"invite_number"`
	Active              int    `json:"active"`
}

type StatsList struct {
	Yesterday StatsData `json:"yesterday"`
	Today     StatsData `json:"today"`
}

func GetPercentList(valueList []string) (PercentList []string) {
	var total float64
	for _, value := range valueList {
		v, _ := strconv.Atoi(value)

		total += float64(v)
	}

	for _, value := range valueList {
		v, _ := strconv.Atoi(value)
		percent := float64(v) / total
		PercentList = append(PercentList, fmt.Sprintf("%.2f", percent*100))
	}

	return
}
