package wechatmp

import (
	"crypto/sha1"
    "github.com/garyburd/redigo/redis"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
    "time"
	"sort"
	"strings"
    "strconv"
    "errors"
)

const (
    WX_ACCESS_TOKEN         = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
    WX_JS_TICKET            = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
    WX_AT_EXPIRES           = 7000
    WX_JS_TICKET_EXPIRES    = 7000
)

func (this *Client) GetJSConfig(url string) (jsConfig JSConfig, err error) {
    jsTicket, err := this.GetJSTicket()
    if err != nil || len(jsTicket) == 0 {
        return jsConfig, err
    }
    now := time.Now()
    jsConfig = JSConfig {
        Nonce:      getNonceStr(),
        Timestamp:  strconv.Itoa(int(now.Unix())),
        Url:        url,
        JSTicket:   jsTicket,
    }
	err = jsConfig.GenSign()
    return jsConfig, err
}

func (this *Client) GetJSTicket() (jsTicket string, err error) {
    redisConn := this.Service.Redis.Master.Get()
    defer redisConn.Close()
    var respBytes []byte
    cacheKey := "wx-mp-js-ticket"
    jsTicket, _ = redis.String(redisConn.Do("GET", cacheKey))
    if len(jsTicket) == 0 {
        at, err := this.GetAccessToken()
        if err != nil {
            return jsTicket, err
        }
        request, err := http.NewRequest("GET", fmt.Sprintf(WX_JS_TICKET, at), nil)
        if err != nil {
            return jsTicket, err
        }
        client := http.Client{}
        resp, err := client.Do(request)
        if err != nil {
            return jsTicket, err
        }
        defer resp.Body.Close()
        respBytes, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            return jsTicket, err
        }
        var res JSTicketResponse
	    err = json.Unmarshal(respBytes, &res)
        if err != nil {
            return jsTicket, err
        }
        if len(res.Ticket) > 0 {
            jsTicket = res.Ticket
            redisConn.Do("SETEX", cacheKey, WX_JS_TICKET_EXPIRES, jsTicket)
        } else if res.ErrorCode > 0 {
            return jsTicket, errors.New(fmt.Sprintf("Get wechat js ticket error(%d): %s", res.ErrorCode, res.ErrorMsg))
        }
    }
	return jsTicket, err
}

func (this *Client) GetAccessToken() (accessToken string, err error) {
    redisConn := this.Service.Redis.Master.Get()
    defer redisConn.Close()
    var respBytes []byte
    cacheKey := "wx-mp-access-token"
    accessToken, _ = redis.String(redisConn.Do("GET", cacheKey))
    if len(accessToken) == 0 {
        request, err := http.NewRequest("GET", fmt.Sprintf(WX_ACCESS_TOKEN, this.AppId, this.AppSecret), nil)
        if err != nil {
            return accessToken, err
        }
        client := http.Client{}
        resp, err := client.Do(request)
        if err != nil {
            return accessToken, err
        }
        defer resp.Body.Close()
        respBytes, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            return accessToken, err
        }
        var res AccessTokenResponse
	    err = json.Unmarshal(respBytes, &res)
        if err != nil {
            return accessToken, err
        }
        if len(res.AccessToken) > 0 {
            accessToken = res.AccessToken
            redisConn.Do("SETEX", cacheKey, WX_AT_EXPIRES, accessToken)
        }
    }
	return accessToken, err
}

func getNonceStr() (nonceStr string) {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(chars) - 1)
		nonceStr += chars[idx : idx+1]
	}
	return
}

func (param *JSConfig) GenSign() error {
	m := param.toMap()
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	c := sha1.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return err
	}
	signByte := c.Sum(nil)
	param.Sign = strings.ToUpper(fmt.Sprintf("%x", signByte))
	return nil

}
