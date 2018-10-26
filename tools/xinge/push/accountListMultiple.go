package push

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type AccountListMultipleRequest struct {
	xinge.BaseRequest
	AccountList []string `json:"account_list"` // Json 数组格式，每个元素是一个 account，string 类型，单次发送 account 不超过 100 个。例：[“account 1”,”account 2”,”account 3”]
	PushId      string   `json:"push_id"`      // 创建批量推送消息，接口的返回值中的 push_id
}

func (this AccountListMultipleRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this AccountListMultipleRequest) Method() string {
	return "account_list_multiple"
}

func (this AccountListMultipleRequest) Class() string {
	return "push"
}
