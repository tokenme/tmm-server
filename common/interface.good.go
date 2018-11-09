package common

import (
	"github.com/shopspring/decimal"
)

type GoodInvest struct {
	Nick   string          `json:"user_name"`
	Points decimal.Decimal `json:"points"`
	Avatar string          `json:"avatar"`
}
