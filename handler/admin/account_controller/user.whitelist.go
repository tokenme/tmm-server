package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

func WhiteListHandler(c *gin.Context) {
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var offset int
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Page > 1 {
		offset = (req.Page - 1) * req.Limit
	}

	query := `
SELECT 
	us.user_id,
	IFNULL(wx.nick,u.nickname),
	u.mobile,
	us.comments
FROM 
	tmm.user_settings  AS us
INNER JOIN ucoin.users AS u ON (u.id = us.user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = us.user_id)

WHERE block_whitelist = 1 LIMIT %d OFFSET %d`

	db := Service.Db
	rows, _, err := db.Query(query, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var list []*admin.UserStats
	for _, row := range rows {
		user := &admin.UserStats{}
		user.Id = row.Uint64(0)
		user.Nick = row.Str(1)
		user.Mobile = row.Str(2)
		user.BlockedMessage = row.Str(3)
		list = append(list, user)
	}

	total := 0
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.user_settings WHERE  block_whitelist = 1`)
	if CheckErr(err, c) {
		return
	}
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"total": total,
			"data":  list,
		},
	})

}
