package user

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func InviteSummaryHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT COUNT(*) FROM tmm.invite_codes WHERE parent_id=%d`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var invites uint
	if len(rows) > 0 {
		invites = rows[0].Uint(0)
	}
	rows, _, err = db.Query(`SELECT SUM(bonus) FROM tmm.invite_bonus WHERE user_id=%d`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var points decimal.Decimal
	if len(rows) > 0 {
		points, _ = decimal.NewFromString(rows[0].Str(0))
	}
	c.JSON(http.StatusOK, gin.H{"invites": invites, "points": points})
}
