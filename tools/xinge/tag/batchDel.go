package tag

import (
	"github.com/tokenme/tmm/tools/xinge"
)

type BatchDelRequest struct {
	xinge.BaseRequest
	TagTokenList [][]string `json:"tag_token_list"` // json 字符串，包含若干标签-token 对，后台将把每一对里面的 token 打上对应的标签。每次调用最多允许设置 20 对，每个对里面标签在前，token 在后。注意标签最长 50 字节，不可包含空格；真实 token 长度至少 40 字节。示例（其中 token 值仅为示意）： [[”tag 1”,”token 1”],[”tag 2”,”token 2”]]
}

func (this BatchDelRequest) HttpMethod() xinge.HttpMethod {
	return xinge.GET
}

func (this BatchDelRequest) Method() string {
	return "batch_del"
}

func (this BatchDelRequest) Class() string {
	return "tags"
}
