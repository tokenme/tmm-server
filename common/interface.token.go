package common

import (
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/lucazulian/cryptocomparego"
	"github.com/lucazulian/cryptocomparego/context"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/ethplorer-api"
	"strings"
)

const (
	TOKEN_INFO_KEY  = "etherplorer-token-%s"
	ETH_PRICE_CACHE = "eth-usd"
	TMM_PRICE_CACHE = "tmm-usd-%d"
)

type PriceType = uint8

const (
	CrowsalePrice PriceType = 1
	MarketPrice   PriceType = 2
	RecyclePrice  PriceType = 3
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
	ctx := context.TODO()
	client := cryptocomparego.NewClient(nil)
	req := cryptocomparego.NewPriceRequest("ETH", []string{"USD"})
	priceList, _, err := client.Price.List(ctx, req)
	if err != nil {
		return rateValue
	}
	for _, price := range priceList {
		if strings.ToLower(price.Name) == "usd" {
			rateValue = decimal.NewFromFloat(price.Value)
			redisConn.Do("SETEX", ETH_PRICE_CACHE, 2*60, rateValue.StringFixed(9))
			break
		}
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

func GetTMMPrice(service *Service, config Config, priceType PriceType) (price decimal.Decimal) {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	cacheKey := fmt.Sprintf(TMM_PRICE_CACHE, priceType)
	priceStr, _ := redis.String(redisConn.Do("GET", cacheKey))
	priceValue, err := decimal.NewFromString(priceStr)
	if err == nil {
		return priceValue
	}
	var (
		salePrice    decimal.Decimal
		recyclePrice decimal.Decimal
	)
	db := service.Db
	rows, _, err := db.Query("SELECT price, recycle_price FROM tmm.erc20 WHERE address='%s' LIMIT 1", config.TMMTokenAddress)
	if err != nil {
		return
	}
	if len(rows) > 0 {
		salePrice, _ = decimal.NewFromString(rows[0].Str(0))
		recyclePrice, _ = decimal.NewFromString(rows[0].Str(1))
	}
	/*
		ethPrice := GetETHPrice(service, config)
		tmmRate := GetTMMRate(service, config)
		tmmPrice := tmmRate.Mul(ethPrice)
		if (priceType == CrowsalePrice || priceType == MarketPrice) && tmmPrice.GreaterThan(price) || priceType == RecyclePrice && tmmPrice.LessThan(price) {
			price = tmmPrice
		}
	*/
	if priceType == CrowsalePrice || priceType == MarketPrice {
		price = salePrice
	} else if priceType == RecyclePrice {
		price = recyclePrice
	}
	if price.GreaterThan(decimal.Zero) {
		redisConn.Do("SETEX", cacheKey, 2*60, price.StringFixed(9))
	}
	return
}

func GetPointPrice(service *Service, config Config) (price decimal.Decimal) {
	exchangeRate, _, err := GetExchangeRate(config, service)
	if err != nil {
		return
	}
	tmmPrice := GetTMMPrice(service, config, RecyclePrice)
	return tmmPrice.Mul(exchangeRate.Rate)
}
