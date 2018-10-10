package recaptcha

import (
	"encoding/json"
	//"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Struct for parsing json in google's response
type Response struct {
	Success    bool     `json:"success"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}

// url to post submitted re-captcha response to
const ENDPOINT = "https://recaptcha.net/recaptcha/api/siteverify"

// VerifyResponse is a method similar to `Verify`; but doesn't parse the form for you.  Useful if
// you're receiving the data as a JSON object from a javascript app or similar.
func Verify(secret string, hostname string, response string) Response {
	var gr Response
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.PostForm(ENDPOINT,
		url.Values{"secret": {secret}, "response": {response}})
	if err != nil {
		gr.ErrorCodes = append(gr.ErrorCodes, err.Error())
		return gr
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gr.ErrorCodes = append(gr.ErrorCodes, err.Error())
		return gr
	}

	err = json.Unmarshal(body, &gr)
	if err != nil {
		gr.ErrorCodes = append(gr.ErrorCodes, err.Error())
		return gr
	}
	if gr.Hostname != hostname {
		gr.ErrorCodes = append(gr.ErrorCodes, "invalid hostname")
		gr.Success = false
	}
	return gr
}
