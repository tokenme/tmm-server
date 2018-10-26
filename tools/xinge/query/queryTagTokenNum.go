package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type QueryTagTokenNumRequest struct {
	xinge.BaseRequest
	Tag string `json:"tag"`
}

func (this QueryTagTokenNumRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this QueryTagTokenNumRequest) Method() string {
	return "query_tag_token_num"
}

func (this QueryTagTokenNumRequest) Class() string {
	return "tags"
}
