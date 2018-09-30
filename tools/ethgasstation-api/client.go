package ethgasstation

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
)

const GATEWAY = "https://ethgasstation.info/json/ethgasAPI.json"

type GasRes struct {
	Average     decimal.Decimal `json:"average"`
	FastestWait decimal.Decimal `json:"fastestWait"`
	FastWait    decimal.Decimal `json:"fastWait"`
	Fast        decimal.Decimal `json:"fast"`
	SafeLowWait decimal.Decimal `json:"safeLowWait"`
	BlockNum    uint64          `json:"blockNum"`
	AvgWait     decimal.Decimal `json:"avgWait"`
	BlockTime   int64           `json:"blockTime"`
	Speed       decimal.Decimal `json:"speed"`
	Fastest     decimal.Decimal `json:"fastest"`
	SafeLow     decimal.Decimal `json:"safeLow"`
}

func Gas() (gas GasRes, err error) {
	resp, err := http.Get(GATEWAY)
	if err != nil {
		return gas, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return gas, err
	}
	err = json.Unmarshal(body, &gas)
	if err != nil {
		return gas, err
	}
	return gas, nil
}
