package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
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
		var headerKeys = []string{"tmm-build", "tmm-ts", "tmm-nounce", "tmm-platform"}
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
			c.Next()
			return
		}
		appKey := c.Request.Header.Get("tmm-appkey")
		ts, _ := strconv.ParseInt(c.Request.Header.Get("tmm-ts"), 10, 64)
		if ts < time.Now().Add(-30*time.Second).Unix() || ts > time.Now().Add(30*time.Second).Unix() {
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "expired request"})
			return
		}
		sign := c.Request.Header.Get("tmm-sign")
		secret := GetAppSecret(appKey)
		if secret == "" {
			log.Error("empty secret")
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "invalid appkey"})
			return
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
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "invalid signature"})
			return
		}
		nounce := c.Request.Header.Get("tmm-nounce")
		redisConn := Service.Redis.Master.Get()
		defer redisConn.Close()
		nounceKey := fmt.Sprintf(REQUEST_NOUNCE_KEY, nounce)
		_, err := redis.String(redisConn.Do("GET", nounceKey))
		if err == nil {
			log.Warn("Duplicate Request Nounce: %s", nounce)
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "invalid nounce"})
			return
		}
		_, err = redisConn.Do("SETEX", nounceKey, 60, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Error(err.Error())
		}
		c.Next()
		return
	}
}
