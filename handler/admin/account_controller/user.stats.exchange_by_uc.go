package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)

func ExchangeByUcHandler(c *gin.Context) {
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
	inserted_at,
	points,
	tmm,
	status
FROM 
	tmm.exchange_records 
WHERE 
	user_id = %d AND direction = -1
ORDER BY inserted_at DESC
LIMIT %d OFFSET %d`
	if req.Id < 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
		})
		return
	}
	rows, _, err := db.Query(query, req.Id, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}
	var exchangeList []*Task
	for _, row := range rows {
		exchange := &Task{
			Type:   ExchangeUc,
			Pay:    fmt.Sprintf("-%.0fUC", row.Float(2)),
			Get:    fmt.Sprintf("+%.0f积分", row.Float(1)),
			When:   row.Str(0),
			Status: MsgMap[row.Int(3)],
		}
		exchangeList = append(exchangeList, exchange)
	}
	rows, _, err = db.Query(`SELECT COUNT(1)  FROM tmm.exchange_records WHERE user_id = %d AND direction = 1`, req.Id)
	if CheckErr(err, c) {
		return
	}
	var total int
	if len(rows) != 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"total": total,
			"data":  exchangeList,
			"page":  req.Page,
		},
	})
}
