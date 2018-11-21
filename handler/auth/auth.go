package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"encoding/base64"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/o1egl/govatar"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/middlewares/jwt"
	"github.com/tokenme/tmm/tools/qiniu"
	"github.com/tokenme/tmm/tools/recaptcha"
	"github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"image/png"
	"strconv"
	"time"
)

type Biometric struct {
	Passwd string `json:"passwd"`
	Ts     int64  `json:"ts"`
}

var AuthenticatorFunc = func(loginInfo jwt.Login, c *gin.Context) (string, int, bool) {
	db := Service.Db
	var where string
	var loginPasswd string
	IsAdmin :=VerifyAdmin(loginInfo.Mobile,loginInfo.CountryCode)
	if !IsAdmin {
		if loginInfo.CountryCode > 0 && loginInfo.Mobile != "" && loginInfo.Password != "" && loginInfo.Captcha != "" {
			where = fmt.Sprintf("u.country_code=%d AND u.mobile='%s'", loginInfo.CountryCode, db.Escape(loginInfo.Mobile))
		} else {
			log.Error("missing params")
			return loginInfo.Mobile, BADREQUEST_ERROR, false
		}
		if loginInfo.Biometric {
			secret := GetAppSecret(loginInfo.Captcha)
			if secret == "" {
				log.Error("invalid captcha")
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			decrepted, err := utils.DesDecrypt(loginInfo.Password, []byte(secret))
			if err != nil {
				log.Error(err.Error())
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			decodedStr, err := base64.StdEncoding.DecodeString(string(decrepted))
			if err != nil {
				log.Error(err.Error())
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			var biometric Biometric
			err = json.Unmarshal([]byte(decodedStr), &biometric)
			if err != nil {
				log.Error(err.Error())
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			if biometric.Ts < time.Now().Add(-10 * time.Minute).Unix() || biometric.Ts > time.Now().Add(10 * time.Minute).Unix() {
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			loginPasswd = biometric.Passwd
		} else {
			captchaRes := recaptcha.Verify(Config.ReCaptcha.Secret, Config.ReCaptcha.Hostname, loginInfo.Captcha)
			if !captchaRes.Success {
				log.Error("invalid captcha")
				return loginInfo.Mobile, INVALID_CAPTCHA_ERROR, false
			}
			loginPasswd = loginInfo.Password
		}
	}else{
		where = fmt.Sprintf("u.country_code=%d AND u.mobile='%s'", loginInfo.CountryCode, db.Escape(loginInfo.Mobile))
		loginPasswd = loginInfo.Password
		}
	query := `SELECT
                u.id,
                u.country_code,
                u.mobile,
                u.nickname,
                u.avatar,
                u.realname,
                u.salt,
                u.passwd,
                u.wallet_addr,
                u.payment_passwd,
                IFNULL(ic.id, 0),
                IFNULL(ic2.id, 0),
                IFNULL(us.exchange_enabled, 0),
                wx.union_id,
                wx.nick,
                wx.avatar,
                wx.gender,
                wx.access_token,
                wx.expires
            FROM ucoin.users AS u
            LEFT JOIN tmm.invite_codes AS ic ON (ic.user_id = u.id)
            LEFT JOIN tmm.invite_codes AS ic2 ON (ic2.user_id = ic.parent_id)
            LEFT JOIN tmm.user_settings AS us ON (us.user_id = u.id)
            LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
            WHERE %s
            AND active = 1
            LIMIT 1`
	rows, _, err := db.Query(query, where)
	if err != nil || len(rows) == 0 {
		if err != nil {
			log.Error(err.Error())
		}
		if len(rows) == 0 {
			return loginInfo.Mobile, UNACTIVATED_USER_ERROR, false
		}
		return loginInfo.Mobile, INTERNAL_ERROR, false
	}
	row := rows[0]
	user := common.User{
		Id:              row.Uint64(0),
		CountryCode:     row.Uint(1),
		Mobile:          row.Str(2),
		Nick:            row.Str(3),
		Avatar:          row.Str(4),
		Name:            row.Str(5),
		Salt:            row.Str(6),
		Password:        row.Str(7),
		Wallet:          row.Str(8),
		InviteCode:      tokenUtils.Token(row.Uint64(10)),
		InviterCode:     tokenUtils.Token(row.Uint64(11)),
		ExchangeEnabled: row.Int(12) == 1 || row.Uint(1) != 86,
	}
	paymentPasswd := row.Str(9)
	if paymentPasswd != "" {
		user.CanPay = 1
	}
	wxUnionId := row.Str(13)
	if wxUnionId != "" {
		wechat := &common.Wechat{
			UnionId:     wxUnionId,
			Nick:        row.Str(14),
			Avatar:      row.Str(15),
			Gender:      row.Uint(16),
			AccessToken: row.Str(17),
			Expires:     row.ForceLocaltime(18),
		}
		user.Wechat = wechat
	}
	if user.Nick == "" {
		for {
			nickname := tokenUtils.New()
			user.Nick = nickname.Encode()
			rows, _, err := db.Query(`UPDATE ucoin.users SET nickname='%s' WHERE id=%d LIMIT 1`, db.Escape(user.Nick), user.Id)
			if err != nil {
				continue
			}
			if len(rows) == 0 {
				break
			}
		}
	}
	if user.InviteCode == 0 {
		for {
			inviteCode := tokenUtils.New()
			_, _, err := db.Query(`INSERT IGNORE INTO tmm.invite_codes (id, user_id) VALUES (%d, %d)`, inviteCode, user.Id)
			if err != nil {
				continue
			}
			user.InviteCode = inviteCode
			break
		}
	}
	if user.Avatar == "" {
		gender := govatar.MALE
		maleOrFemale := utils.RangeRandUint64(0, 1)
		if maleOrFemale == 1 {
			gender = govatar.FEMALE
		}
		avatarImg, err := govatar.GenerateFromUsername(gender, user.Wallet)
		if err != nil {
			log.Error(err.Error())
			return loginInfo.Mobile, INTERNAL_ERROR, false
		}
		avatarBuf := new(bytes.Buffer)
		err = png.Encode(avatarBuf, avatarImg)
		if err != nil {
			log.Error(err.Error())
			return loginInfo.Mobile, INTERNAL_ERROR, false
		}

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		avatar, _, err := qiniu.Upload(c, Config.Qiniu, fmt.Sprintf("%s/%s", Config.Qiniu.AvatarPath, user.Wallet), timestamp, avatarBuf.Bytes())
		if err != nil {
			log.Error(err.Error())
			return loginInfo.Mobile, INTERNAL_ERROR, false
		}
		user.Avatar = avatar
		db.Query(`UPDATE ucoin.users SET avatar='%s' WHERE id=%d LIMIT 1`, db.Escape(user.Avatar), user.Id)
	}
	user.ShowName = user.GetShowName()
	user.Avatar = user.GetAvatar(Config.CDNUrl)
	c.Set("USER", user)
	passwdSha1 := utils.Sha1(fmt.Sprintf("%s%s%s", user.Salt, loginPasswd, user.Salt))
	js, err := json.Marshal(user)
	if err != nil {
		log.Error(err.Error())
		return loginInfo.Mobile, INTERNAL_ERROR, false
	}
	if passwdSha1 != user.Password {
		return string(js), INVALID_PASSWD_ERROR, false
	}
	return string(js), 200, passwdSha1 == user.Password
}

var AuthorizatorFunc = func(data string, c *gin.Context) bool {
	var user common.User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return false
	}
	db := Service.Db
	query := `SELECT 1 FROM ucoin.users WHERE id=%d AND active = 1 LIMIT 1`
	rows, _, err := db.Query(query, user.Id)
	if err != nil || len(rows) == 0 {
		if err != nil {
			raven.CaptureError(err, nil)
		}
		return false
	}
	c.Set("USER", user)
	return true
}
