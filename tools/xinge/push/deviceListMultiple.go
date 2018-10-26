package push

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type DeviceListMultipleRequest struct {
	xinge.BaseRequest
	DeviceList []string `json:"device_list"` // Json 数组格式，每个元素是一个 token，string 类型，单次发送 token 不超过 1000 个。例：[“token 1”,”token 2”,”token 3”]
	PushId     string   `json:"push_id"`     // 创建批量推送消息，接口的返回值中的 push_id
}

func (this DeviceListMultipleRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this DeviceListMultipleRequest) Method() string {
	return "device_list_multiple"
}

func (this DeviceListMultipleRequest) Class() string {
	return "push"
}
