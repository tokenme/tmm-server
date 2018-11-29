package handler

import (
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/articlesuggest"
	"github.com/tokenme/tmm/tools/blowup"
	"github.com/tokenme/tmm/utils"
	//"github.com/tokenme/ucoin/tools/sqs"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"runtime"
	"path"
)

var (
	Service       *common.Service
	Config        common.Config
	GlobalLock    *sync.Mutex
	BlowupService *blowup.Server
	SuggestEngine *articlesuggest.Engine
	ExitCh        chan struct{}
	//Queues  map[string]sqs.Queue
)

func InitHandler(s *common.Service, c common.Config) {
	Service = s
	Config = c
	GlobalLock = new(sync.Mutex)
	BlowupService = blowup.NewServer(s, c)
	SuggestEngine = articlesuggest.NewEngine(s, c)
	//Queues = queues
	raven.SetDSN(Config.SentryDSN)
	ExitCh = make(chan struct{}, 1)
}

func Start() {
	BlowupService.Start()
	go SuggestEngine.Start()
}

func Close() {
	BlowupService.Stop()
	SuggestEngine.Stop()
}

type APIResponse struct {
	Msg string `json:"message,omitempty"`
}

type ErrorCode = int

const (
	BADREQUEST_ERROR            ErrorCode = 400
	INTERNAL_ERROR              ErrorCode = 500
	NOTFOUND_ERROR              ErrorCode = 404
	UNAUTHORIZED_ERROR          ErrorCode = 401
	FEATURE_NOT_AVAILABLE_ERROR ErrorCode = 402
	INVALID_PASSWD_ERROR        ErrorCode = 409
	INVALID_CAPTCHA_ERROR       ErrorCode = 408
	INVALID_PASSWD_LENGTH       ErrorCode = 407
	DUPLICATE_USER_ERROR        ErrorCode = 202
	UNACTIVATED_USER_ERROR      ErrorCode = 502
	NOT_ENOUGH_TOKEN_ERROR      ErrorCode = 600
	DAILY_BONUS_COMMITTED_ERROR ErrorCode = 601
	NOT_ENOUGH_POINTS_ERROR     ErrorCode = 700
	INVALID_MIN_POINTS_ERROR    ErrorCode = 701
	INVALID_MIN_TOKEN_ERROR     ErrorCode = 702
	WECHAT_UNAUTHORIZED_ERROR   ErrorCode = 703
	WECHAT_PAYMENT_ERROR        ErrorCode = 704
	WECHAT_OPENID_ERROR         ErrorCode = 705
	NOT_ENOUGH_ETH_ERROR        ErrorCode = 800
	INVALID_INVITE_CODE_ERROR   ErrorCode = 1000
	MAX_BIND_DEVICE_ERROR       ErrorCode = 1100
	OTHER_BIND_DEVICE_ERROR     ErrorCode = 1101
	INVALID_CDP_VENDOR_ERROR    ErrorCode = 1200
	BLOWUP_ESCAPE_LATE_ERROR    ErrorCode = 1300
	BLOWUP_ESCAPE_EARLY_ERROR   ErrorCode = 1301
)

type APIError struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"message,omitempty"`
}

func (this APIError) Error() string {
	return fmt.Sprintf("CODE:%d, MSG:%s", this.Code, this.Msg)
}

func Check(flag bool, err string, c *gin.Context) (ret bool) {
	ret = flag
	if ret {
		_, file, line, _ := runtime.Caller(1)
		log.Error("[%s:%d]: %s", path.Base(file), line, err)
		c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err})
	}
	return
}

func CheckErr(err error, c *gin.Context) (ret bool) {
	ret = err != nil
	if ret {
		_, file, line, _ := runtime.Caller(1)
		log.Error("[%s:%d]: %s", path.Base(file), line, err.Error())
		if _, ok := err.(APIError); ok {
			c.JSON(http.StatusOK, err)
		} else {
			c.JSON(http.StatusOK, APIError{Code: BADREQUEST_ERROR, Msg: err.Error()})
		}
	}
	return
}

func CheckWithCode(flag bool, code ErrorCode, err string, c *gin.Context) (ret bool) {
	ret = flag
	if ret {
		log.Error(err)
		c.JSON(http.StatusOK, APIError{Code: code, Msg: err})
	}
	return
}

func Uint64Value(val string, defaultVal uint64) (uint64, error) {
	if val == "" {
		return defaultVal, nil
	}

	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func Uint64NonZero(val string, err string) (uint64, error) {
	if val == "" {
		return 0, errors.New(err)
	}

	i, e := strconv.ParseUint(val, 10, 64)
	if e != nil {
		return 0, e
	}

	return i, nil
}

func ClientIP(c *gin.Context) string {
	if values, _ := c.Request.Header["X-Forwarded-For"]; len(values) > 0 {
		clientIP := values[0]
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	if values, _ := c.Request.Header["X-Real-Ip"]; len(values) > 0 {
		clientIP := strings.TrimSpace(values[0])
		if len(clientIP) > 0 {
			return clientIP
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func IP2Long(IpStr string) (int64, error) {
	bits := strings.Split(IpStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}

type APIRequest struct {
	APPKey  string `form:"k" json:"k" binding:"required"`
	Version string `form:"v" json:"v" binding:"required"`
	Ts      int64  `form:"t" json:"t" binding:"required"`
	Rand    string `form:"r" json:"r" binding:"required"`
	Payload string `form:"p" json:"p" binding:"required"`
	Sign    string `form:"s" json:"s" binding:"required"`
}

func ApiCheckFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req APIRequest
		if err := c.Bind(&req); err != nil {
			log.Error(err.Error())
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": err.Error()})
			return
		}
		if req.Ts < time.Now().Add(-10 * time.Minute).Unix() || req.Ts > time.Now().Add(10 * time.Minute).Unix() {
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "expired request"})
			return
		}
		secret := GetAppSecret(req.APPKey)
		if secret == "" {
			log.Error("empty secret")
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "invalid appkey"})
			return
		}
		if !verifySign(secret, req) {
			c.Abort()
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    401,
				"message": "invalid signature"})
			return
		}
		c.Set("Request", req)
		c.Set("Secret", secret)
		c.Next()
		return
	}
}

func GetAppSecret(appKey string) string {
	redisMasterConn := Service.Redis.Master.Get()
	defer redisMasterConn.Close()
	redisKey := fmt.Sprintf("tmm-app-%s", appKey)
	secret, _ := redis.String(redisMasterConn.Do("GET", redisKey))
	if secret == "" {
		db := Service.Db
		rows, _, err := db.Query(`SELECT secret FROM tmm.appkeys WHERE appkey='%s' AND active=1 LIMIT 1`, db.Escape(appKey))
		if err != nil || len(rows) == 0 {
			if err != nil {
				log.Error(err.Error())
			}
			return ""
		}
		secret = rows[0].Str(0)
		redisMasterConn.Do("SETEX", redisKey, 60*60*24, secret)
	}
	return secret
}

func verifySign(secret string, req APIRequest) bool {

	reqStr := fmt.Sprintf("k=%s&p=%s&r=%s&t=%d&v=%s", req.APPKey, req.Payload, req.Rand, req.Ts, req.Version)
	sign := utils.Sha1(fmt.Sprintf("%s%s%s", secret, reqStr, secret))
	return sign == strings.ToLower(req.Sign)
}
