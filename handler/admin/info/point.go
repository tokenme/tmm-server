package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"fmt"
	"math"
)




func PointInfoHandler(c *gin.Context) {

	db :=Service.Db
	var req InfoRequest
	if CheckErr(c.Bind(&req),c){
		return
	}
	query :=``

}

func PointHandler(c *gin.Context) {
	db := Service.Db
	query := `SELECT
    COUNT(*) AS users,
    l
FROM (
    SELECT
        d.user_id,
        FLOOR(LOG10(SUM(d.points))) AS l
    FROM tmm.devices AS d
    GROUP BY d.user_id
) AS tmp
GROUP BY l ORDER BY l
`

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName []string
	var valueList []int
	for _, row := range rows {
		valueList = append(valueList, row.Int(0))
		Name := fmt.Sprintf(`%d-%d`, int(math.Pow10(row.Int(1)-1)), int(math.Pow10(row.Int(1))))
		indexName = append(indexName, Name)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: Data{
			Title:     "积分",
			IndexName: indexName,
			Value:     valueList,
		},
	})
}

func Top10PointHandler(c *gin.Context) {
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
		point, err := decimal.NewFromString(row.Str(4))
		if CheckErr(err, c) {
			return
		}
		User := User{
			Id:          row.Int(0),
			CountryCode: row.Int(1),
			Mobile:      row.Str(2),
			Nick:        row.Str(3),
			Point:       point,
		}
		userList = append(userList, User)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    userList,
	})
}
