package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"strconv"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)


func ExchangeByPointHandler(c *gin.Context) {
	db := Service.Db
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, `-1`))
	var offset int
	limit := 10
	if page > 0 {
		offset = limit * (page - 1)
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
	if id < 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
		})
		return
	}
	rows, _, err := db.Query(query, id, limit, offset)
	if CheckErr(err, c) {
		return
	}
	var exchangeList []*Exchange
	for _, row := range rows {
		exchange := &Exchange{
			Type:   ExchangePoint,
			Pay:    fmt.Sprintf("-%.2fUC", row.Float(2)),
			Get:    fmt.Sprintf("+%.2f积分", row.Float(1)),
			When:   row.Str(0),
			Status: MsgMap[row.Int(3)],
		}
		exchangeList = append(exchangeList, exchange)
	}
	rows, _, err = db.Query(`SELECT COUNT(1)  FROM tmm.exchange_records WHERE user_id = %d AND direction = -1`, id)
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
			"page":  page,
		},
	})
}
