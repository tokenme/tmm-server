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

type TMMWithdrawRecord struct {
	TMM            decimal.Decimal `json:"tmm"`
	Cash           decimal.Decimal `json:"cash"`
	Tx             string          `json:"tx"`
	TxStatus       uint            `json:"tx_status"`
	WithdrawStatus uint            `json:"withdraw_status"`
	InsertedAt     string          `json:"inserted_at"`
}

func ExceededDailyWithdraw(service *Service, config Config) (exceeded bool, total decimal.Decimal, err error) {
	db := service.Db
	rows, _, err := db.Query(`SELECT SUM(cny) FROM
(
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.withdraw_txs AS wt WHERE tx_status!=0 AND inserted_at>=DATE(NOW())
UNION ALL
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.point_withdraws AS pw WHERE inserted_at>=DATE(NOW())
) AS t`)
	if err != nil {
		return exceeded, total, err
	}
	if len(rows) == 0 {
		return exceeded, total, nil
	}
	maxDailyWithdraw := decimal.New(config.MaxDailWithdraw, 0)
	total, _ = decimal.NewFromString(rows[0].Str(0))
	return total.GreaterThan(maxDailyWithdraw), total, nil
}
