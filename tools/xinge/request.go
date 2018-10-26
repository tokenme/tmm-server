package xinge

import (
	"fmt"
	"time"
)

type Request interface {
	HttpMethod() HttpMethod
	Method() string
	Class() string
}

type BaseRequest struct {
	AccessId  uint64  `json:"access_id"`
	CalType   CalType `json:"cal_type"`
	Timestamp int64   `json:"timestamp"`
	ValidTime uint    `json:"valid_time"`
}

func RequestURI(method string, class string) string {
	return fmt.Sprintf("http://%s", RequestGateway(method, class))
}

func RequestGateway(method string, class string) string {
	return fmt.Sprintf("%s/%s/%s/%s", DOMAIN, VERSION, class, method)
}

func (this *Client) DefaultBaseRequest() BaseRequest {
	return BaseRequest{
		AccessId:  this.AccessId,
		CalType:   RealtimeCalType,
		Timestamp: time.Now().Unix(),
		ValidTime: 600,
	}
}
