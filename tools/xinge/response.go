package xinge

type Response interface {
	Code() int
}

type BaseResponse struct {
	RetCode int    `json:"ret_code"`
	ErrMsg  string `json:"err_msg,omitempty"`
}

func (this BaseResponse) Code() int {
	return this.RetCode
}

type PushIdResponse struct {
	RetCode int          `json:"ret_code"`
	ErrMsg  string       `json:"err_msg,omitempty"`
	Result  PushIdResult `json:"result,omitempty"`
}

func (this PushIdResponse) Code() int {
	return this.RetCode
}

type PushIdResult struct {
	PushId string `json:"push_id"`
}

type DeviceNumResponse struct {
	RetCode int             `json:"ret_code"`
	ErrMsg  string          `json:"err_msg,omitempty"`
	Result  DeviceNumResult `json:"result,omitempty"`
}

func (this DeviceNumResponse) Code() int {
	return this.RetCode
}

type DeviceNumResult struct {
	DeviceNum uint64 `json:"device_num"`
}

type DeviceTokenResponse struct {
	RetCode int               `json:"ret_code"`
	ErrMsg  string            `json:"err_msg,omitempty"`
	Result  DeviceTokenResult `json:"result,omitempty"`
}

type DeviceTokenResult struct {
	IsReg         int   `json:"isReg"`         // 1为token已注册，0为未注册
	ConnTimestamp int64 `json:"connTimestamp"` // 最新活跃时间戳
	MsgsNum       uint  `json:"msgsNum"`       // 该应用的离线消息数
}

func (this DeviceTokenResponse) Code() int {
	return this.RetCode
}

type TokensResponse struct {
	RetCode int          `json:"ret_code"`
	ErrMsg  string       `json:"err_msg,omitempty"`
	Result  TokensResult `json:"result,omitempty"`
}

type TokensResult struct {
	Tokens []string `json:"tokens"`
}

func (this TokensResponse) Code() int {
	return this.RetCode
}

type TagsResponse struct {
	RetCode int        `json:"ret_code"`
	ErrMsg  string     `json:"err_msg,omitempty"`
	Result  TagsResult `json:"result,omitempty"`
}

type TagsResult struct {
	Total uint64   `json:"total"`
	Tags  []string `json:"tags"`
}

func (this TagsResponse) Code() int {
	return this.RetCode
}
