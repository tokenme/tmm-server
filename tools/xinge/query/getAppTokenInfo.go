package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type GetAppTokenInfoRequest struct {
	xinge.BaseRequest
	DeviceToken string `json:"device_token"`
}

func (this GetAppTokenInfoRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this GetAppTokenInfoRequest) Method() string {
	return "get_app_token_info"
}

func (this GetAppTokenInfoRequest) Class() string {
	return "application"
}
