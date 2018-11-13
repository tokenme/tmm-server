package ykt

import (
	"net/url"
	"strconv"
)

type GoodInfoResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data GoodInfoData `json:"data"`
}

func (this GoodInfoResponse) StatusCode() int {
	return this.Code
}

func (this GoodInfoResponse) Error() Error {
	return Error{Code: this.Code, Msg: this.Msg}
}

type GoodInfoData struct {
	Data Good `json:"data"`
}

type GoodInfoRequest struct {
	Id  uint64
	Uid uint64
}

func (this GoodInfoRequest) Run() (*GoodInfoResponse, error) {
	params := url.Values{}
	params.Add("id", strconv.FormatUint(this.Id, 10))
	params.Add("m_uid", strconv.FormatUint(this.Uid, 10))
	req := Request{
		Method: "g/item",
		Params: params,
	}
	res := &GoodInfoResponse{}
	err := Exec(req, res)
	return res, err
}
