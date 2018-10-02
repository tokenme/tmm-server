package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fabioberger/coinbase-go"
	"github.com/garyburd/redigo/redis"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/ethplorer-api"
)

const (
	TOKEN_INFO_KEY  = "etherplorer-token-%s"
	ETH_PRICE_CACHE = "eth-usd"
)

type Token struct {
	Address      string            `json:"address,omitempty"`
	Name         string            `json:"name,omitempty"`
	Symbol       string            `json:"symbol,omitempty"`
	Decimals     uint              `json:"decimals,omitempty"`
	Balance      decimal.Decimal   `json:"balance,omitempty"`
	Price        decimal.Decimal   `json:"price,omitempty"`
	Icon         string            `json:"icon,omitempty"`
	InitialPrice map[string]string `json:"initial_price,omitempty"`
	Overview     map[string]string `json:"overview,omitempty"`
	Email        string            `json:"email,omitempty"`
	Website      string            `json:"website,omitempty"`
	State        string            `json:"state,omitempty"`
	PublishOn    string            `json:"published_on,omitempty"`
	Links        map[string]string `json:"links,omitempty"`
	MinGas       decimal.Decimal   `json:"min_gas,omitempty"`
}

func GetTokenInfo(apiKey, tokenAddress string, service *Service) (ethplorer.Token, error) {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	infoKey := fmt.Sprintf(TOKEN_INFO_KEY, tokenAddress)
	buf, err := redis.Bytes(redisConn.Do("GET", infoKey))
	if err == nil {
		var token ethplorer.Token
		err = json.Unmarshal(buf, &token)
		if err == nil {
			return token, nil
		}
	}
	client := ethplorer.NewClient(apiKey)
	token, err := client.GetTokenInfo(tokenAddress)
	if token.Address == "" {
		return token, errors.New("Invalid Token")
	}
	js, err := json.Marshal(token)
	if err == nil {
		redisConn.Do("SETEX", infoKey, 24*60, string(js))
	}
	return token, nil
}

func GetETHPrice(service *Service, config Config) decimal.Decimal {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	rateStr, _ := redis.String(redisConn.Do("GET", ETH_PRICE_CACHE))
	rateValue, err := decimal.NewFromString(rateStr)
	if err == nil {
		return rateValue
	}
	coinbaseClient := coinbase.ApiKeyClient(config.CoinbaseAPI.Key, config.CoinbaseAPI.Secret)
	exchange, err := coinbaseClient.GetExchangeRate("eth", "usd")
	rateValue = decimal.NewFromFloat(exchange)
	if err == nil {
		redisConn.Do("SETEX", ETH_PRICE_CACHE, 2*60, rateValue.StringFixed(9))
	}
	return rateValue
}

func GetTMMRate(service *Service, config Config) (rate decimal.Decimal) {
	db := service.Db
	rows, _, err := db.Query("SELECT price FROM tmm.orderbook_trades ORDER BY id DESC LIMIT 1")
	if err != nil || len(rows) == 0 {
		return
	}
	rate, _ = decimal.NewFromString(rows[0].Str(0))
	return
}
