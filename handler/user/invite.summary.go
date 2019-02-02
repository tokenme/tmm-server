package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
)

type InviteSummaryRequest struct {
	WithUserList bool `json:"with_user_list" form:"with_user_list"`
}

type InviteSummaryResponse struct {
	Invites           uint            `json:"invites"`
    FamilyInvites     uint            `json:"family_invites"`
	Points            decimal.Decimal `json:"points"`
	FriendsContribute decimal.Decimal `json:"friends_contribute"`
	Users             []common.User   `json:"users,omitempty"`
	NextLevelInvites  uint            `json:"next_level_invites"`
	InviteBonusRate   decimal.Decimal `json:"invite_bonus_rate"`
	InviteBonus       decimal.Decimal `json:"invite_bonus"`
	InviterBonus      decimal.Decimal `json:"inviter_bonus"`
	InviteCashBonus   decimal.Decimal `json:"invite_cash_bonus"`
	InviterCashBonus  decimal.Decimal `json:"inviter_cash_bonus"`
}

func InviteSummaryHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req InviteSummaryRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	query := `SELECT COUNT(*) AS num, us.level
FROM tmm.invite_codes AS ic
INNER JOIN tmm.wx AS wx ON (wx.user_id=ic.user_id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id=ic.parent_id)
LEFT JOIN tmm.user_settings AS us2 ON (us2.user_id=ic.user_id)
WHERE ic.parent_id=%d AND (IFNULL(us2.blocked, 0)=0 OR us2.block_whitelist=1)`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}
	var invites uint
	var creditLevel uint
	if len(rows) > 0 {
		invites = rows[0].Uint(0)
		creditLevel = rows[0].Uint(1)
	}

    rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.invite_codes AS ic LEFT JOIN tmm.invite_submissions AS iss ON (ic.id=iss.code) WHERE ic.user_id=%d AND iss.is_family=1`, user.Id)
    if CheckErr(err, c) {
        return
    }
    var familyInvites uint
    if len(rows) > 0 {
        familyInvites = rows[0].Uint(0)
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
			friendsContribute = friendsContribute.Add(bonus)
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
	var users []common.User
	if req.WithUserList {
		query := `SELECT
                u.country_code,
                u.mobile,
                u.nickname,
                u.avatar,
                wx.nick,
                wx.avatar
            FROM tmm.invite_codes AS ic
            INNER JOIN ucoin.users AS u ON (u.id=ic.user_id)
            LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
            WHERE ic.parent_id=%d ORDER BY ic.user_id DESC LIMIT 10`
		rows, _, err = db.Query(query, user.Id)
		if CheckErr(err, c) {
			return
		}

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
	}
	summary := InviteSummaryResponse{
		Invites:           invites,
        FamilyInvites:     familyInvites,
		Points:            points,
		FriendsContribute: friendsContribute,
		Users:             users,
		NextLevelInvites:  nextLevelInvites,
		InviteBonusRate:   decimal.NewFromFloat(Config.InviteBonusRate * 100),
		InviteBonus:       decimal.New(int64(Config.InviteBonus), 0),
		InviterBonus:      decimal.New(int64(Config.InviterBonus), 0),
		InviteCashBonus:   decimal.New(int64(Config.InviteCashBonus), 0),
		InviterCashBonus:  decimal.New(int64(Config.InviterCashBonus), 0),
	}
	c.JSON(http.StatusOK, summary)
}
