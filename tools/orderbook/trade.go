package orderbook

import (
	"github.com/shopspring/decimal"
)

type Trade struct {
	Timestamp          int
	Id                 uint64
	Wallet             string
	CounterParty       uint64
	CounterPartyWallet string
	Quantity           decimal.Decimal
	Price              decimal.Decimal
	Side               Side
}
