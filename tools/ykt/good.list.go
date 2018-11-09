package ykt

import (
	"fmt"
	"github.com/shopspring/decimal"
	"net/url"
)

type GoodListResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data GoodListPage `json:"data"`
}

func (this GoodListResponse) StatusCode() int {
	return this.Code
}

func (this GoodListResponse) Error() Error {
	return Error{Code: this.Code, Msg: this.Msg}
}

type GoodListPage struct {
	Page     uint   `json:"curr_page"`
	PageSize uint   `json:"page_size"`
	Total    uint64 `json:"total_rows"`
	Items    []Good `json:"data"`
}

type Good struct {
	Id               uint64          `json:"id"`
	SkuId            uint64          `json:"sku_id"`
	WareId           uint64          `json:"ware_id"`
	AccountId        uint64          `json:"account_id,omitempty"`
	OriPrice         decimal.Decimal `json:"ori_price,omitempty"`
	Price            decimal.Decimal `json:"price,omitempty"`
	CommissionPrice  decimal.Decimal `json:"commision_price,omitempty"`
	PurchaseWithdraw decimal.Decimal `json:"purchase_withdraw,omitempty"`
	Name             string          `json:"goods_name,omitempty"`
	Pic              string          `json:"goods_pic,omitempty"`
	ShareLink        string          `json:"share_link,omitempty"`
}

type GoodListRequest struct {
	Source   uint
	Page     uint
	PageSize uint
}

func (this GoodListRequest) Run() (*GoodListResponse, error) {
	params := url.Values{}
	params.Add("source", fmt.Sprintf("%d", this.Source))
	params.Add("page", fmt.Sprintf("%d", this.Page))
	params.Add("page_size", fmt.Sprintf("%d", this.PageSize))
	req := Request{
		Method: "g/list",
		Params: params,
	}
	res := &GoodListResponse{}
	err := Exec(req, res)
	return res, err
}
