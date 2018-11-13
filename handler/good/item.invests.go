package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strconv"
)

func ItemInvestsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	itemId, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.ParseUint(c.Param("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(c.Param("pageSize"), 10, 64)
	if page == 0 {
		page = 1
	}
	defaultPageSize := uint64(DEFAULT_PAGE_SIZE)
	if pageSize == 0 || pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT u.id, u.country_code, u.mobile, wx.nick, wx.avatar, gi.points FROM tmm.good_invests AS gi LEFT JOIN ucoin.users AS u ON (u.id=gi.user_id) LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id) WHERE gi.good_id=%d AND gi.redeem_status IN (0, 1) ORDER BY gi.points DESC LIMIT %d, %d`, itemId, (page-1)*pageSize, pageSize)
	if CheckErr(err, c) {
		return
	}
	var invests []common.GoodInvest
	for _, row := range rows {
		mobile := row.Str(2)
		if row.Uint64(0) != user.Id {
			mobile = utils.HideMobile(mobile)
		}
		u := common.User{
			CountryCode: row.Uint(1),
			Mobile:      mobile,
			Wechat: &common.Wechat{
				Nick:   row.Str(3),
				Avatar: row.Str(4),
			},
		}
		u.ShowName = u.GetShowName()
		u.Avatar = u.GetAvatar(Config.CDNUrl)
		points, _ := decimal.NewFromString(row.Str(5))
		inv := common.GoodInvest{
			Nick:   u.ShowName,
			Avatar: u.Avatar,
			Points: points,
		}
		invests = append(invests, inv)
	}
	c.JSON(http.StatusOK, invests)
}
