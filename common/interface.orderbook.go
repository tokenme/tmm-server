package common

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/orderbook"
)

type Order struct {
	TradeId      uint64                `json:"trand_id,omitempty"`
	Quantity     decimal.Decimal       `json:"quantity,omitempty"`
	Price        decimal.Decimal       `json:"price,omitempty"`
	Side         orderbook.Side        `json:"side,omitempty"`
	ProcessType  orderbook.ProcessType `json:"process_type,omitempty"`
	DealQuantity decimal.Decimal       `json:"deal_quantity,omitempty"`
	DealEth      decimal.Decimal       `json:"deal_eth,omitempty"`
}
