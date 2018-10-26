package query

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type GetAppDeviceNumRequest struct {
	xinge.BaseRequest
}

func (this GetAppDeviceNumRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this GetAppDeviceNumRequest) Method() string {
	return "get_app_device_num"
}

func (this GetAppDeviceNumRequest) Class() string {
	return "application"
}
