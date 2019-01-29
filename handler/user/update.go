package user

import (
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	commonutils "github.com/tokenme/tmm/utils"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type UpdateRequest struct {
	Mobile        string           `form:"mobile" json:"mobile"`
	CountryCode   uint             `form:"country_code" json:"country_code"`
	VerifyCode    string           `form:"verify_code" json:"verify_code"`
	Password      string           `form:"passwd" json:"passwd"`
	RePassword    string           `form:"repasswd" json:"repasswd"`
	Realname      string           `form:"realname" json:"realname"`
	PaymentPasswd string           `form:"payment_passwd" json:"payment_passwd"`
	InviterCode   tokenUtils.Token `form:"inviter_code" json:"inviter_code"`
	WxUnionId     string           `form:"wx_union_id" json:"wx_union_id"`
	WxOpenId      string           `form:"wx_open_id" json:"wx_open_id"`
	WxNick        string           `form:"wx_nick" json:"wx_nick"`
	WxAvatar      string           `form:"wx_avatar" json:"wx_avatar"`
	WxGender      int              `form:"wx_gender" json:"wx_gender"`
	WxToken       string           `form:"wx_token" json:"wx_token"`
	WxExpires     int64            `form:"wx_expires" json:"wx_expires"`
}

func UpdateHandler(c *gin.Context) {
	var req UpdateRequest
	if CheckWithCode(c.Bind(&req) != nil, BADREQUEST_ERROR, "invalid request", c) {
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
		salt := commonutils.Sha1(token.String())
		passwd = commonutils.Sha1(fmt.Sprintf("%s%s%s", salt, req.Password, salt))
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
		paymentPasswd := commonutils.Sha1(fmt.Sprintf("%s%s%s", user.Salt, req.PaymentPasswd, user.Salt))
		updateFields = append(updateFields, fmt.Sprintf("payment_passwd='%s'", db.Escape(paymentPasswd)))
	} else if req.InviterCode > 0 {
		_, ret, err := db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.parent_id=t2.user_id, t1.grand_id=t2.parent_id WHERE (t1.parent_id!=t2.user_id OR t1.grand_id!=t2.parent_id) AND t2.user_id!=t1.user_id AND t2.parent_id!= t1.user_id AND t1.parent_id=0 AND t2.id != t1.id AND t2.id=%d AND t1.user_id=%d`, req.InviterCode, user.Id)
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return
		}
		if CheckWithCode(ret.AffectedRows() == 0, INVALID_INVITE_CODE_ERROR, "invalid invite code", c) {
			return
		}

		err = _inviteBonus(c, user.Id, Service)
		if CheckErr(err, c) {
			return
		}
	} else if req.WxUnionId != "" {
		expires := time.Unix(req.WxExpires/1000, 0)
		rows, _, err := db.Query(`SELECT union_id FROM tmm.wx WHERE user_id=%d LIMIT 1`, user.Id)
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			unionId := rows[0].Str(0)
			if unionId == req.WxUnionId {
				_, _, err := db.Query(`UPDATE tmm.wx SET nick='%s', avatar='%s', gender=%d, access_token='%s', expires='%s' WHERE user_id=%d`, db.Escape(req.WxNick), db.Escape(req.WxAvatar), req.WxGender, db.Escape(req.WxToken), expires.Format("2006-01-02 15:04:05"), user.Id)
				if CheckErr(err, c) {
					return
				}
			}
			c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
			return
		}
		_, _, err = db.Query(`INSERT INTO tmm.wx (user_id, union_id, open_id, nick, avatar, gender, access_token, expires) VALUES (%d, '%s', '%s', '%s', '%s', %d, '%s', '%s') ON DUPLICATE KEY UPDATE union_id=VALUES(union_id), open_id=VALUES(open_id), nick=VALUES(nick), avatar=VALUES(avatar), gender=VALUES(gender), access_token=VALUES(access_token), expires=VALUES(expires)`, user.Id, db.Escape(req.WxUnionId), db.Escape(req.WxOpenId), db.Escape(req.WxNick), db.Escape(req.WxAvatar), req.WxGender, db.Escape(req.WxToken), expires.Format("2006-01-02 15:04:05"))
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return
		}

		err = _inviteBonus(c, user.Id, Service)
		if CheckErr(err, c) {
			return
		}
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

func _inviteBonus(c *gin.Context, userId uint64, service *common.Service) error {
	db := service.Db
	{ // Check if have been gave bonus
		rows, _, err := db.Query(`SELECT IFNULL(ib.user_id, 0) FROM tmm.invite_codes AS ic INNER JOIN tmm.wx ON (wx.user_id=ic.user_id) LEFT JOIN tmm.invite_bonus AS ib ON (ib.from_user_id=ic.user_id AND ib.task_type=0) WHERE ic.user_id=%d AND ic.parent_id>0 LIMIT 1`, userId)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		ibUser := rows[0].Uint64(0)
		if ibUser > 0 {
			return nil
		}
	}

	query := `SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id=d.user_id)
WHERE ic.user_id = %d AND (IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1)
ORDER BY d.lastping_at DESC LIMIT 1`
	rows, _, err := db.Query(query, userId)
	if err != nil {
		raven.CaptureError(err, nil)
		return err
	}
	if len(rows) == 0 {
		return errors.New("not found")
	}
	inviterPointBonus := decimal.New(int64(Config.InviterBonus), 0)
	inviterDeviceId := rows[0].Str(0)
	inviterUserId := rows[0].Uint64(1)
	pointsPerTs, _ := common.GetPointsPerTs(Service)
	inviterTs := inviterPointBonus.Div(pointsPerTs)
	forexRate := forex.Rate(Service, "USD", "CNY")
	tx, tmm, err := _transferToken(userId, forexRate, c)
	if err != nil {
		log.Error("Bonus Transfer failed")
		raven.CaptureError(err, nil)
		return err
	}
	//log.Warn("Inviter bonus: %s, inviter:%d", inviterPointBonus.String(), inviterUserId)
	{ // UPDATE inviter bonus
		_, _, err := db.Query(`UPDATE tmm.devices AS d2 SET d2.points = d2.points + %s, d2.total_ts = d2.total_ts + %d WHERE d2.id='%s'`, inviterPointBonus.String(), inviterTs.IntPart(), db.Escape(inviterDeviceId))
		if err != nil {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return err
		}
		_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, tmm, tmm_tx) VALUES (%d, %d, 0, %s, '%s'), (%d, %d, %s, 0, '')`, userId, userId, tmm.String(), db.Escape(tx), inviterUserId, userId, inviterPointBonus.String())
		if err != nil {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
		}
	}
	{ // UPDATE user grand id
		_, _, err := db.Query(`UPDATE tmm.invite_codes AS t1, tmm.invite_codes AS t2 SET t1.grand_id=t2.parent_id WHERE t2.user_id=t1.parent_id AND t2.parent_id!=t1.user_id AND t2.user_id=%d`, userId)
		if err != nil {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return err
		}
	}

	{ // UPDATE inviter credit level
		_, _, err := db.Query(`INSERT INTO user_settings (user_id, level)
(
SELECT
i.parent_id, ul.id
FROM tmm.user_levels AS ul
INNER JOIN (
    SELECT ic.parent_id, COUNT(DISTINCT ic.user_id) AS invites
    FROM tmm.invite_codes AS ic
    INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
    LEFT JOIN tmm.user_settings AS us ON (us.user_id=ic.user_id)
    WHERE ic.parent_id=%d AND (IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1)
) AS i ON (i.invites >= ul.invites AND parent_id IS NOT NULL)
ORDER BY ul.id DESC LIMIT 1
) ON DUPLICATE KEY UPDATE level=VALUES(level)`, inviterUserId)
		if CheckErr(err, c) {
			log.Error(err.Error())
			raven.CaptureError(err, nil)
			return err
		}
	}
	return nil
}

func _transferToken(userId uint64, forexRate decimal.Decimal, c *gin.Context) (receipt string, tokenAmount decimal.Decimal, err error) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT u.wallet_addr FROM ucoin.users AS u LEFT JOIN tmm.user_settings AS us ON (us.user_id=u.id) WHERE u.id=%d AND (IFNULL(us.blocked, 0)=0 OR us.block_whitelist=1)`, userId)
	if err != nil {
		return receipt, tokenAmount, err
	}
	if len(rows) == 0 {
		return receipt, tokenAmount, errors.New(fmt.Sprintf("not found: %d", userId))
	}
	userWallet := rows[0].Str(0)
	tokenPrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	tokenPriceCny := tokenPrice.Mul(forexRate)
	tokenAmount = decimal.New(3, 0).Div(tokenPriceCny)

	token, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if err != nil {
		return receipt, tokenAmount, err
	}
	tokenDecimal, err := utils.TokenDecimal(token)
	if err != nil {
		return receipt, tokenAmount, err
	}
	tmmInt := tokenAmount.Mul(decimal.New(1, int32(tokenDecimal)))
	amount, ok := new(big.Int).SetString(tmmInt.Floor().String(), 10)
	if !ok {
		return receipt, tokenAmount, nil
	}

	agentPrivKey, err := commonutils.AddressDecrypt(Config.TMMAgentWallet.Data, Config.TMMAgentWallet.Salt, Config.TMMAgentWallet.Key)
	if err != nil {
		return receipt, tokenAmount, err
	}
	agentPubKey, err := eth.AddressFromHexPrivateKey(agentPrivKey)
	if err != nil {
		return receipt, tokenAmount, err
	}
	tokenBalance, err := utils.TokenBalanceOf(token, agentPubKey)
	if err != nil {
		return receipt, tokenAmount, err
	}
	if amount.Cmp(tokenBalance) == 1 {
		return receipt, tokenAmount, nil
	}
	transactor := eth.TransactorAccount(agentPrivKey)
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	nonce, err := eth.Nonce(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		return receipt, tokenAmount, err
	}
	gasPrice, err := Service.Geth.SuggestGasPrice(c)
	if err == nil && gasPrice.Cmp(eth.MinGas) == -1 {
		gasPrice = eth.MinGas
	} else {
		gasPrice = nil
	}
	transactorOpts := eth.TransactorOptions{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: 210000,
	}
	eth.TransactorUpdate(transactor, transactorOpts, c)
	tx, err := utils.Transfer(token, transactor, userWallet, amount)
	if err != nil {
		return receipt, tokenAmount, err
	}
	err = eth.NonceIncr(c, Service.Geth, Service.Redis.Master, agentPubKey, Config.Geth)
	if err != nil {
		log.Error(err.Error())
	}
	receipt = tx.Hash().Hex()
	return receipt, tokenAmount, nil
}
