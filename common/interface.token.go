package common

import (
	"github.com/shopspring/decimal"
)

type Token struct {
	Address  string          `json:"address,omitempty"`
	Name     string          `json:"name,omitempty"`
	Symbol   string          `json:"symbol,omitempty"`
	Decimals uint            `json:"decimals,omitempty"`
	Balance  decimal.Decimal `json:"balance,omitempty"`
	Price    decimal.Decimal `json:"price,omitempty"`
	Icon     string          `json:"icon,omitempty"`
}
