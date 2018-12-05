package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
)

func InviteSummaryHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT COUNT(*) AS num, us.level FROM tmm.invite_codes AS ic LEFT JOIN tmm.user_settings AS us ON (us.user_id=ic.parent_id) WHERE ic.parent_id=%d`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var invites uint
	var creditLevel uint
	if len(rows) > 0 {
		invites = rows[0].Uint(0)
		creditLevel = rows[0].Uint(1)
	}
	rows, _, err = db.Query(`SELECT SUM(bonus), task_type FROM tmm.invite_bonus WHERE user_id=%d GROUP BY task_type`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var (
		points            decimal.Decimal
		friendsContribute decimal.Decimal
	)
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(0))
		if row.Uint(1) != 0 {
			friendsContribute.Add(bonus)
		}
		points = points.Add(bonus)
	}

	rows, _, err = db.Query(`SELECT invites FROM tmm.user_levels WHERE id>%d ORDER BY id ASC LIMIT 1`, creditLevel)
	if CheckErr(err, c) {
		return
	}
	var nextLevelInvites uint
	if len(rows) > 0 {
		nextLevel := rows[0].Uint(0)
		if nextLevel > invites {
			nextLevelInvites = nextLevel - invites
		} else {
			nextLevelInvites = 0
		}
	}
	rows, _, err = db.Query(`SELECT
                u.country_code,
                u.mobile,
                u.nickname,
                u.avatar,
                wx.nick,
                wx.avatar
            FROM tmm.invite_codes AS ic
            INNER JOIN ucoin.users AS u ON (u.id=ic.user_id)
            LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
            WHERE ic.parent_id=%d ORDER BY ic.user_id DESC LIMIT 10`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var users []common.User
	for _, row := range rows {
		u := common.User{
			CountryCode: row.Uint(0),
			Mobile:      utils.HideMobile(row.Str(1)),
			Nick:        row.Str(2),
			Avatar:      row.Str(3),
		}
		wxNick := row.Str(4)
		if wxNick != "" {
			u.Wechat = &common.Wechat{
				Nick:   wxNick,
				Avatar: row.Str(5),
			}
		}
		u.ShowName = u.GetShowName()
		u.Avatar = u.GetAvatar(Config.CDNUrl)
		users = append(users, u)
	}
	c.JSON(http.StatusOK, gin.H{"invites": invites, "points": points, "friends_contribute": friendsContribute, "users": users, "next_level_invites": nextLevelInvites})
}
