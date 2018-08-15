package common

import (
	"fmt"
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
	IsTablet     bool            `json:"is_tablet,omitempty"`
	TotalTs      uint64          `json:"total_ts,omitempty"`
	TotalApps    uint            `json:"total_apps,omitempty"`
	Points       decimal.Decimal `json:"points,omitempty"`
	Balance      decimal.Decimal `json:"balance,omitempty"`
	TMMBalance   decimal.Decimal `json:"-"`
	GrowthFactor decimal.Decimal `json:"gf,omitempty"`
	LastPingAt   string          `json:"lastping_at,omitempty"`
	InsertedAt   string          `json:"inserted_at,omitempty"`
	UpdatedAt    string          `json:"updated_at,omitempty"`
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
	OpenIDFA        string   `json:"openIDFA,omitempty"`
	DeviceType      string   `json:"deviceType,omitempty"`
	OsVersion       string   `json:"osVersion,omitempty"`
	Platform        Platform `json:"platform,omitempty"`
}

func (this DeviceRequest) DeviceId() string {
	if this.Platform == IOS {
		return utils.Sha1(this.Idfa)
	}
	return ""
}

func (this DeviceRequest) AppId() string {
	if this.Platform == IOS {
		return utils.Sha1(fmt.Sprintf("%s-%s", this.Platform, this.AppBundleId))
	}
	return ""
}

type PingRequest struct {
	Ts     int64         `json:"duration,omitempty"`
	Device DeviceRequest `json:"device,omitempty"`
}
