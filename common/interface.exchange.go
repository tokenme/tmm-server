package common

import (
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	commonutils "github.com/tokenme/tmm/utils"
)

type Transaction struct {
	Receipt    string          `json:"receipt"`
	Status     int             `json:"status"`
	Value      decimal.Decimal `json:"value"`
	InsertedAt string          `json:"inserted_at"`
}

type ExchangeRate struct {
	Rate      decimal.Decimal `json:"rate"`
	MinPoints decimal.Decimal `json:"min_points"`
}

func GetExchangeRate(config Config, service *Service) (ExchangeRate, error) {
	exchRate := ExchangeRate{}
	privateKey, err := commonutils.AddressDecrypt(config.TMMPoolWallet.Data, config.TMMPoolWallet.Salt, config.TMMPoolWallet.Key)
	if err != nil {
		return exchRate, err
	}
	pubKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if err != nil {
		return exchRate, err
	}
	token, err := utils.NewToken(config.TMMTokenAddress, service.Geth)
	if err != nil {
		return exchRate, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return exchRate, err
	}
	balance, err := utils.TokenBalanceOf(token, pubKey)
	if err != nil {
		return exchRate, err
	}
	remainSeconds := commonutils.YearRemainSeconds()
	if err != nil {
		return exchRate, err
	}
	balanceDecimal := decimal.NewFromBigInt(balance, 0)
	remainSecondsDecimal := decimal.NewFromFloat(remainSeconds)
	tmmPerSecond := balanceDecimal.Div(remainSecondsDecimal).Div(decimal.New(1, int32(tokenDecimal)))
	db := service.Db
	rows, _, err := db.Query(`SELECT SUM(d.points) AS points, SUM(d.total_ts) AS ts FROM tmm.devices AS d`)
	if err != nil {
		return exchRate, err
	}
	points, err := decimal.NewFromString(rows[0].Str(0))
	if err != nil {
		return exchRate, err
	}
	ts := decimal.New(rows[0].Int64(1), 0)
	{
		query := `SELECT
    SUM( IFNULL(points, 0) ) AS points,
    SUM( IFNULL(ts, 0) ) AS ts
FROM
    (
    SELECT
        1,
        SUM( dat.points ) AS points,
        COUNT(*) * %d AS ts
    FROM
        tmm.device_app_tasks AS dat
    WHERE
        dat.STATUS = 1 UNION
    SELECT
        2,
        SUM( dst.points ) AS points,
        COUNT(*) * %d AS ts
    FROM
    tmm.device_share_tasks AS dst
    ) AS tmp`
		rows, _, err := db.Query(query, config.DefaultAppTaskTS, config.DefaultShareTaskTS)
		if err != nil {
			return exchRate, err
		}
		taskPoints, err := decimal.NewFromString(rows[0].Str(0))
		if err != nil {
			return exchRate, err
		}
		taskTS := decimal.New(rows[0].Int64(1), 0)
		points = points.Add(taskPoints)
		ts = ts.Add(taskTS)
	}
	pointsPerTs := points.Div(ts)
	exchRate.Rate = tmmPerSecond.Div(pointsPerTs)
	minTMMExchange := decimal.New(int64(config.MinTMMExchange), 0)
	exchRate.MinPoints = pointsPerTs.Div(tmmPerSecond).Mul(minTMMExchange)
	return exchRate, nil
}
