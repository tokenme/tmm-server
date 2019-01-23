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

	var (
		froms      []string
		countForms []string
	)
	db := Service.Db
	if req.Types == Invite || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT
		bonus AS point,
		DATE_ADD(inserted_at,INTERVAL 8 HOUR) AS inserted_at,
		1 AS type ,
		0  AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm,
		0  AS extra
	FROM tmm.invite_bonus
	WHERE user_id=%d AND task_type=0 AND user_id!=from_user_id`, req.Id))
	}
	countForms = append(countForms, fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.invite_bonus WHERE user_id=%d AND task_type=0 AND user_id!=from_user_id`, req.Id))
	if req.Types == Reading || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(`
	SELECT
		point AS point,
		DATE_ADD(inserted_at,INTERVAL 8 HOUR) AS inserted_at,
		0 AS type,
		0 AS device_id,
		0  AS invite_bonus_types,
		0  AS tmm,
		ts AS extra
	FROM tmm.reading_logs
 	WHERE user_id=%d`, req.Id))
		countForms = append(countForms, fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.reading_logs WHERE user_id=%d`, req.Id))
	}
	if req.Types == Share || req.Types == -1 {
		var (
			q      string
			countQ string
		)
		if req.Devices != "" {
			q = fmt.Sprintf(`
    SELECT
        sha.points AS point,
        DATE_ADD(sha.inserted_at,INTERVAL 8 HOUR) AS inserted_at,
        2 AS type,
        sha.device_id AS device_id,
        0  AS invite_bonus_types,
        0  AS tmm,
        sha.viewers AS extra
    FROM
        tmm.device_share_tasks AS sha
    INNER JOIN tmm.devices AS dev ON (dev.id=sha.device_id)
    WHERE dev.user_id=%d AND dev.device_id='%s'`, req.Id, db.Escape(req.Devices))
			countQ = fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.device_share_tasks AS dst INNER JOIN tmm.devices AS d ON (d.id=dst.device_id) WHERE d.user_id=%d AND d.id='%s'`, req.Id, db.Escape(req.Devices))
		} else {
			q = fmt.Sprintf(`
    SELECT
        sha.points AS point,
        DATE_ADD(sha.inserted_at,INTERVAL 8 HOUR) AS inserted_at,
        2 AS type,
        sha.device_id AS device_id,
        0  AS invite_bonus_types,
        0  AS tmm,
        sha.viewers AS extra
    FROM
        tmm.device_share_tasks AS sha
    INNER JOIN tmm.devices AS dev ON (dev.id=sha.device_id)
    WHERE dev.user_id=%d`, req.Id)
			countQ = fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.device_share_tasks AS dst INNER JOIN tmm.devices AS d ON (d.id=dst.device_id) WHERE d.user_id=%d`, req.Id)
		}
		froms = append(froms, q)
		countForms = append(countForms, countQ)
	}
	if req.Types == BfBouns || req.Types == -1 {
		froms = append(froms, fmt.Sprintf(
			`
	SELECT
		bonus AS point,
		DATE_ADD(inserted_at,INTERVAL 8 HOUR) AS inserted_at,
		3 AS type ,
		0  AS device_id,
		task_type  AS invite_bonus_types,
		tmm  AS tmm,
		0 AS extra
	FROM tmm.invite_bonus
	WHERE user_id=%d AND task_type!=0`, req.Id))
		countForms = append(countForms, fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.invite_bonus WHERE user_id=%d AND task_type!=0`, req.Id))
	}
	if req.Types == AppTask || req.Types == -1 {
		var (
			countQ string
			q      string
		)
		if req.Devices != "" {
			q = fmt.Sprintf(`
    SELECT
        app.points AS point,
        DATE_ADD(app.inserted_at,INTERVAL 8 HOUR) AS inserted_at,
        4 AS type,
        app.device_id AS device_id,
        0  AS invite_bonus_types,
        0  AS tmm,
        0 AS extra
    FROM tmm.device_app_tasks AS app
    INNER JOIN tmm.devices AS dev ON (dev.id=app.device_id)
    WHERE dev.user_id=%d AND app.status=1 AND d.id='%s'`, req.Id, db.Escape(req.Devices))
			countQ = fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.device_app_tasks AS dat INNER JOIN tmm.devices AS d ON (d.id=dat.device_id) WHERE d.user_id=%d AND dat.status=1 AND d.id='%s'`, req.Id, db.Escape(req.Devices))
		} else {
			q = fmt.Sprintf(`
    SELECT
        app.points AS point,
        DATE_ADD(app.inserted_at,INTERVAL 8 HOUR) AS inserted_at,
        4 AS type,
        app.device_id AS device_id,
        0  AS invite_bonus_types,
        0  AS tmm,
        0 AS extra
    FROM tmm.device_app_tasks AS app
    INNER JOIN tmm.devices AS dev ON (dev.id=app.device_id)
    WHERE dev.user_id=%d AND app.status=1`, req.Id)
			countQ = fmt.Sprintf(`SELECT COUNT(1) AS num FROM tmm.device_app_tasks AS dat INNER JOIN tmm.devices AS d ON (d.id=dat.device_id) WHERE d.user_id=%d AND dat.status=1`, req.Id)
		}
		froms = append(froms, q)
		countForms = append(countForms, countQ)
	}
	query := `%s ORDER BY inserted_at DESC LIMIT %d OFFSET %d`
	rows, _, err := db.Query(query, strings.Join(froms, " UNION "), req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var taskList []*Task
	var taskType string
	for _, row := range rows {
		extra := ""
		get := fmt.Sprintf("+%.2f积分", row.Float(0))
		if row.Int(4) > 0 {
			taskType = InviteMap[row.Int(4)]
			if row.Int(4) == 3 {
				get = fmt.Sprintf("+%.2fUC", row.Float(5))
			}
		} else {
			taskType = typeMap[row.Int(2)]
		}

		if value, ok := ForMatMap[row.Int(2)]; ok {
			extra = fmt.Sprintf(value, row.Int(6))
		}
		task := &Task{
			Get:    get,
			When:   row.Str(1),
			Type:   taskType,
			Status: TaskSuccessful,
			Extra:  extra,
		}
		taskList = append(taskList, task)
	}

	var total int
	countQuery := `SELECT SUM(num) FROM (%s) AS tmp`
	rows, _, err = db.Query(countQuery, strings.Join(countForms, " UNION "))
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
