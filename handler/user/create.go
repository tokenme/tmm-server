package user

import (
	"bytes"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/govatar"
	"image/png"
	//"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/tokenme/tmm/coins/eth"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/qiniu"
	"github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"github.com/tokenme/tmm/utils/twilio"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateRequest struct {
	Mobile      string `form:"mobile" json:"mobile"`
	CountryCode uint   `form:"country_code" json:"country_code"`
	VerifyCode  string `form:"verify_code" json:"verify_code"`
	Password    string `form:"passwd" json:"passwd"`
	RePassword  string `form:"repasswd" json:"repasswd"`
}

func CreateHandler(c *gin.Context) {
	var req CreateRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	createByMobile(c, req)
}

func createByMobile(c *gin.Context, req CreateRequest) {
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
	passwd := utils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
	mobile := strings.Replace(req.Mobile, " ", "", 0)

	ret, err := twilio.AuthVerification(Config.TwilioToken, mobile, req.CountryCode, req.VerifyCode)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	if Check(!ret.Success, ret.Message, c) {
		return
	}

	privateKey, _, err := eth.GenerateAccount()
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	walletSalt, wallet, err := utils.AddressEncrypt(privateKey, Config.TokenSalt)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	pubKey, err := eth.AddressFromHexPrivateKey(privateKey)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	gender := govatar.MALE
	maleOrFemale := utils.RangeRandUint64(0, 1)
	if maleOrFemale == 1 {
		gender = govatar.FEMALE
	}
	avatarImg, err := govatar.GenerateFromUsername(gender, pubKey)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	avatarBuf := new(bytes.Buffer)
	err = png.Encode(avatarBuf, avatarImg)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	avatar, _, err := qiniu.Upload(c, Config.Qiniu, fmt.Sprintf("%s/%s", Config.Qiniu.AvatarPath, pubKey), timestamp, avatarBuf.Bytes())
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}

	db := Service.Db
	_, userRet, err := db.Query(`INSERT INTO ucoin.users (country_code, mobile, passwd, avatar, salt, active, wallet, wallet_salt, wallet_addr) VALUES (%d, '%s', '%s', '%s', '%s', 1, '%s', '%s', '%s')`, req.CountryCode, db.Escape(mobile), db.Escape(passwd), db.Escape(avatar), db.Escape(salt), db.Escape(wallet), db.Escape(walletSalt), db.Escape(pubKey))
	if err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY {
		c.JSON(http.StatusOK, APIResponse{Msg: "account already exists"})
		return
	}
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	userId := userRet.InsertId()

	for {
		nickname := tokenUtils.New()
		rows, _, err := db.Query(`UPDATE ucoin.users SET nickname='%s' WHERE id=%d LIMIT 1`, db.Escape(nickname.Encode()), userId)
		if err != nil {
			continue
		}
		if len(rows) == 0 {
			break
		}
	}

	for {
		inviteCode := tokenUtils.New()
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.invite_codes (id, user_id) VALUES (%d, %d)`, inviteCode, userId)
		if err != nil {
			continue
		}
		break
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
