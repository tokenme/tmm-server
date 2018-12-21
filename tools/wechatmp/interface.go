package wechatmp

import (
    "github.com/tokenme/tmm/common"
	"github.com/gin-gonic/gin"
)

type Client struct {
	AppId       string
	AppSecret   string
    Service     *common.Service
    Context     *gin.Context
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

type OAuthAccessToken struct {
    AccessToken     string  `json:"access_token"`
    ExpiresIn       uint    `json:"expires_in"`
    RefreshToken    string  `json:"refresh_token"`
    Openid          string  `json:"openid"`
    Scope           string  `json:"scope"`
    ErrorCode       uint    `json:"errcode"`
    ErrorMsg        string  `json:"errmsg"`
}

func NewClient(appId string, appSecret string, service *common.Service, context *gin.Context) *Client {
	return &Client{
		AppId:      appId,
        AppSecret:  appSecret,
        Service:    service,
        Context:    context,
	}
}