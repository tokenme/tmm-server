package user

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/mobilecode"
	"github.com/tokenme/tmm/utils"
	"github.com/tokenme/tmm/utils/twilio"
	"net/http"
	"strings"
)

type ResetPasswordRequest struct {
	Mobile      string `form:"mobile" json:"mobile" binding:"required"`
	CountryCode uint   `form:"country_code" json:"country_code" binding:"required"`
	VerifyCode  string `form:"verify_code" json:"verify_code" binding:"required"`
	Password    string `form:"passwd" json:"passwd" binding:"required"`
	RePassword  string `form:"repasswd" json:"repasswd" binding:"required"`
}

func ResetPasswordHandler(c *gin.Context) {
	var req ResetPasswordRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	passwdLength := len(req.Password)
	if Check(passwdLength < 8 || passwdLength > 64, "password length must between 8-32", c) {
		return
	}
	if Check(req.Password != req.RePassword, "repassword!=password", c) {
		return
	}
	mobile := strings.Replace(req.Mobile, " ", "", -1)
	mobile = strings.Replace(mobile, "-", "", -1)
	if req.CountryCode == 86 {
		authClient := mobilecode.NewClient(Service, Config)
		ret := authClient.Verify(mobile, req.VerifyCode)
		if Check(!ret.Success, ret.Message, c) {
			log.Error("Auth Send Failed: %s", ret.Message)
			return
		}
	} else {
		ret, err := twilio.AuthVerification(Config.TwilioToken, mobile, req.CountryCode, req.VerifyCode)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if Check(!ret.Success, ret.Message, c) {
			return
		}
	}

	db := Service.Db
	rows, _, err := db.Query(`SELECT id, salt FROM ucoin.users WHERE country_code=%d AND mobile='%s' LIMIT 1`, req.CountryCode, db.Escape(mobile))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(len(rows) == 0, "user doesn't exists", c) {
		return
	}
	userId := rows[0].Uint64(0)
	salt := rows[0].Str(1)
	passwd := utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
	_, _, err = db.Query(`UPDATE ucoin.users SET passwd='%s' WHERE id=%d`, db.Escape(passwd), userId)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
