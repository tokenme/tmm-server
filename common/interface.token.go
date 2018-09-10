package common

import (
	"github.com/shopspring/decimal"
)

type Token struct {
	Name     string          `json:"name,omitempty"`
	Symbol   string          `json:"symbol,omitempty"`
	Decimals uint            `json:"decimals,omitempty"`
	Balance  decimal.Decimal `json:"balance,omitempty"`
}
