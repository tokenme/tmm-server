package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"strconv"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func DrawMoneyByPointHandler(c *gin.Context) {
	db := Service.Db
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, `-1`))
	if CheckErr(err, c) {
		return
	}
	limit:=10
	var offset int
	if page > 0 {
		offset = (page-1)*limit
	}
	query := `
SELECT 
	points,
	cny,
	inserted_at
FROM 
	tmm.point_withdraws 
WHERE 
	user_id = %d
ORDER BY inserted_at DESC
LIMIT %d OFFSET %d`
	rows, _, err := db.Query(query, id, limit, offset)
	if CheckErr(err, c) {
		return
	}
	var DrawMoneyList []*DrawMoney
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
			Data: gin.H{
				"data":  DrawMoneyList,
				"total": 0,
			},
		})
	}
	for _, row := range rows {
		drawMoney := &DrawMoney{}
		drawMoney.Type = DrawMoneyByPoint
		drawMoney.Pay = row.Str(0)
		drawMoney.Get = fmt.Sprintf("%.2f", row.Float(1))
		drawMoney.When = row.Str(2)
		drawMoney.Status = MsgMap[Success]
		DrawMoneyList = append(DrawMoneyList, drawMoney)
	}
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.point_withdraws WHERE user_id = %d ORDER BY trade_num  `, id)
	if CheckErr(err, c) {
		return
	}
	var total = 0
	if len(rows) != 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"data":  DrawMoneyList,
			"total": total,
			"page":page,
		},
	})
}
