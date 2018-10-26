package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type GetMsgStatusRequest struct {
	xinge.BaseRequest
	PushIds []PushIdMap `map:"push_ids"`
}

type PushIdMap struct {
	PushId string `json:"push_id"`
}

func (this GetMsgStatusRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this GetMsgStatusRequest) Method() string {
	return "get_msg_status"
}

func (this GetMsgStatusRequest) Class() string {
	return "push"
}
