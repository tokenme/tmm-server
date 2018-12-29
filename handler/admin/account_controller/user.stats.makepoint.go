package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"strings"
)

func MakePointHandler(c *gin.Context) {
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
	var froms []string
	if req.Types == Invite || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT 
		bonus AS point,
		inserted_at AS inserted_at,	
		1 AS type ,
		0  AS device_id
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_id = 0`, req.Id))
	}
	if req.Types == Reading || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`	
	SELECT 
		point AS point,
		inserted_at AS inserted_at,
		0 AS type,
		0 AS device_id 
	FROM 
 		tmm.reading_logs
 	WHERE
		user_id = %d`, req.Id))
	}
	if req.Types == Share || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT 
		sha.points AS point,
		sha.inserted_at AS inserted_at,
		2 AS type,
		sha.device_id AS device_id 
	FROM 
		tmm.device_share_tasks AS sha
	INNER JOIN tmm.devices AS dev ON (dev.id = sha.device_id)
	WHERE 
		dev.user_id = %d  `, req.Id))
	}
	if req.Types == BfBouns || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(
			`
	SELECT 
		bonus AS point,
		inserted_at AS inserted_at,	
		3 AS type ,
		0  AS device_id 
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_id != 0`, req.Id))
	}
	if req.Types == AppTask || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT 
		app.points AS point,
		app.inserted_at AS inserted_at,
		4 AS type, 
		app.device_id AS device_id
	FROM
		tmm.device_app_tasks AS app
	INNER JOIN tmm.devices AS dev ON (dev.id = app.device_id)
	WHERE 
		dev.user_id = %d  AND app.status = 1 `, req.Id))
	}
	query := `
	SELECT 
		IFNULL(tmp.point,0) AS point,
		tmp.inserted_at AS inserted_at,
		tmp.type  AS type
	FROM(
		%s
	) AS tmp
	 %s
	ORDER BY 
		tmp.inserted_at DESC
	LIMIT %d OFFSET %d
	`
	var where string
	if req.Devices != "" {
		where = fmt.Sprintf(" WHERE tmp.device_id IN (%d,%s)", 0, req.Devices)
	}
	rows, _, err := db.Query(query, strings.Join(froms, " UNION "), where, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}
	var taskList []*Task
	for _, row := range rows {
		task := &Task{
			Get:    fmt.Sprintf("+%.2f积分", row.Float(0)),
			When:   row.Str(1),
			Type:   typeMap[row.Int(2)],
			Status: TaskSuccessful,
		}
		taskList = append(taskList, task)
	}
	var total int
	rows, _, err = db.Query(`
	SELECT 
		COUNT(1)
	FROM(
		%s
	) AS tmp 
	%s `, strings.Join(froms, " UNION "), where)
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"data":  taskList,
			"total": total,
			"page":  req.Page,
		},
	})
}
