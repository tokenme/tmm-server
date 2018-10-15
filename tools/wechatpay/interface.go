package wechatpay

import (
	"encoding/xml"
	"fmt"
)

type CheckNameType = string

const (
	NO_CHECK    CheckNameType = "NO_CHECK"
	FORCE_CHECK CheckNameType = "FORCE_CHECK"
	GATEWAY_PAY               = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

type Client struct {
	AppId   string
	MchId   string
	Key     string
	CertCrt string
	CertKey string
}

func NewClient(appId string, mchId string, key string, crt string, cliCrt string) *Client {
	return &Client{
		AppId:   appId,
		MchId:   mchId,
		Key:     key,
		CertCrt: crt,
		CertKey: cliCrt,
	}
}

type Request struct {
	XMLName     xml.Name      `xml:"xml"`
	AppId       string        `xml:"mch_appid"`
	MchId       string        `xml:"mchid"`
	DeviceInfo  string        `xml:"device_info,omitempty"`
	TradeNum    string        `xml:"partner_trade_no"`
	Amount      int64         `xml:"amount"`
	CallbackURL string        `xml:"-"`
	CheckName   CheckNameType `xml:"check_name"`
	Desc        string        `xml:"desc,omitempty"`
	OpenId      string        `xml:"openid"`
	Ip          string        `xml:"spbill_create_ip"`
	Nonce       string        `xml:"nonce_str"`
	Username    string        `xml:"re_user_name,omitempty"`
	Sign        string        `xml:"sign"`
}

func (this Request) toMap() map[string]string {
	return map[string]string{
		"mch_appid":        this.AppId,
		"mchid":            this.MchId,
		"partner_trade_no": this.TradeNum,
		"amount":           fmt.Sprintf("%d", this.Amount),
		"check_name":       this.CheckName,
		"desc":             this.Desc,
		"openid":           this.OpenId,
		"spbill_create_ip": this.Ip,
		"nonce_str":        this.Nonce,
		"re_user_name":     this.Username,
	}
}

type Response struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg,omitempty"`
	ErrCode     string `xml:"err_code,omitempty"`
	ErrCodeDesc string `xml:"err_code_des,omitempty"`
	AppId       string `xml:"mch_appid"`
	MchId       string `xml:"mchid"`
	DeviceInfo  string `xml:"device_info,omitempty"`
	TradeNum    string `xml:"partner_trade_no"`
	PaymentNum  string `xml:"payment_no,omitempty"`
	PaymentTime string `xml:"payment_time"`
}
