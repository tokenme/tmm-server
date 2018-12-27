package auth

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/afs"
	"github.com/tokenme/tmm/tools/mobilecode"
	"github.com/tokenme/tmm/utils/twilio"
	"net/http"
	"strings"
)

const (
	SEND_CODE_KEY    = "SCK-%s"
	SEND_CODE_IP_KEY = "SCKIP-%s"
)

type SendRequest struct {
	Mobile     string `form:"mobile" json:"mobile" binding:"required"`
	Country    uint   `form:"country" json:"country" binding:"required"`
	AfsSession string `form:"afs_session" json:"afs_session" binding:"required"`
	AfsToken   string `form:"afs_token" json:"afs_token"`
	AfsSig     string `form:"afs_sig" json:"afs_sig"`
}

func SendHandler(c *gin.Context) {
	var req SendRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.AfsSession == "" || req.AfsToken != "" && req.AfsSig == "", "missing params", c) {
		return
	}

	mobile := strings.Replace(req.Mobile, " ", "", -1)
	mobile = strings.Replace(mobile, "-", "", -1)
	locale := "en"
	if req.Country == 86 {
		locale = "zh-CN"
	} else if req.Country == 886 || req.Country == 852 {
		locale = "zh-HK"
	} else if req.Country == 81 {
		locale = "ja"
	} else if req.Country == 82 {
		locale = "ko"
	}
	l := len(mobile)
	if Check(mobile == "" || l < 7 || l > 11, "invalid phone number", c) {
		return
	}

	if req.AfsSession != "" && req.AfsToken == "" {
		afsClient, err := afs.NewClientWithAccessKey(Config.Aliyun.RegionId, Config.Aliyun.AK, Config.Aliyun.AS)
		if CheckWithCode(err != nil, INVALID_CAPTCHA_ERROR, "Invalid captcha", c) {
			return
		}
		afsRequest := afs.CreateCreateAfsAppCheckRequest()
		afsRequest.Session = req.AfsSession
		afsResponse, err := afsClient.CreateAfsAppCheck(afsRequest)
		if CheckWithCode(err != nil || afsResponse.Data.SecondCheckResult != 1, INVALID_CAPTCHA_ERROR, "Invalid captcha", c) {
			return
		}
	} else if req.AfsSession != "" && req.AfsToken != "" && req.AfsSig != "" {
		afsClient, err := afs.NewClientWithAccessKey(Config.Aliyun.RegionId, Config.Aliyun.AK, Config.Aliyun.AS)
		if CheckWithCode(err != nil, INVALID_CAPTCHA_ERROR, "Invalid captcha", c) {
			return
		}
		afsRequest := afs.CreateAuthenticateSigRequest()
		afsRequest.SessionId = req.AfsSession
		afsRequest.Token = req.AfsToken
		afsRequest.Sig = req.AfsSig
		afsRequest.AppKey = Config.Aliyun.AfsAppKey
		afsRequest.Scene = "android"
		afsRequest.RemoteIp = ClientIP(c)
		afsResponse, err := afsClient.AuthenticateSig(afsRequest)
		if CheckWithCode(err != nil || afsResponse.Code != 100, INVALID_CAPTCHA_ERROR, "Invalid captcha", c) {
			return
		}
	}

	if req.Country == 86 {
		var deviceId string
		platform := c.GetString("tmm-platform")
		buildVersionStr := c.GetString("tmm-build")
		if platform == "android" {
			imei := c.GetString("tmm-imei")
			device := common.DeviceRequest{
				Platform: common.ANDROID,
				Imei:     imei,
			}
			deviceId = device.DeviceId()
			if Check(deviceId == "", "bad request", c) {
				return
			}
			redisConn := Service.Redis.Master.Get()
			defer redisConn.Close()
			sendCodeKey := fmt.Sprintf(SEND_CODE_KEY, deviceId)
			lastMobile, _ := redis.String(redisConn.Do("GET", sendCodeKey))
			_, err := redisConn.Do("SETEX", sendCodeKey, 60*5, mobile)
			if err != nil {
				log.Error(err.Error())
			}
			if Check(lastMobile != "" && lastMobile != mobile, "bad request", c) {
				log.Error("Auth Send last:%s, now:%s", lastMobile, mobile)
				return
			}
			sendCodeIPKey := fmt.Sprintf(SEND_CODE_IP_KEY, ClientIP(c))
			lastMobile, _ = redis.String(redisConn.Do("GET", sendCodeIPKey))
			_, err = redisConn.Do("SETEX", sendCodeIPKey, 60*5, mobile)
			if err != nil {
				log.Error(err.Error())
			}
			if Check(lastMobile != "" && lastMobile != mobile, "bad request", c) {
				log.Error("Auth Send last:%s, now:%s", lastMobile, mobile)
				return
			}
		}

		authClient := mobilecode.NewClient(Service, Config)
		_, err := authClient.Send(mobile)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			log.Error("Auth Send Failed: %s", err.Error())
			return
		}

		log.Warn("Auth Send: +%d-%s, platform: %s, build: %s, deviceId: %s", req.Country, mobile, platform, buildVersionStr, deviceId)
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	ret, err := twilio.AuthSend(Config.TwilioToken, mobile, req.Country, locale)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		log.Error("Auth Send Failed: %s", err.Error())
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		log.Error("Auth Send Failed: %s", ret.Message)
		return
	}
	log.Warn("Auth Send: +%d-%s", req.Country, mobile)
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
