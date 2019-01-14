package account_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
)

func MakePointHandler(c *gin.Context) {
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
		0  AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_type = 0`, req.Id))
	}
	if req.Types == Reading || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`	
	SELECT 
		point AS point,
		inserted_at AS inserted_at,
		0 AS type,
		0 AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm
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
		sha.device_id AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm
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
		0  AS device_id,
		task_type  AS invite_bonus_types,
		tmm  AS tmm
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_type != 0 `, req.Id))
	}
	if req.Types == AppTask || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT 
		app.points AS point,
		app.inserted_at AS inserted_at,
		4 AS type, 
		app.device_id AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm 
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
		tmp.type  AS type,
		tmp.invite_bonus_types AS invite_bonus_types,
		tmp.tmm  AS tmm
	FROM(
		%s
	) AS tmp
	 %s
	ORDER BY 
		tmp.inserted_at DESC
	LIMIT %d OFFSET %d
	`

	var where string
	db := Service.Db
	if req.Devices != "" && (req.Types == AppTask || req.Types == Share) {
		where = fmt.Sprintf(" WHERE tmp.device_id ='%s'", db.Escape(req.Devices))
	}
	rows, _, err := db.Query(query, strings.Join(froms, " UNION "), where, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var taskList []*Task
	var taskType string
	var get string
	for _, row := range rows {
		get = fmt.Sprintf("+%.2f积分", row.Float(0))
		if row.Int(3) > 0 {
			taskType = InviteMap[row.Int(3)]
			if row.Int(3) == 3 {
				get = fmt.Sprintf("+%.2fUC", row.Float(4))
			}
		} else {
			taskType = typeMap[row.Int(2)]
		}
		task := &Task{
			Get:    get,
			When:   row.Str(1),
			Type:   taskType,
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
