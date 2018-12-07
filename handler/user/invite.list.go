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

const DEFAULT_PAGE_SIZE = 10

type InviteListRequest struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

func InviteListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req InviteListRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}

	db := Service.Db

	query := `SELECT
    u.country_code,
    u.mobile,
    u.nickname,
    u.avatar,
    wx.nick,
    wx.avatar,
    IF(ic.id IS NULL, false, true) is_parent,
    tb.bonus
FROM
(SELECT SUM(ib.bonus) AS bonus, ib.from_user_id AS from_user_id, ib.user_id AS user_id
FROM tmm.invite_bonus AS ib
WHERE ib.user_id=%d
GROUP BY ib.from_user_id) AS tb
INNER JOIN ucoin.users AS u ON (u.id=tb.from_user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
LEFT JOIN tmm.invite_codes AS ic ON (ic.user_id=tb.from_user_id AND ic.parent_id=tb.user_id)
ORDER BY tb.bonus DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, (req.Page-1)*req.PageSize, req.PageSize)
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
		u.DirectFriend = row.Bool(6)
		u.Contribute, _ = decimal.NewFromString(row.Str(7))
		users = append(users, u)
	}

	c.JSON(http.StatusOK, users)
}
