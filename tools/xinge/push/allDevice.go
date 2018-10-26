package push

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type AllDeviceRequest struct {
	xinge.BaseRequest
	DeviceToken  string `json:"device_token"`
	MessageType  uint   `json:"message_type"` //消息类型：1. 通知 2. 透传消息。iOS 平台请填 0
	Message      string `json:"message"`
	ExpireTime   int64  `json:"expire_time,omitempty"`   // 消息离线存储时间（单位为秒），最长存储时间 3 天。若设置为 0，则使用默认值（3 天）
	SendTime     string `json:"send_time,omitempty"`     // 指定推送时间，格式为 year-mon-day hour:min:sec， 若小于服务器当前时间，则会立即推送
	MultiPkg     uint   `json:"multi_pkg,omitempty"`     // 0表示按注册时提供的包名分发消息；1 表示按 access id 分发消息，所有以该 access id 成功注册推送的 App 均可收到消息。本字段对 iOS 平台无效
	Environment  uint   `json:"environment,omitempty"`   // 向 iOS 设备推送时必填，1 表示推送生产环境；2 表示推送开发环境。推送 Android 平台不填或填 0
	LoopTimes    uint   `json:"loop_times,omitempty"`    // 循环任务执行的次数，取值[1, 15]
	LoopInterval uint   `json:"loop_interval,omitempty"` // 循环任务的执行间隔，以天为单位，取值[1, 14]。loop_times 和 loop_interval 一起表示任务的生命周期，不可超过 14 天
}

func (this AllDeviceRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this AllDeviceRequest) Method() string {
	return "all_device"
}

func (this AllDeviceRequest) Class() string {
	return "push"
}
