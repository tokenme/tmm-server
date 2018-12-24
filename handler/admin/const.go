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
	Point                decimal.Decimal `json:"point"`
	InviteBonus          decimal.Decimal `json:"invite_bonus,omitempty"`
	DrawCash             string          `json:"draw_cash"`
	InviteCount          int             `json:"invite_count,omitempty"`
	Tmm                  decimal.Decimal `json:"tmm"`
	ExchangePointToUcoin decimal.Decimal `json:"exchange_point_to_ucoin"`
	ExchangeCount        int             `json:"exchange_count"`
	CompletedTaskCount   int             `json:"completed_task_count,omitempty"`
	OnlineBFNumber       int             `json:"online_bf_number"`
	OffLineBFNumber      int             `json:"off_line_bf_number"`

	Blocked         int    `json:"blocked,omitempty"`
	PointByShare    string `json:"point_by_share,omitempty"`
	PointByReading  string `json:"point_by_reading,omitempty"`
	PointByInvite   string `json:"point_by_invite,omitempty"`
	DrawCashByPoint string `json:"draw_cash_by_point,omitempty"`
	DrawCashByUc    string `json:"draw_cash_by_uc,omitempty"`
	DirectFriends   int    `json:"direct_friends,omitempty"`
	IndirectFriends int    `json:"indirect_friends,omitempty"`
	ActiveFriends   int    `json:"active_friends,omitempty"`
	common.User
}