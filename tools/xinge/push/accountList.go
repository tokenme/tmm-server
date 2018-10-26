package push

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type AccountListRequest struct {
	xinge.BaseRequest
	AccountList []string `json:"account_list"` // Json 数组格式，每个元素是一个 account，string 类型，单次发送 account 不超过 100 个。例：[“account 1”,”account 2”,”account 3”]
	MessageType uint     `json:"message_type"` //消息类型：1. 通知 2. 透传消息。iOS 平台请填 0
	Message     string   `json:"message"`
	ExpireTime  int64    `json:"expire_time,omitempty"` // 消息离线存储时间（单位为秒），最长存储时间 3 天。若设置为 0，则使用默认值（3 天）
	MultiPkg    uint     `json:"multi_pkg,omitempty"`   // 0表示按注册时提供的包名分发消息；1 表示按 access id 分发消息，所有以该 access id 成功注册推送的 App 均可收到消息。本字段对 iOS 平台无效
	Environment uint     `json:"environment,omitempty"` // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台不填或填 0
}

func (this AccountListRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this AccountListRequest) Method() string {
	return "account_list"
}

func (this AccountListRequest) Class() string {
	return "push"
}
