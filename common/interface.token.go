package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/tools/ethplorer-api"
)

const (
	TOKEN_INFO_KEY = "etherplorer-token-%s"
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
