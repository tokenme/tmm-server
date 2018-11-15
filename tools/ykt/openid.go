package ykt

import (
	"net/url"
)

type OpenIdResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data OpenIdData `json:"data"`
}

func (this OpenIdResponse) StatusCode() int {
	return this.Code
}

func (this OpenIdResponse) Error() Error {
	return Error{Code: this.Code, Msg: this.Msg}
}

type OpenIdData struct {
	OpenId string `json:"openid"`
}

type OpenIdRequest struct {
	UnionId string
}

func (this OpenIdRequest) Run() (*OpenIdResponse, error) {
	params := url.Values{}
	params.Add("unionid", this.UnionId)
	req := Request{
		Method: "u/getopenid",
		Params: params,
	}
	res := &OpenIdResponse{}
	err := Exec(req, res)
	return res, err
}
