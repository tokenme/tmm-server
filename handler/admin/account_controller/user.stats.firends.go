package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
)

type types int

const (
	Direct types = iota
	Indirect
	Children
	Active
)

func FriendsHandler(c *gin.Context) {
	db := Service.Db
	var req PageOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.Id < 0, admin.Not_Found, c) {
		return
	}
	var offset int
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	}
	var totalquery, query string
	query = `
	SELECT
		inv.user_id,
		u.mobile,
		wx.nick
	FROM 
		tmm.invite_codes AS inv
	INNER JOIN 
		ucoin.users AS u ON u.id = inv.user_id 
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = inv.user_id 
	WHERE %s
	GROUP BY inv.user_id
	LIMIT %d OFFSET %d`
	totalquery = `
	SELECT
		COUNT(distinct inv.user_id)
	FROM 
		tmm.invite_codes AS inv
	INNER JOIN 
		ucoin.users AS u ON u.id = inv.user_id 
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = inv.user_id 
	WHERE 
		 %s
`
	switch types(req.Types) {
	case Direct:
		direct := fmt.Sprintf(" inv.parent_id = %d ", req.Id)
		query = fmt.Sprintf(query, direct, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, direct)
	case Indirect:
		indirect := fmt.Sprintf(" inv.grand_id = %d", req.Id)
		query = fmt.Sprintf(query, indirect, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, indirect)
	case Children:
		online := fmt.Sprintf("  inv.parent_id = %d OR inv.grand_id = %d ", req.Id, req.Id)
		query = fmt.Sprintf(query, online, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, online)

	case Active:
		active := fmt.Sprintf(` (inv.parent_id = %d OR inv.grand_id = %d)  AND EXISTS(
		SELECT 
		1
		FROM tmm.devices AS dev 
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at > DATE(NOW()))
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id  AND  app.inserted_at > DATE(NOW()) AND app.status = 1)
		LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id  AND (reading.updated_at > DATE(NOW()) OR  reading.inserted_at > DATE(NOW())))
		WHERE dev.user_id = inv.user_id AND ( 
		IFNULL(sha.points,0) > 0  
		OR IFNULL(app.points,0) > 0   OR IFNULL(reading.point,0) > 0)
		LIMIT 1 )`, req.Id, req.Id)
		query = fmt.Sprintf(query, active, req.Limit, offset)
		totalquery = fmt.Sprintf(totalquery, active)
	default:
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data: gin.H{
				"data":  nil,
				"total": nil,
			},
		})
		return
	}

	var List []*admin.Users
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data: gin.H{
				"data":  List,
				"total": 0,
			},
		})
		return
	}
	for _, row := range rows {
		user := &admin.Users{}
		user.Id = row.Uint64(0)
		user.Mobile = row.Str(1)
		user.Nick = row.Str(2)
		List = append(List, user)
	}

	var total int
	rows, _, err = db.Query(totalquery)
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
			"data":  List,
			"total": total,
		},
	})
}
