package wechatmp

import (
	// "fmt"
    "github.com/tokenme/tmm/common"
)

const (
	GATEWAY_PAY = ""
)

type Client struct {
	AppId       string
	AppSecret   string
    Service     *common.Service
}

type JSConfig struct {
    JSTicket    string  `json:"jsapi_ticket"`
	Nonce       string  `json:"noncestr"`
    Timestamp   string  `json:"timestamp"`
    Url         string  `json:"url"`
    Sign        string  `json:"sign"`
}

func (this JSConfig) toMap() map[string]string {
	return map[string]string{
		"jsapi_ticket": this.JSTicket,
        "noncestr":     this.Nonce,
        "timestamp":    this.Timestamp,
        "url":          this.Url,
        "sign":         this.Sign,
    }
}

type AccessTokenResponse struct {
    AccessToken string  `json:"access_token"`
    ExpiresIn   uint    `json:"expires_in"`
}

type JSTicketResponse struct {
    ErrorCode   uint    `json:"errcode"`
    ErrorMsg    string  `json:"errmsg"`
    Ticket      string  `json:"ticket"`
    ExpiresIn   uint    `json:"expires_in"`
}

func NewClient(appId string, appSecret string, service *common.Service) *Client {
	return &Client{
		AppId:      appId,
        AppSecret:  appSecret,
        Service:    service,
	}
}
