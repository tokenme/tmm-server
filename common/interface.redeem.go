package common

import (
	"github.com/shopspring/decimal"
)

type RedeemCdp struct {
	OfferId uint64          `json:"offer_id"`
	Grade   string          `json:"grade"`
	Price   decimal.Decimal `json:"price"`
	Points  decimal.Decimal `json:"points"`
}

type RedeemCdpSlice []*RedeemCdp

func (c RedeemCdpSlice) Len() int {
	return len(c)
}
func (c RedeemCdpSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c RedeemCdpSlice) Less(i, j int) bool {
	return c[i].Price.LessThan(c[j].Price)
}

type RedeemCdpRecord struct {
	DeviceId   string          `json:"device_id"`
	Points     decimal.Decimal `json:"points"`
	Grade      string          `json:"grade"`
	InsertedAt string          `json:"inserted_at"`
}

type TMMWithdrawResponse struct {
	TMM      decimal.Decimal `json:"tmm"`
	Cash     decimal.Decimal `json:"cash"`
	Currency string          `json:"currency"`
}
