package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type QueryAppTagsRequest struct {
	xinge.BaseRequest
	Start uint `json:"start,omitempty"`
	Limit uint `json:"limit,omitempty"`
}

func (this QueryAppTagsRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this QueryAppTagsRequest) Method() string {
	return "query_app_tags"
}

func (this QueryAppTagsRequest) Class() string {
	return "tags"
}
