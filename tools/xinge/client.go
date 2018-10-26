package xinge

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	SecretKey string
	AccessId  uint64
}

func NewClient(accessId uint64, secretKey string) *Client {
	return &Client{AccessId: accessId, SecretKey: secretKey}
}

func (this *Client) Run(req Request, ret Response) error {
	values := this.genQueries(req)
	uri := RequestURI(req.Method(), req.Class())
	resp, err := http.Get(fmt.Sprintf("%s?%s", uri, values.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBytes, ret)
}

func (this *Client) genSign(req Request) string {
	mp := toMap(req)
	var signData []string
	for k, v := range mp {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case uint:
			val = strconv.FormatUint(uint64(v.(uint)), 10)
		case uint64:
			val = strconv.FormatUint(v.(uint64), 10)
		case int:
			val = strconv.FormatInt(int64(v.(int)), 10)
		case int64:
			val = strconv.FormatInt(v.(int64), 10)
		case interface{}:
			js, _ := json.Marshal(v.(interface{}))
			val = string(js)
		default:
			val = fmt.Sprintf("%v", v)
		}
		if val != "" {
			signData = append(signData, fmt.Sprintf("%s=%v", k, val))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "")
	signStr = fmt.Sprintf("%s%s%s%s", req.HttpMethod(), RequestGateway(req.Method(), req.Class()), signStr, this.SecretKey)
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return ""
	}
	signByte := c.Sum(nil)
	return strings.ToLower(fmt.Sprintf("%x", signByte))
}

func (this *Client) genQueries(req Request) url.Values {
	sign := this.genSign(req)
	mp := toMap(req)
	values := url.Values{}
	values.Add("sign", sign)
	for k, v := range mp {
		var val string
		switch v.(type) {
		case string:
			val = v.(string)
		case uint:
			val = strconv.FormatUint(uint64(v.(uint)), 10)
		case uint64:
			val = strconv.FormatUint(v.(uint64), 10)
		case int:
			val = strconv.FormatInt(int64(v.(int)), 10)
		case int64:
			val = strconv.FormatInt(v.(int64), 10)
		case interface{}:
			js, _ := json.Marshal(v.(interface{}))
			val = string(js)
		default:
			val = fmt.Sprintf("%v", v)
		}
		if val != "" {
			values.Add(k, val)
		}

	}
	return values
}

func toMap(obj interface{}) map[string]interface{} {
	js, _ := json.Marshal(obj)
	var ret map[string]interface{}
	json.Unmarshal(js, &ret)
	return ret
}
