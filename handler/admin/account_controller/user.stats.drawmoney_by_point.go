package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func DrawMoneyByPointHandler(c *gin.Context) {
	db := Service.Db
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var offset int
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	} else {
		offset = 0
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
	rows, _, err := db.Query(query, req.Id, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}
	var DrawMoneyList []*Task
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
		drawMoney := &Task{}
		drawMoney.Type = DrawMoneyByPoint
		drawMoney.Pay = fmt.Sprintf("-%.2f 积分", row.Float(0))
		drawMoney.Get = fmt.Sprintf("+%.2f CNY", row.Float(1))
		drawMoney.When = row.Str(2)
		drawMoney.Status = MsgMap[Success]
		DrawMoneyList = append(DrawMoneyList, drawMoney)
	}
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.point_withdraws WHERE user_id = %d ORDER BY trade_num  `, req.Id)
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
			"page":  req.Page,
		},
	})
}
