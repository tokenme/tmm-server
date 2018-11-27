package user

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func CreditLevelsHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query("SELECT id, `name`, enname, `desc`, endesc, invites FROM tmm.user_levels ORDER BY id ASC")
	if CheckErr(err, c) {
		return
	}
	var levels []common.CreditLevel
	for _, row := range rows {
		l := common.CreditLevel{
			Id:      row.Uint(0),
			Name:    row.Str(1),
			Enname:  row.Str(2),
			Desc:    row.Str(3),
			Endesc:  row.Str(4),
			Invites: row.Uint(5),
		}
		levels = append(levels, l)
	}
	c.JSON(http.StatusOK, levels)
}
