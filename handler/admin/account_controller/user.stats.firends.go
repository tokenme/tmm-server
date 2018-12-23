package account_controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const (
	Direct   = `inv.parent_id`
	Indirect = `inv.grand_id`
)

func FriendsHandler(c *gin.Context) {
	db := Service.Db
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	query := `
SELECT
	inv.user_id,
	u.mobile,
	wx.nick
FROM 
	 tmm.invite_codes AS inv
INNER JOIN ucoin.users AS u ON u.id = inv.user_id 
LEFT JOIN tmm.wx AS wx ON wx.user_id = inv.user_id 
WHERE %s = %d`
	var directList []*admin.Users
	var indirectList []*admin.Users
	if id > 0 {
		rows, _, err := db.Query(query, db.Escape(Direct), id)
		if CheckErr(err, c) {
			return
		}
		for _, row := range rows {
			user := &admin.Users{}
			user.Id = row.Uint64(0)
			user.Mobile = row.Str(1)
			user.Nick = row.Str(2)
			directList = append(directList, user)
		}

		rows, _, err = db.Query(query, db.Escape(Indirect), id)
		for _, row := range rows {
			user := &admin.Users{}
			user.Id = row.Uint64(0)
			user.Mobile = row.Str(1)
			user.Nick = row.Str(2)
			indirectList = append(indirectList, user)
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"direct":   directList,
			"indirect": indirectList,
		},
	})
}
