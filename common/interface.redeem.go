package common

import (
	"github.com/shopspring/decimal"
	"time"
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

func ExceededDailyWithdraw(cash decimal.Decimal, service *Service, config Config) (exceeded bool, dailyTotal decimal.Decimal, chunkBudget decimal.Decimal, nextHour time.Time, err error) {
	db := service.Db
	{
		rows, _, err := db.Query(`SELECT SUM(cny) FROM
(
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.withdraw_txs AS wt WHERE tx_status!=0 AND inserted_at>=DATE(NOW())
UNION ALL
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.point_withdraws AS pw WHERE inserted_at>=DATE(NOW())
) AS t`)
		if err != nil {
			return exceeded, dailyTotal, chunkBudget, nextHour, err
		}
		if len(rows) == 0 {
			return exceeded, dailyTotal, chunkBudget, nextHour, nil
		}
		dailyTotal, _ = decimal.NewFromString(rows[0].Str(0))
	}
	var (
		now                   = time.Now()
		tomorrow              = now.Add(24 * time.Hour)
		chunkSize         int = 3
		maxChunks         int = 24 / chunkSize
		currentChunk          = now.Hour() / chunkSize
		startHour             = currentChunk * chunkSize
		chunckTime            = time.Date(now.Year(), now.Month(), now.Day(), startHour, 0, 0, 0, now.Location())
		currentChunkTotal decimal.Decimal
	)
	maxDailyWithdraw := decimal.New(config.MaxDailWithdraw, 0)
	exceeded = dailyTotal.GreaterThan(maxDailyWithdraw)
	if exceeded {
		nextHour = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
		return exceeded, dailyTotal, chunkBudget, nextHour, nil
	}
	nextHour = chunckTime.Add(time.Duration(chunkSize) * time.Hour)
	{

		rows, _, err := db.Query(`SELECT SUM(cny) FROM
(
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.withdraw_txs AS wt WHERE tx_status!=0 AND inserted_at>='%s'
UNION ALL
SELECT IFNULL(SUM(cny), 0) AS cny FROM tmm.point_withdraws AS pw WHERE inserted_at>='%s'
) AS t`, chunckTime.Format("2006-01-02 15:04:05"), chunckTime.Format("2006-01-02 15:04:05"))
		if err != nil {
			return exceeded, dailyTotal, chunkBudget, nextHour, err
		}
		if len(rows) == 0 {
			return exceeded, dailyTotal, chunkBudget, nextHour, nil
		}
		currentChunkTotal, _ = decimal.NewFromString(rows[0].Str(0))
	}

	todayLeft := maxDailyWithdraw.Sub(dailyTotal.Sub(currentChunkTotal))
	chunksLeft := decimal.New(int64(maxChunks-currentChunk), 0)
	chunkBudget = todayLeft.Div(chunksLeft)
	exceeded = chunkBudget.LessThan(cash)
	return exceeded, dailyTotal, chunkBudget, nextHour, nil
}
