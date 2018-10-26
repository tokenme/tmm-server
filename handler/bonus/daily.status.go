package bonus

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func DailyStatusHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT days FROM tmm.daily_bonus_logs WHERE user_id=%d LIMIT 1`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var days int
	if len(rows) > 0 {
		days = rows[0].Int(0)
	}
	if days == 7 {
		days = 0
	}
	c.JSON(http.StatusOK, gin.H{"days": days})
}
