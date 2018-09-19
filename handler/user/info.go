package user

import (
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
                IFNULL(ic2.id, 0)
            FROM ucoin.users AS u
            LEFT JOIN tmm.invite_codes AS ic ON (ic.user_id = u.id)
            LEFT JOIN tmm.invite_codes AS ic2 ON (ic2.user_id = ic.parent_id)
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
			Id:          row.Uint64(0),
			CountryCode: row.Uint(1),
			Mobile:      row.Str(2),
			Nick:        row.Str(3),
			Avatar:      row.Str(4),
			Name:        row.Str(5),
			Salt:        row.Str(6),
			Password:    row.Str(7),
			Wallet:      row.Str(8),
			InviteCode:  tokenUtils.Token(row.Uint64(10)),
			InviterCode: tokenUtils.Token(row.Uint64(11)),
		}
		paymentPasswd := row.Str(9)
		if paymentPasswd != "" {
			user.CanPay = 1
		}
		user.ShowName = user.GetShowName()
		user.Avatar = user.GetAvatar(Config.CDNUrl)
	}
	c.JSON(http.StatusOK, user)
}
