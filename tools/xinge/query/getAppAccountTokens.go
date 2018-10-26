package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type GetAppAccountTokensRequest struct {
	xinge.BaseRequest
	Account string `json:"account"`
}

func (this GetAppAccountTokensRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this GetAppAccountTokensRequest) Method() string {
	return "get_app_account_tokens"
}

func (this GetAppAccountTokensRequest) Class() string {
	return "application"
}
