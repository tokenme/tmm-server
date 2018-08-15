package user

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strings"
)

type UpdateRequest struct {
	Mobile        string `form:"mobile" json:"mobile"`
	CountryCode   uint   `form:"country_code" json:"country_code"`
	VerifyCode    string `form:"verify_code" json:"verify_code"`
	Password      string `form:"passwd" json:"passwd"`
	RePassword    string `form:"repasswd" json:"repasswd"`
	Realname      string `form:"realname" json:"realname"`
	PaymentPasswd string `form:"payment_passwd" json:"payment_passwd"`
}

func UpdateHandler(c *gin.Context) {
	var req UpdateRequest
	if CheckErr(c.Bind(&req), c) {
		log.Error("UNBINED")
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db

	var updateFields []string
	if req.Realname != "" {
		updateFields = append(updateFields, fmt.Sprintf("realname='%s'", db.Escape(req.Realname)))
	}
	var (
		passwd string
		mobile string
	)
	if req.Mobile != "" {
		if Check(req.Mobile == "" || req.CountryCode == 0 || req.VerifyCode == "" || req.Password == "" || req.RePassword == "", "missing params", c) {
			return
		}
		if Check(req.Password != req.RePassword, "repassword!=password", c) {
			return
		}
		passwdLength := len(req.Password)
		if Check(passwdLength < 8 || passwdLength > 64, "password length must between 8-32", c) {
			return
		}
		token, err := uuid.NewV4()
		if CheckErr(err, c) {
			return
		}
		salt := utils.Sha1(token.String())
		passwd = utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
		mobile = strings.Replace(req.Mobile, " ", "", 0)
		rows, _, err := db.Query(`SELECT 1 FROM ucoin.auth_verify_codes WHERE country_code=%d AND mobile='%s' AND code='%s' LIMIT 1`, req.CountryCode, db.Escape(mobile), db.Escape(req.VerifyCode))
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if Check(len(rows) == 0, "unverified phone number", c) {
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("country_code=%d, mobile='%s', salt='%s', passwd='%s'", req.CountryCode, db.Escape(mobile), db.Escape(salt), db.Escape(passwd)))
	} else if req.PaymentPasswd != "" {
		paymentPasswd := utils.Sha1(fmt.Sprintf("%s%s%s", user.Salt, req.PaymentPasswd, user.Salt))
		updateFields = append(updateFields, fmt.Sprintf("payment_passwd='%s'", db.Escape(paymentPasswd)))
	}
	if len(updateFields) == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}

	_, _, err := db.Query(`UPDATE ucoin.users SET %s WHERE id=%d LIMIT 1`, strings.Join(updateFields, ","), user.Id)
	if CheckErr(err, c) {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
