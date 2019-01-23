package admin

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
)

const (
	API_OK    = "OK"
	Not_Found = "没有查找到数据"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

type AddRequest struct {
	Title         string          `json:"title" form:"title" `
	Summary       string          `json:"summary" form:"summary" `
	Link          string          `json:"link" form:"link" `
	Image         string          `json:"image" form:"image"`
	FileExtension string          `json:"image_extension" from:"image_extension"`
	Points        decimal.Decimal `json:"points" form:"points" `
	Bonus         decimal.Decimal `json:"bonus" form:"bonus" `
	MaxViewers    uint            `json:"max_viewers" form:"max_viewers" `
	Cid           []int           `json:"cid" form:"cid"`
}

type UserStats struct {
	TotalMakePoint        int       `json:"total_make_point"`
	DeviceList            []*Device `json:"device_list"`
	IsHaveEmulatorDevices bool      `json:"is_have_emulator_devices"`
	IsActive              bool      `json:"is_active"`
	FirstDayActive        bool      `json:"first_day_active"`
	SecondDayActive       bool      `json:"second_day_active"`
	ThreeDayActive        bool      `json:"three_day_active"`
	NotActive             string    `json:"not_active"`
	IsHaveAppId           bool      `json:"is_have_app_id"`
	FriendBonus           string    `json:"bonus,omitempty"`
	FriendType            string    `json:"firend_type,omitempty"`
	Root                  *User     `json:"root,omitempty"`
	OtherAccount          []string  `json:"other_account"`
	WxInsertedAt          string    `json:"wx_inserted_at"`
	WxOpenId              string    `json:"wx_open_id"`
	WxUnionId             string    `json:"wx_union_id"`
	BlockedMessage        string    `json:"blocked_message"`
	MakePoint
	WithDrawCash
	InviteStatus
	User
	ShareBlocked
}

type Pages struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type MakePoint struct {
	PointByShare             int `json:"point_by_share"`
	PointByReading           int `json:"point_by_reading"`
	PointByInvite            int `json:"point_by_invite"`
	PointByDownLoadApp       int `json:"point_by_down_load_app"`
	DaySign                  int `json:"day_sign"`
	InviteNewUserActiveCount int `json:"invite_new_user_active_number"`
	InviteNewUserByThreeDays int `json:"invite_new_user_by_three_days"`
}

type InviteStatus struct {
	DirectFriends        int `json:"direct_friends"`
	IndirectFriends      int `json:"indirect_friends"`
	DirectBlockedCount   int `json:"direct_blocked_count"`
	InDirectBlockedCount int `json:"indirect_blocked_count"`
	ActiveFriends        int `json:"active_friends"`
}
type WithDrawCash struct {
	DrawCashByPoint string `json:"draw_cash_by_point,omitempty"`
	DrawCashByUc    string `json:"draw_cash_by_uc,omitempty"`
}

type ShareBlocked struct {
	TenMinuteCount int `json:"ten_minute_count,omitempty"`
	OneHourCount   int `json:"one_hour_count,omitempty"`
}

type User struct {
	Blocked              int             `json:"blocked"`
	Point                string          `json:"point"`
	DrawCash             string          `json:"draw_cash"`
	Tmm                  string          `json:"tmm"`
	InviteBonus          decimal.Decimal `json:"invite_bonus,omitempty"`
	InviteCount          int             `json:"invite_count,omitempty"`
	CompletedTaskCount   int             `json:"completed_task_count,omitempty"`
	ExchangePointToUcoin decimal.Decimal `json:"exchange_point_to_ucoin"`
	ExchangeCount        int             `json:"exchange_count"`
	OnlineBFNumber       int             `json:"online_bf_number"`
	OffLineBFNumber      int             `json:"off_line_bf_number"`
	Created              string          `json:"inserted_at,omitempty"`
	Parent               *User           `json:"pre_user,omitempty"`
	CurrentPoint         string          `json:"current_point,omitempty"`
	common.User
}

type Device struct {
	common.Device
	SystemVersion string  `json:"system_version"`
	Language      string  `json:"language"`
	OsVersion     string  `json:"os_version"`
	Timezone      string  `json:"timezone"`
	Country       string  `json:"country"`
	IsJailbrojen  bool    `json:"is_jailbrojen"`
	ConsumedTs    float64 `json:"consumed_ts"`
	TmpTs         uint64  `json:"tmp_ts"`
}
