package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
)


type SearchOptions struct {
	Id     int    `form:"id"`
	Mobile string `form:"mobile"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

func GetAccountList(c *gin.Context) {
	var search SearchOptions
	if CheckErr(c.Bind(&search), c) {
		return
	}

	var offset int
	if search.Limit < 1 {
		search.Limit = 10
	}
	if search.Page > 0 {
		offset = (search.Page - 1) * search.Limit
	} else {
		offset = 0
	}

	var where []string

	if search.Id > 0 {
		where = append(where, fmt.Sprintf(" AND u.id = %d ", search.Id))
	}
	if search.Mobile != "" {
		where = append(where, fmt.Sprintf(" AND u.mobile = '%s' ", search.Mobile))
	}

	db := Service.Db
	query := `
SELECT 
	u.id,
	IFNULL(wx.nick,u.nickname),
	u.mobile,
	IF(us.user_id > 0 , IF(us.block_whitelist = us.blocked,0,1) ,0)
FROM 
	ucoin.users AS u 
LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
LEFT JOIN tmm.user_settings AS us ON (us.user_id = u.id)
WHERE 
	1 = 1 %s
ORDER BY u.created DESC
LIMIT %d OFFSET %d`

	rows, _, err := db.Query(query, strings.Join(where, " "), search.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var list []*admin.User
	for _, row := range rows {
		user := &admin.User{}
		user.Id = row.Uint64(0)
		user.Nick = row.Str(1)
		user.Mobile = row.Str(2)
		user.Blocked = row.Int(3)
		list = append(list, user)
	}

	var total int

	rows, _, err = db.Query(` SELECT COUNT(*) FROM ucoin.users AS u LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id) WHERE 1 = 1 %s`, strings.Join(where, " "))
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
			"data":  list,
			"total": total,
		},
	})
}
