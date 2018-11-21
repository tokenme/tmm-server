package auth

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils/twilio"
	"net/http"
	"strings"
)

type SendRequest struct {
	Mobile  string `form:"mobile" json:"mobile" binding:"required"`
	Country uint   `form:"country" json:"country" binding:"required"`
}

func SendHandler(c *gin.Context) {
	var req SendRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	mobile := strings.Replace(req.Mobile, " ", "", 0)
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
	ret, err := twilio.AuthSend(Config.TwilioToken, mobile, req.Country, locale)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
