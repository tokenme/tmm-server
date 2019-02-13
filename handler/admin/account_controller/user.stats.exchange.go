package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
)

const (
	Point = 1
	Uc    = -1
)

func ExchangeHandler(c *gin.Context) {
	db := Service.Db
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}

	if Check(req.Id < 0 || req.Types != Uc && req.Types != Point, admin.Error_Param, c) {
		return
	}

	var offset int
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	}

	query := `
SELECT
	DATE_ADD(inserted_at,INTERVAL 8 HOUR),
	points,
	tmm,
	status
FROM tmm.exchange_records
WHERE user_id=%d AND %s
ORDER BY inserted_at DESC
LIMIT %d OFFSET %d`

	var where []string
	if req.Types != 0 {
		where = append(where, fmt.Sprintf("direction=%d ", req.Types))
	}

	if req.Devices != "" {
		where = append(where, fmt.Sprintf("device_id='%s'", db.Escape(req.Devices)))
	}
	rows, _, err := db.Query(query, req.Id, strings.Join(where, " AND "), req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var exchangeList []*Task

	for _, row := range rows {

		exchange := &Task{
			When:   row.Str(0),
			Status: MsgMap[row.Int(3)],
		}

		switch req.Types {
		case Uc:
			exchange.Type = ExchangeUc
			exchange.Pay = fmt.Sprintf("-%.0fUC", row.Float(2))
			exchange.Get = fmt.Sprintf("+%.0f积分", row.Float(1))
		case Point:
			exchange.Type = ExchangePoint
			exchange.Pay = fmt.Sprintf("-%.0f积分", row.Float(1))
			exchange.Get = fmt.Sprintf("+%.0fUC", row.Float(2))
		}

		exchangeList = append(exchangeList, exchange)

	}

	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.exchange_records WHERE user_id=%d AND %s`, req.Id, strings.Join(where, " AND "))
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
