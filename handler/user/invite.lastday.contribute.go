package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"net/http"
)

func InviteLastdayContributeHandler(c *gin.Context) {
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
	rows, _, err := db.Query(`SELECT SUM(bonus) FROM tmm.invite_bonus WHERE user_id=%d AND user_id!=from_user_id AND inserted_at>=DATE_SUB(NOW(), INTERVAL 1 DAY)`, user.Id)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, gin.H{"contribute": 0, "cny": 0})
		return
	}
	row := rows[0]
	friendsContribute, _ := decimal.NewFromString(row.Str(0))
	pointPrice := common.GetPointPrice(Service, Config)
	currency := c.Query("currency")
	if currency == "" {
		currency = "USD"
	}
	if currency != "USD" {
		rate := forex.Rate(Service, "USD", currency)
		pointPrice = pointPrice.Mul(rate)
	}
	cny := friendsContribute.Mul(pointPrice)
	c.JSON(http.StatusOK, gin.H{"contribute": friendsContribute, "cny": cny})
}
