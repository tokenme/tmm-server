package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
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
	DATE_ADD(inserted_at,INTERVAL 8 HOUR),
	verified,
	IF(verified=-1,0,IF(verified = 1,IF(trade_num != "",1,2),3)) 
FROM tmm.point_withdraws
WHERE user_id=%d
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
		drawMoney.Pay = fmt.Sprintf("-%.0f 积分", row.Float(0))
		drawMoney.Get = fmt.Sprintf("+%.2f CNY", row.Float(1))
		drawMoney.When = row.Str(2)
		drawMoney.Extra = AuditMsgMap[row.Int(3)]
		drawMoney.Status = MsgMap[row.Int(4)]
		DrawMoneyList = append(DrawMoneyList, drawMoney)
	}
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.point_withdraws WHERE user_id=%d`, req.Id)
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
