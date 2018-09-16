package ethplorer

import (
	"encoding/json"
	"fmt"
)

func (this *Client) GetTokenInfo(tokenAddress string) (token Token, err error) {
	uri := fmt.Sprintf("/getTokenInfo/%s", tokenAddress)
	data, err := this.Exec(uri, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return
	}
	mp := make(map[string]interface{})
	err = json.Unmarshal(data, &mp)
	if err != nil {
		return
	}
	return token, nil
}
