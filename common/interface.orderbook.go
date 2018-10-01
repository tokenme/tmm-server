package common

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/orderbook"
)

type Order struct {
	TradeId      uint64                `json:"trade_id,omitempty"`
	Quantity     decimal.Decimal       `json:"quantity,omitempty"`
	Price        decimal.Decimal       `json:"price,omitempty"`
	Side         orderbook.Side        `json:"side,omitempty"`
	ProcessType  orderbook.ProcessType `json:"process_type,omitempty"`
	DealQuantity decimal.Decimal       `json:"deal_quantity,omitempty"`
	DealEth      decimal.Decimal       `json:"deal_eth,omitempty"`
	OnlineStatus int8                  `json:"online_status,omitempty"`
	InsertedAt   string                `json:"inserted_at,omitempty"`
	UpdatedAt    string                `json:"updated_at,omitempty"`
}

type MarketGraph struct {
	Trades   uint64          `json:"trades,omitempty"`
	Quantity decimal.Decimal `json:"quantity,omitempty"`
	Price    decimal.Decimal `json:"price,omitempty"`
	Low      decimal.Decimal `json:"low,omitempty"`
	High     decimal.Decimal `json:"high,omitempty"`
	At       string          `json:"at,omitempty"`
}
