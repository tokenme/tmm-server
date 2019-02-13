package withdraw

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"github.com/tokenme/tmm/handler/admin/account_controller"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const (
	Point = 0
	UC    = 1
)

type SearchOption struct {
	Page   int `form:"page"`
	Limit  int `form:"limit"`
	Status int `form:"status"`
}

type WithDraw struct {
	Id             string `json:"id"`
	UserId         int    `json:"user_id"`
	Nick           string `json:"nick"`
	Mobile         string `json:"mobile"`
	InsertedAt     string `json:"inserted_at"`
	DrawCashType   string `json:"draw_cash_type"`
	Cny            string `json:"cny"`
	WithDrawStatus string `json:"with_draw_status"`
	Status         string `json:"status"`
}

func GetWithDrawList(c *gin.Context) {
	db := Service.Db

	var search SearchOption
	if CheckErr(c.Bind(&search), c) {
		return
	}

	if search.Limit < 1 {
		search.Limit = 20
	}

	var offset int
	if search.Page > 1 {
		offset = (search.Page - 1) * search.Limit
	}
	query := `
SELECT 
	tmp._index,
	tmp.user_id,
	us.mobile,
	IFNULL(wx.nick,us.nickname),
	tmp.inserted_at,
	tmp.types,
	tmp.cny,
	tmp.verified,
	tmp._status 
FROM (
	SELECT 
		id AS _index,
		user_id AS user_id ,
		inserted_at AS inserted_at,
		0 AS types,
		verified AS verified,
		cny AS cny,
		IF(verified=-1,0,IF(verified = 1,IF(trade_num != "",1,2),3))  AS _status
	FROM 
		tmm.point_withdraws  
	WHERE 
		verified = %d  UNION ALL 
	SELECT 
		tx AS _index,
		user_id AS user_id ,
		inserted_at AS inserted_at,
		1 AS types, 
		verified AS verified,
		cny AS cny,
		IF(verified=-1,0,IF(verified = 1,IF(trade_num != "",1,2),3))  AS _status 
	FROM 
		tmm.withdraw_txs
	WHERE 
		verified = %d
) AS tmp 
INNER JOIN ucoin.users AS us ON (us.id = tmp.user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = tmp.user_id)
ORDER BY tmp.inserted_at DESC 
LIMIT %d OFFSET %d
	`

	totalQuery := `
SELECT 
	COUNT(1)
FROM (
	SELECT 
		id AS _index,
		user_id AS user_id ,
		inserted_at AS inserted_at,
		0 AS types,
		verified AS verified,
		cny AS cny
	FROM 
		tmm.point_withdraws  
	WHERE 
		verified = %d  UNION ALL 
	SELECT 
		tx AS _index,
		user_id AS user_id ,
		inserted_at AS inserted_at,
		1 AS types, 
		verified AS verified,
		cny AS cny
	FROM 
		tmm.withdraw_txs
	WHERE 
		verified = %d
) AS tmp 
INNER JOIN ucoin.users AS us ON (us.id = tmp.user_id)
LEFT JOIN tmm.wx AS wx ON (wx.user_id = tmp.user_id)
ORDER BY tmp.inserted_at DESC
`
	rows, _, err := db.Query(query, search.Status, search.Status, search.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var withDrawList []*WithDraw

	var types string
	for _, row := range rows {
		if row.Int(5) == Point {
			types = "积分提现"
		} else if row.Int(5) == UC {
			types = "UC提现"
		}

		withDrawList = append(withDrawList, &WithDraw{
			Id:             row.Str(0),
			UserId:         row.Int(1),
			Mobile:         row.Str(2),
			Nick:           row.Str(3),
			InsertedAt:     row.Str(4),
			DrawCashType:   types,
			Cny:            fmt.Sprintf("%.2f", row.Float(6)),
			Status:         account_controller.AuditMsgMap[row.Int(7)],
			WithDrawStatus: account_controller.MsgMap[row.Int(8)],
		})
	}

	var total int

	rows, _, err = db.Query(totalQuery, search.Status, search.Status)
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
			"total": total,
			"data":  withDrawList,
		},
	})
}
