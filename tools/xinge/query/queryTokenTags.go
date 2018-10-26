package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type QueryTokenTagsRequest struct {
	xinge.BaseRequest
	DeviceToken string `json:"device_token"`
}

func (this QueryTokenTagsRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this QueryTokenTagsRequest) Method() string {
	return "query_token_tags"
}

func (this QueryTokenTagsRequest) Class() string {
	return "tags"
}
