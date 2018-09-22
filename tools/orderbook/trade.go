package orderbook

import (
	"github.com/shopspring/decimal"
)

type Trade struct {
	Timestamp    int
	Id           uint64
	CounterParty uint64
	Quantity     decimal.Decimal
	Price        decimal.Decimal
}
