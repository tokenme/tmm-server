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
	ExchangePointToUcoin     decimal.Decimal `json:"exchange_point_to_ucoin"`
	ExchangeCount            int             `json:"exchange_count"`
	OnlineBFNumber           int             `json:"online_bf_number"`
	OffLineBFNumber          int             `json:"off_line_bf_number"`
	PointByShare             int             `json:"point_by_share"`
	PointByReading           int             `json:"point_by_reading"`
	PointByInvite            int             `json:"point_by_invite"`
	PointByDownLoadApp       int             `json:"point_by_down_load_app"`
	DrawCashByPoint          string          `json:"draw_cash_by_point,omitempty"`
	DrawCashByUc             string          `json:"draw_cash_by_uc,omitempty"`
	DirectFriends            int             `json:"direct_friends"`
	IndirectFriends          int             `json:"indirect_friends"`
	DirectBlockedCount       int             `json:"direct_blocked_count"`
	InDirectBlockedCount     int             `json:"indirect_blocked_count"`
	ActiveFriends            int             `json:"active_friends"`
	TotalMakePoint           int             `json:"total_make_point"`
	DeviceList               []*Device       `json:"device_list"`
	InviteNewUserActiveCount int             `json:"invite_new_user_active_number"`
	InviteNewUserByThreeDays int             `json:"invite_new_user_by_three_days"`
	IsHaveEmulatorDevices    bool            `json:"is_have_emulator_devices"`
	InsertedAt               string          `json:"inserted_at,omitempty"`
	IsActive                 bool            `json:"is_active"`
	FirstDayActive           bool            `json:"first_day_active"`
	SecondDayActive          bool            `json:"second_day_active"`
	ThreeDayActive           bool            `json:"three_day_active"`
	NotActive                string          `json:"not_active"`
	IsHaveAppId              bool            `json:"is_have_app_id"`
	Bonus                    string          `json:"bonus,omitempty"`
	FirendType               string          `json:"firend_type,omitempty"`
	Parent                   User            `json:"pre_user,omitempty"`
	Root                     User            `json:"root,omitempty"`
	OtherAccount             []string        `json:"other_account"`
	BlockedMessage           string          `json:"blocked_message"`
	WxInsertedAt             string          `json:"wx_inserted_at"`
	WxOpenId                 string          `json:"wx_open_id"`
	WxUnionId                string          `json:"wx_union_id"`
	WxExpires                string          `json:"wx_expires"`
	User
}
type User struct {
	common.User
	Blocked            int             `json:"blocked"`
	Point              decimal.Decimal `json:"point"`
	InviteBonus        decimal.Decimal `json:"invite_bonus,omitempty"`
	DrawCash           string          `json:"draw_cash"`
	InviteCount        int             `json:"invite_count,omitempty"`
	Tmm                string          `json:"tmm"`
	CompletedTaskCount int             `json:"completed_task_count,omitempty"`
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
