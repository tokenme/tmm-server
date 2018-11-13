package ykt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	GATEWAY = "https://jkgj-isv.isvjcloud.com/rest/m"
)

type Request struct {
	Method string
	Params url.Values
}

func (this Request) Uri() string {
	return fmt.Sprintf("%s/%s?%s", GATEWAY, this.Method, this.Params.Encode())
}

type Response interface {
	StatusCode() int
	Error() Error
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (this Error) Error() string {
	return this.Msg
}

func Exec(req Request, res Response) error {
	resp, err := http.DefaultClient.Get(req.Uri())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	if res.StatusCode() != 0 {
		return res.Error()
	}
	return nil
}
