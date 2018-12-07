package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func UcInfoHandler(c *gin.Context) {
	db := Service.Db

	query := ``
	}

func Top10Handler(c *gin.Context) {
	db := Service.Db

	query := `SELECT
    u.id AS id,
    u.country_code AS country_code,
    u.mobile AS mobile,
    u.nickname AS nick,
    wx.nick AS wx_nick,
    SUM(d.points) AS points
FROM tmm.devices AS d
INNER JOIN ucoin.users AS u ON (u.id = d.user_id)
LEFT JOIN tmm.wx AS wx ON ( wx.user_id = u.id )
GROUP BY u.id
ORDER BY points DESC LIMIT 10`

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var userList []User
	for _, row := range rows {
		point, err := decimal.NewFromString(row.Str(5))
		if CheckErr(err, c) {
			return
		}
		user := User{
			Id:          row.Int(0),
			CountryCode: row.Int(1),
			Mobile:      row.Str(2),
			Nick:        row.Str(3),
			WxNick:      row.Str(4),
			Point:       point,
		}
		userList = append(userList, user)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    userList,
	})
}
