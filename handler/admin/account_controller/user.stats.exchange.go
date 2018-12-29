package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
)

const (
	_Point = 1
	Uc     = -1
)

func ExchangeHandler(c *gin.Context) {
	db := Service.Db
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Id < 0 || req.Types != Uc && req.Types != _Point {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
		})
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
	user_id = %d AND %s
ORDER BY inserted_at DESC
LIMIT %d OFFSET %d`
	var where []string
	if req.Types == Uc {
		where = append(where, fmt.Sprintf(" direction = %d ", Uc))
	} else {
		where = append(where, fmt.Sprintf(" direction = %d ", Point))
	}
	if req.Devices != "" {
		where = append(where, fmt.Sprintf(" device_id = %s", req.Devices))
	}
	rows, _, err := db.Query(query, req.Id, strings.Join(where, " AND "), req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var exchangeList []*Task
	if req.Types == Uc {
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
	} else {
		for _, row := range rows {
			exchange := &Task{
				Type:   ExchangeUc,
				Pay:    fmt.Sprintf("-%.0f积分", row.Float(1)),
				Get:    fmt.Sprintf("+%.0fUC", row.Float(2)),
				When:   row.Str(0),
				Status: MsgMap[row.Int(3)],
			}
			exchangeList = append(exchangeList, exchange)
		}
	}
	rows, _, err = db.Query(`SELECT COUNT(1)  FROM tmm.exchange_records WHERE user_id = %d AND %s`, req.Id, strings.Join(where, " AND "))
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
