package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/utils"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const REQUEST_NOUNCE_KEY = "ReqNounce-%s"

func ApiSignFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ApiCheckError(c); err != nil {
			c.Abort()
			c.JSON(http.StatusBadRequest, err)
		}
		c.Next()
		return
	}
}

func ApiSignPassFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ApiCheckError(c)
		c.Next()
		return
	}
}

func ApiCheckError(c *gin.Context) *APIError {
	var headerKeys = []string{"tmm-build", "tmm-ts", "tmm-nounce", "tmm-platform", "tmm-mac", "tmm-imei"}
	var requestParams = make(map[string]string)
	for _, k := range headerKeys {
		v := c.Request.Header.Get(k)
		if v == "" {
			continue
		}
		requestParams[k] = v
		c.Set(k, v)
	}
	if len(requestParams) == 0 {
		return nil
	}

	ipInfo, err := Service.Ip2Region.MemorySearch(ClientIP(c))
	if err != nil {
		log.Error(err.Error())
		return &APIError{
			Code: 400,
			Msg:  err.Error()}
	}
	if strings.Contains(ipInfo.ISP, "阿里云") {
		log.Warn("%s", ipInfo)
		return &APIError{
			Code: 400,
			Msg:  "invalid signature"}
	}
	appKey := c.Request.Header.Get("tmm-appkey")
	ts, _ := strconv.ParseInt(c.Request.Header.Get("tmm-ts"), 10, 64)
	if ts < time.Now().Add(-10*time.Minute).Unix() || ts > time.Now().Add(10*time.Minute).Unix() {
		return &APIError{
			Code: 400,
			Msg:  "Invalid timestamp, you may need to correct your system clock."}
	}
	sign := c.Request.Header.Get("tmm-sign")
	var secret string
	platform := c.Request.Header.Get("tmm-platform")
	if platform == "android" && appKey == Config.AndroidSig.Key {
		secret = Config.AndroidSig.Secret
	} else if appKey == Config.IOSSig.Key {
		secret = Config.IOSSig.Secret
	}
	if secret == "" {
		return &APIError{
			Code: 400,
			Msg:  "invalid appkey"}
	}
	queries := c.Request.URL.Query()
	for k, _ := range queries {
		requestParams[k] = fmt.Sprintf("%v", queries.Get(k))
	}
	if c.Request.Method == "POST" {
		postData := make(map[string]interface{})
		buf, _ := ioutil.ReadAll(c.Request.Body)
		if err := json.Unmarshal(buf, &postData); err != nil {
			log.Error(err.Error())
		} else {
			for k, v := range postData {
				var param string
				switch v.(type) {
				case bool:
					if v.(bool) {
						param = "1"
					} else {
						param = "0"
					}
				case float64:
					d := decimal.NewFromFloat(v.(float64))
					param = d.String()
				case float32:
					d := decimal.NewFromFloat(float64(v.(float32)))
					param = d.String()
				default:
					param = fmt.Sprintf("%v", v)
				}
				requestParams[k] = param
			}
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	}

	var signKeys []string
	for k, v := range requestParams {
		if v != "" {
			signKeys = append(signKeys, k)
		}
	}
	sort.Strings(signKeys)
	var signData []string
	for _, key := range signKeys {
		v := requestParams[key]
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s%s", key, v))
		}
	}
	urlStr := fmt.Sprintf("%s%s", Config.BaseUrl, c.Request.URL.Path)
	rawSign := fmt.Sprintf("%s%s%s%s", urlStr, appKey, strings.Join(signData, ""), secret)
	verifySign := utils.Md5(rawSign)
	if sign != verifySign {
		log.Info("invalid sign")
		log.Info(rawSign)
		log.Info(verifySign)
		return &APIError{
			Code: 400,
			Msg:  "invalid signature"}
	}
	nounce := c.Request.Header.Get("tmm-nounce")
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	nounceKey := fmt.Sprintf(REQUEST_NOUNCE_KEY, nounce)
	_, err = redis.String(redisConn.Do("GET", nounceKey))
	if err == nil {
		log.Warn("Duplicate Request Nounce: %s", nounce)
		return &APIError{
			Code: 400,
			Msg:  "invalid nounce"}
	}
	_, err = redisConn.Do("SETEX", nounceKey, 60, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error(err.Error())
	}
	return nil
}
