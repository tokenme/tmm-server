package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/afs"
	"github.com/tokenme/tmm/tools/mobilecode"
	"github.com/tokenme/tmm/utils/twilio"
	"net/http"
	"strings"
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
		authClient := mobilecode.NewClient(Service, Config)
		_, err := authClient.Send(req.Mobile)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			log.Error("Auth Send Failed: %s", err.Error())
			return
		}
		platform := c.GetString("tmm-platform")
		buildVersionStr := c.GetString("tmm-build")
		log.Warn("Auth Send: +%d-%s, platform: %s, build: %s", req.Country, mobile, platform, buildVersionStr)
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
