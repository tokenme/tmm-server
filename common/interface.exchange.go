package common

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	commonutils "github.com/tokenme/tmm/utils"
)

type TMMExchangeDirection = int8

const (
	TMMExchangeIn  TMMExchangeDirection = 1
	TMMExchangeOut TMMExchangeDirection = -1
)

type Transaction struct {
	Receipt           string          `json:"receipt"`
	Status            int             `json:"status"`
	Value             decimal.Decimal `json:"value"`
	From              string          `json:"from,omitempty"`
	To                string          `json:"to,omitempty"`
	Gas               decimal.Decimal `json:"gas,omitempty"`
	GasPrice          decimal.Decimal `json:"gas_price,omitempty"`
	GasUsed           decimal.Decimal `json:"gas_used,omitempty"`
	CumulativeGasUsed decimal.Decimal `json:"cumulative_gas_used,omitempty"`
	Confirmations     int             `json:"confirmations,omitempty"`
	InsertedAt        string          `json:"inserted_at"`
}

type ExchangeRate struct {
	Rate      decimal.Decimal `json:"rate"`
	MinPoints decimal.Decimal `json:"min_points"`
}

func GetExchangeRate(config Config, service *Service) (ExchangeRate, decimal.Decimal, error) {
	exchRate := ExchangeRate{}
	pointsPerTs := decimal.New(0, 0)
	tmmPerTs, err := GetTMMPerTs(config, service)
	if err != nil {
		return exchRate, pointsPerTs, err
	}
	pointsPerTs, err = GetPointsPerTs(service)
	if err != nil {
		return exchRate, pointsPerTs, err
	}
	exchRate.Rate = tmmPerTs.Div(pointsPerTs)
	minTMMExchange := decimal.New(int64(config.MinTMMExchange), 0)
	exchRate.MinPoints = pointsPerTs.Div(tmmPerTs).Mul(minTMMExchange)
	return exchRate, pointsPerTs, nil
}

func GetPointsPerTs(service *Service) (decimal.Decimal, error) {
	pointsPerTs := decimal.New(0, 0)
	db := service.Db
	rows, _, err := db.Query(`SELECT SUM(d.points) AS points, SUM(d.total_ts) - SUM(d.consumed_ts) AS ts FROM tmm.devices AS d`)
	if err != nil {
		return pointsPerTs, err
	}
	points, err := decimal.NewFromString(rows[0].Str(0))
	if err != nil {
		return pointsPerTs, err
	}
	ts := decimal.New(rows[0].Int64(1), 0)
	pointsPerTs = points.Div(ts)
	return pointsPerTs, nil
}

func GetTMMPerTs(config Config, service *Service) (decimal.Decimal, error) {
	tmmPerTs := decimal.New(0, 0)
	privateKey, err := commonutils.AddressDecrypt(config.TMMPoolWallet.Data, config.TMMPoolWallet.Salt, config.TMMPoolWallet.Key)
	if err != nil {
		return tmmPerTs, err
	}
	pubKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return tmmPerTs, err
	}
	token, err := utils.NewToken(config.TMMTokenAddress, service.Geth)
	if err != nil {
		return tmmPerTs, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return tmmPerTs, err
	}
	balance, err := utils.TokenBalanceOf(token, pubKey)
	if err != nil {
		return tmmPerTs, err
	}
	remainSeconds := commonutils.YearRemainSeconds()
	if err != nil {
		return tmmPerTs, err
	}
	balanceDecimal := decimal.NewFromBigInt(balance, 0)
	remainSecondsDecimal := decimal.NewFromFloat(remainSeconds)
	tmmPerTs = balanceDecimal.Div(remainSecondsDecimal).Div(decimal.New(1, int32(tokenDecimal)))
	return tmmPerTs, nil
}
