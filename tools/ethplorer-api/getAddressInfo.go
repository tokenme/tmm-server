package ethplorer

import (
	"encoding/json"
	"fmt"
)

func (this *Client) GetAddressInfo(address string, tokenAddress string) (addressInfo Address, err error) {
	uri := fmt.Sprintf("/getAddressInfo/%s", address)
	var params map[string]string
	if tokenAddress != "" {
		params["token"] = tokenAddress
	}
	data, err := this.Exec(uri, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &addressInfo)
	if err != nil {
		return
	}
	return addressInfo, nil
}
