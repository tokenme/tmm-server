package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"net/http"
)

func InfoGetHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	if c.Query("refresh") != "" {
		db := Service.Db
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
                IFNULL(us.level, 0),
                ul.name,
                ul.enname,
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
            LEFT JOIN tmm.user_levels AS ul ON (ul.id=IFNULL(us.level, 0))
            LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
            WHERE u.id = %d
            AND active = 1
            LIMIT 1`
		rows, _, err := db.Query(query, user.Id)
		if CheckErr(err, c) {
			raven.CaptureError(err, nil)
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		row := rows[0]
		user = common.User{
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
			Level: common.CreditLevel{
				Id:     row.Uint(13),
				Name:   row.Str(14),
				Enname: row.Str(15),
			},
		}
		paymentPasswd := row.Str(9)
		if paymentPasswd != "" {
			user.CanPay = 1
		}
		wxUnionId := row.Str(16)
		if wxUnionId != "" {
			wechat := &common.Wechat{
				UnionId:     wxUnionId,
				Nick:        row.Str(17),
				Avatar:      row.Str(18),
				Gender:      row.Uint(19),
				AccessToken: row.Str(20),
				Expires:     row.ForceLocaltime(21),
			}
			user.Wechat = wechat
			user.WxBinded = true
		}
		user.ShowName = user.GetShowName()
		user.Avatar = user.GetAvatar(Config.CDNUrl)
	}
	c.JSON(http.StatusOK, user)
}
