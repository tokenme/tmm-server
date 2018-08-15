package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
)

const (
	TOTAL_APP_TMM_BALANCE_KEY = "total_app_tmm_balance"
	APP_LOOKUP_KEY            = "applookup%s"
	APP_LOOKUP_GATEWAY        = "https://itunes.apple.com/lookup?bundleId=%s"
)

type App struct {
	Id           string          `json:"id"`
	Name         string          `json:"name,omitempty"`
	Version      string          `json:"version,omitempty"`
	Platform     Platform        `json:"platform,omitempty"`
	BundleId     string          `json:"bundle_id,omitempty"`
	StoreId      uint64          `json:"store_id,omitempty"`
	Icon         string          `json:"icon,omitempty"`
	Ts           uint64          `json:"ts,omitempty"`
	TMMBalance   decimal.Decimal `json:"-"`
	GrowthFactor decimal.Decimal `json:"gf,omitempty"`
	BuildVersion string          `json:"build_version,omitempty"`
	LastPingAt   string          `json:"lastping_at,omitempty"`
	InsertedAt   string          `json:"inserted_at,omitempty"`
	UpdatedAt    string          `json:"updated_at,omitempty"`
}

func GetTotalAppTMMBalance(service *Service) (total decimal.Decimal, err error) {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	totalTMMStr, _ := redis.String(redisConn.Do("GET", TOTAL_APP_TMM_BALANCE_KEY))
	total, err = decimal.NewFromString(totalTMMStr)
	if err == nil {
		return total, nil
	}
	db := service.Db
	rows, _, err := db.Query(`SELECT SUM(tmm) FROM tmm.apps WHERE is_active=1`)
	if err != nil {
		return total, err
	}
	total, err = decimal.NewFromString(rows[0].Str(0))
	if err != nil {
		return total, err
	}
	redisConn.Do("SETEX", TOTAL_APP_TMM_BALANCE_KEY, 24*60, total.StringFixed(9))
	return total, nil
}

func (this App) GetGrowthFactor(service *Service) (decimal.Decimal, error) {
	total, err := GetTotalAppTMMBalance(service)
	if err != nil {
		return decimal.New(0, 0), err
	}
	if total.IsZero() {
		return decimal.New(0, 0), err
	}
	return this.TMMBalance.DivRound(total, 9), nil
}

type AppLookUp struct {
	ResultCount uint              `json:"result_count"`
	Results     []AppLookUpResult `json:"results,omitempty"`
}

type AppLookUpResult struct {
	TrackId       uint64 `json:"trackId,omitempty"`
	ArtworkUrl512 string `json:"artworkUrl512,omitempty"`
}

func (this App) LookUp(service *Service) (lookupResult AppLookUpResult, err error) {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	redisKey := fmt.Sprintf(APP_LOOKUP_KEY, this.BundleId)
	lookupBytes, err := redis.Bytes(redisConn.Do("GET", redisKey))
	if err == nil {
		err = json.Unmarshal(lookupBytes, &lookupResult)
		if err == nil && lookupResult.TrackId > 0 {
			return lookupResult, err
		}
	}
	resp, err := http.Get(fmt.Sprintf(APP_LOOKUP_GATEWAY, this.BundleId))
	if err != nil {
		return lookupResult, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return lookupResult, err
	}
	js := bytes.Replace(body, []byte{'\n'}, []byte(""), -1)
	var lookup AppLookUp
	err = json.Unmarshal(js, &lookup)
	if err != nil {
		return lookupResult, err
	}
	if len(lookup.Results) == 0 {
		return lookupResult, errors.New("empty result")
	}
	lookupResult = lookup.Results[0]
	if lookupResult.TrackId == 0 {
		return lookupResult, errors.New("empty result")
	}
	js, err = json.Marshal(lookupResult)
	if err == nil {
		redisConn.Do("SETEX", redisKey, 60*24*7, string(js))
	}
	return lookupResult, nil
}
