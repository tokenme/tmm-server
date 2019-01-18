package common

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
)

type Platform = string

const (
	IOS     Platform = "ios"
	ANDROID Platform = "android"
)

type Device struct {
	Id           string          `json:"id"`
	Name         string          `json:"name,omitempty"`
	Model        string          `json:"model,omitempty"`
	Platform     Platform        `json:"platform,omitempty"`
	Idfa         string          `json:"idfa,omitempty"`
	Mac          string          `json:"mac,omitempty"`
	Imei         string          `json:"imei,omitempty"`
	IsTablet     bool            `json:"is_tablet,omitempty"`
	IsEmulator   bool            `json:"is_emulator,omitempty"`
	TotalTs      uint64          `json:"total_ts,omitempty"`
	TotalApps    uint            `json:"total_apps,omitempty"`
	Points       decimal.Decimal `json:"points,omitempty"`
	Balance      decimal.Decimal `json:"balance,omitempty"`
	TMMBalance   decimal.Decimal `json:"-"`
	GrowthFactor decimal.Decimal `json:"gf,omitempty"`
	LastPingAt   string          `json:"lastping_at,omitempty"`
	InsertedAt   string          `json:"inserted_at,omitempty"`
	UpdatedAt    string          `json:"updated_at,omitempty"`
	Creative     *Creative       `json:"creative,omitempty"`
}

func (this Device) GetGrowthFactor(service *Service) (decimal.Decimal, error) {
	total, err := GetTotalAppTMMBalance(service)
	if err != nil {
		return decimal.New(0, 0), err
	}
	if total.IsZero() {
		return decimal.New(0, 0), err
	}
	return this.TMMBalance.DivRound(total, 9), nil
}

type DeviceRequest struct {
	Id              string   `json:"id,omitempty"`
	IsEmulator      bool     `json:"isEmulator,omitempty"`
	IsJailBrojen    bool     `json:"isJailBrojen,omitempty"`
	IsTablet        bool     `json:"isTablet,omitempty"`
	DeviceName      string   `json:"deviceName,omitempty"`
	Carrier         string   `json:"carrier,omitempty"`
	Country         string   `json:"country,omitempty"`
	Timezone        string   `json:"timezone,omitempty"`
	SystemVersion   string   `json:"systemVersion,omitempty"`
	AppName         string   `json:"appName,omitempty"`
	AppVersion      string   `json:"appVersion,omitempty"`
	AppBundleId     string   `json:"appBundleID,omitempty"`
	AppBuildVersion string   `json:"appBuildVersion,omitempty"`
	Ip              string   `json:"ip,omitempty"`
	Language        string   `json:"language,omitempty"`
	Idfa            string   `json:"idfa,omitempty"`
	Imei            string   `json:"imei,omitempty"`
	Mac             string   `json:"mac,omitempty"`
	OpenIDFA        string   `json:"openIDFA,omitempty"`
	DeviceType      string   `json:"deviceType,omitempty"`
	OsVersion       string   `json:"osVersion,omitempty"`
	Platform        Platform `json:"platform,omitempty"`
}

func UnmarshalDeviceRequest(val []byte) DeviceRequest {
	return DeviceRequest{
		Id:              jsoniter.Get(val, "id").ToString(),
		IsEmulator:      jsoniter.Get(val, "isEmulator").ToBool(),
		IsJailBrojen:    jsoniter.Get(val, "isJailBrojen").ToBool(),
		IsTablet:        jsoniter.Get(val, "isTablet").ToBool(),
		DeviceName:      jsoniter.Get(val, "deviceName").ToString(),
		Carrier:         jsoniter.Get(val, "carrier").ToString(),
		Country:         jsoniter.Get(val, "country").ToString(),
		Timezone:        jsoniter.Get(val, "timezone,omitempty").ToString(),
		SystemVersion:   jsoniter.Get(val, "systemVersion").ToString(),
		AppName:         jsoniter.Get(val, "appName").ToString(),
		AppVersion:      jsoniter.Get(val, "appVersion").ToString(),
		AppBundleId:     jsoniter.Get(val, "appBundleID").ToString(),
		AppBuildVersion: jsoniter.Get(val, "appBuildVersion").ToString(),
		Ip:              jsoniter.Get(val, "ip").ToString(),
		Language:        jsoniter.Get(val, "language").ToString(),
		Idfa:            jsoniter.Get(val, "idfa").ToString(),
		Imei:            jsoniter.Get(val, "imei").ToString(),
		Mac:             jsoniter.Get(val, "mac").ToString(),
		OpenIDFA:        jsoniter.Get(val, "openIDFA").ToString(),
		DeviceType:      jsoniter.Get(val, "deviceType").ToString(),
		OsVersion:       jsoniter.Get(val, "osVersion").ToString(),
		Platform:        jsoniter.Get(val, "platform").ToString(),
	}
}

func (this DeviceRequest) DeviceId() string {
	if len(this.Idfa) > 0 {
		this.Platform = IOS
		return utils.Sha1(this.Idfa)
	} else if len(this.Imei) > 0 {
		this.Platform = ANDROID
		str := this.Imei
		/*
		   if len(this.Mac) > 0 && this.Mac != "02:00:00:00:00:00" {
		       str = str + strings.Replace(this.Mac, ":", "", -1)
		   }
		*/
		return utils.Sha1(str)
	}
	return ""
}

func (this DeviceRequest) AppId() string {
	/*
		if this.Platform == IOS {
			return utils.Sha1(fmt.Sprintf("%s-%s", this.Platform, this.AppBundleId))
		}
	*/
	if len(this.AppBundleId) > 0 {
		return utils.Sha1(fmt.Sprintf("%s-%s", this.Platform, this.AppBundleId))
	}
	return ""
}

type PingRequest struct {
	Logs   string        `json:"logs,omitempty"`
	Ts     int64         `json:"duration,omitempty"`
	Device DeviceRequest `json:"device,omitempty"`
}

func UnmarshalPingRequest(val []byte) PingRequest {
	var req PingRequest
	req.Logs = jsoniter.Get(val, "logs").ToString()
	req.Ts = jsoniter.Get(val, "duration").ToInt64()
	j := jsoniter.Get(val, "device")
	req.Device = DeviceRequest{
		Id:              j.Get("id").ToString(),
		IsEmulator:      j.Get("isEmulator").ToBool(),
		IsJailBrojen:    j.Get("isJailBrojen").ToBool(),
		IsTablet:        j.Get("isTablet").ToBool(),
		DeviceName:      j.Get("deviceName").ToString(),
		Carrier:         j.Get("carrier").ToString(),
		Country:         j.Get("country").ToString(),
		Timezone:        j.Get("timezone,omitempty").ToString(),
		SystemVersion:   j.Get("systemVersion").ToString(),
		AppName:         j.Get("appName").ToString(),
		AppVersion:      j.Get("appVersion").ToString(),
		AppBundleId:     j.Get("appBundleID").ToString(),
		AppBuildVersion: j.Get("appBuildVersion").ToString(),
		Ip:              j.Get("ip").ToString(),
		Language:        j.Get("language").ToString(),
		Idfa:            j.Get("idfa").ToString(),
		Imei:            j.Get("imei").ToString(),
		Mac:             j.Get("mac").ToString(),
		OpenIDFA:        j.Get("openIDFA").ToString(),
		DeviceType:      j.Get("deviceType").ToString(),
		OsVersion:       j.Get("osVersion").ToString(),
		Platform:        j.Get("platform").ToString(),
	}
	return req
}

type PingCache struct {
	Logs string `json:"l,omitempty"`
	Ts   int64  `json:"ts,omitempty"`
	Cap  int64  `json:"cap,omitempty"`
}
