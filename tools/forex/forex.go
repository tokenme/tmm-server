package forex

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	ex "github.com/me-io/go-swap/pkg/exchanger"
	"github.com/me-io/go-swap/pkg/swap"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
)

const (
	FOREX_CACHE_KEY = "forex_cache_%s_%s"
)

func Rate(service *common.Service, from string, to string) decimal.Decimal {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	cacheKey := fmt.Sprintf(FOREX_CACHE_KEY, from, to)
	rateStr, err := redis.String(redisConn.Do("GET", cacheKey))
	var rateValue decimal.Decimal
	if err == nil {
		rateValue, err = decimal.NewFromString(rateStr)
		if err == nil {
			return rateValue
		}
	}
	client := swap.NewSwap()
	client.AddExchanger(ex.NewGoogleApi(nil)).AddExchanger(ex.NewYahooApi(nil)).AddExchanger(ex.NewTheMoneyConverterApi(nil)).Build()
	exchanger := client.Latest(fmt.Sprintf("%s/%s", from, to))
	rateValue = decimal.NewFromFloat(exchanger.GetRateValue())
	redisConn.Do("SETEX", cacheKey, 12*60*60, rateValue.StringFixed(9))
	return rateValue
}
