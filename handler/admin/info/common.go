package info

import (
	"github.com/shopspring/decimal"
)

const (
	KeyAlive = 60 * 60 * 1
)

//活跃
type Alive struct {
	Day  float64 `json:"day"`
	Week float64 `json:"week"`
}

//页面打开率
type PageInfo struct {
	Wallet      float64 `json:"wallet"`
	PointMake   float64 `json:"point_make"`
	Mall        float64 `json:"mall"`
	Help        float64 `json:"help"`
	Interests   float64 `json:"interests"`
	FeedBack    float64 `json:"feed_back"`
	ChannelOpen map[string]float64
}

//分享拉新
type InviteInfo struct {
	TotalInvite          int             `json:"total_invite"`
	InviteProportionRate float64         `json:"invite_proportion_rate"`
	InviteSucRate        float64         `json:"invite_suc_rate"`
	TotalCost            decimal.Decimal `json:"total_cost"`
	ClickWebRate         float64         `json:"click_web_rate"`
}

//投资
type InvestsInfo struct {
	TotalPoint           decimal.Decimal `json:"total_point"`
	TotalGoods           int             `json:"total_goods"`
	AvgGoodsInvestsPoint decimal.Decimal `json:"avg_goods_invests_point"`
	Goodshare            GoodShare       `json:"goodshare"`
}

//商品分享次数和总流量
type GoodShare struct {
	ShareTimes int `json:"share_times"`
	TotalClick int `json:"total_click"`
}

//提现
type TixianInfo struct {
	TotalTimes          int             `json:"total_times"`
	TotalUser           int             `json:"total_user"`
	TotalMoney          decimal.Decimal `json:"total_money"`
	AvgUserTimes        float64         `json:"avg_user_times"`
	AvgUserMoney        decimal.Decimal `json:"avg_user_money"`
	LessTen             int             `json:"less_ten"`
	LessHundred         int             `json:"less_hundred"`
	LessThousand        int             `json:"less_thousand"`
	LessTenThousand     int             `json:"less_ten_thousand"`
	MoreThanTenThousand int             `json:"more_than_ten_thousand"`
}

type PointInfo struct {
	TotalPoint          decimal.Decimal `json:"total_point"`
	TotalRecovery       decimal.Decimal `json:"total_recovery"`
	CurrentPoint        decimal.Decimal `json:"current_point"`
	AvgDayPoint         decimal.Decimal `json:"avg_day_point"`
	AvgUserPoint        decimal.Decimal `json:"avg_user_point"`
	LessHundred         int             `json:"less_hundred"`
	LessThousand        int             `json:"less_thousand"`
	LessTenThousand     int             `json:"less_ten_thousand"`
	MoreThanTenThousand int             `json:"more_than_ten_thousand"`
}

//任务
type TaskInfo struct {
	TotalTask    int             `json:"total_task"`
	TotalPoint   decimal.Decimal `json:"total_point"`
	ReadTask     TaskDate        `json:"read_task"`
	ShareTask    TaskDate        `json:"share_task"`
	AvgUserRead  float64         `json:"avg_read_task"`
	AvgUserShare float64         `json:"avg_share_task"`
	SignTask     float64         `json:"sign_task"`
}

type TaskDate struct {
	TotalTask int
	Point     decimal.Decimal
}

//Ucoin
type UcInfo struct {
	TotalSupply     decimal.Decimal `json:"total_supply"`
	TotalRecovery   decimal.Decimal `json:"totalrecovery"`
	CurrentSupply   decimal.Decimal `json:"current_supply"`
	DaySupply       decimal.Decimal `json:"day_supply"`
	AvgPersonSupply decimal.Decimal `json:"avg_person_supply"`
}

type ExChangeInfo struct {
	PointToTmm Data `json:"point_to_tmm"`
	TmmToPoint Data `json:"tmm_to_point"`
	Daytimes   int  `json:"daytimes"`
	AvgUser    Data `json:"avg_user"`
}

type data struct {
	Total decimal.Decimal `json:"total"`
	Times float64         `json:"times"`
}

type DevicesInfo struct {
	TotalDevices        int                 `json:"total_devices"`
	TotalIosDevices     int                 `json:"total_ios_devices"`
	TotalAndroidDevices int                 `json:"total_android_devices"`
	UserDownloadChannel UserDownloadChannel `json:"user_download_channel"`
	NewUser             NewUser             `json:"new_user"`
}

//用户下载渠道
type UserDownloadChannel struct {
	InviteDownload int `json:"invite_download"`
	NormalDownload int `json:"normal_download"`
}

//新增用户
type NewUser struct {
	Month int `json:"month"`
	Week  int `json:"week"`
	Day   int `json:"day"`
}

type Data struct {
	Title     string   `json:"title"`
	IndexName []string `json:"index_name"`
	Value     []int    `json:"value"`
}

type User struct {
	Id          int             `json:"id,omitempty"`
	CountryCode int             `json:"country_code,omitempty"`
	Mobile      string          `json:"mobile,omitempty"`
	Nick        string          `json:"nick,omitempty"`
	WxNick      string          `json:"wx_nick"`
	Point       decimal.Decimal `json:"point"`
}

type InfoRequest struct {
	StartTime string `json:"start_time",form:"start_time"`
	EndTime   string `json:"end_time",form:"end_time" `
}
