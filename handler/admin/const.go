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

type Users struct {
	Point                    decimal.Decimal  `json:"point"`
	InviteBonus              decimal.Decimal  `json:"invite_bonus,omitempty"`
	DrawCash                 string           `json:"draw_cash"`
	InviteCount              int              `json:"invite_count,omitempty"`
	Tmm                      decimal.Decimal  `json:"tmm"`
	ExchangePointToUcoin     decimal.Decimal  `json:"exchange_point_to_ucoin"`
	ExchangeCount            int              `json:"exchange_count"`
	CompletedTaskCount       int              `json:"completed_task_count,omitempty"`
	OnlineBFNumber           int              `json:"online_bf_number"`
	OffLineBFNumber          int              `json:"off_line_bf_number"`
	ChildrenNumber           int              `json:"children_number"`
	Blocked                  int              `json:"blocked"`
	PointByShare             int              `json:"point_by_share"`
	PointByReading           int              `json:"point_by_reading"`
	PointByInvite            int              `json:"point_by_invite"`
	PointByDownLoadApp       int              `json:"point_by_down_load_app"`
	DrawCashByPoint          string           `json:"draw_cash_by_point,omitempty"`
	DrawCashByUc             string           `json:"draw_cash_by_uc,omitempty"`
	DirectFriends            int              `json:"direct_friends"`
	IndirectFriends          int              `json:"indirect_friends"`
	ActiveFriends            int              `json:"active_friends"`
	TotalMakePoint           int              `json:"total_make_point"`
	DeviceList               []*common.Device `json:"device_list"`
	InviteNewUserActiveCount int              `json:"invite_new_user_active_number"`
	InviteNewUserByThreeDays int              `json:"invite_new_user_by_three_days"`
	IsHaveEmulatorDevices    bool             `json:"is_have_emulator_devices"`
	InsertedAt               string           `json:"inserted_at"`
	IsActive                 bool             `json:"is_active"`
	FirstDayActive           bool             `json:"first_day_active"`
	SecondDayActive          bool             `json:"second_day_active"`
	ThreeDayActive           bool             `json:"three_day_active"`
	Bonus                    string           `json:"bonus"`
	FirendType               string           `json:"firend_type"`
	common.User
}
