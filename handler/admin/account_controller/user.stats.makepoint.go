package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
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
	if req.Page> 0 {
		offset = (req.Page - 1) * req.Limit
	} else {
		offset = 0
	}
	query := `
SELECT 
	IFNULL(tmp.point,0) AS point,
	tmp.inserted_at AS inserted_at,
	tmp.type  AS type
FROM(
	SELECT 
		bonus AS point,
		inserted_at AS inserted_at,	
		1 AS type 
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_id = 0
	UNION 
	SELECT 
		bonus AS point,
		inserted_at AS inserted_at,	
		3 AS type 
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d AND task_id != 0
	UNION
	SELECT 
		point AS point,
		inserted_at AS inserted_at,
		0 AS type
	FROM 
 		tmm.reading_logs
 	WHERE
		user_id = %d
	UNION 
	SELECT 
		sha.points AS point,
		sha.inserted_at AS inserted_at,
		2 AS type
	FROM 
		tmm.device_share_tasks AS sha
	INNER JOIN tmm.devices AS dev ON (dev.id = sha.device_id)
	WHERE 
		dev.user_id = %d
) AS tmp
ORDER BY tmp.inserted_at DESC
LIMIT %d OFFSET %d
	`

	rows, _, err := db.Query(query, req.Id, req.Id, req.Id,req.Id, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}
	var taskList []*Task
	for _, row := range rows {
		task := &Task{
			Get:  fmt.Sprintf("+%.2f积分", row.Float(0)),
			When:   row.Str(1),
			Type:   typeMap[row.Int(2)],
			Status: TaskSuccessful,
		}
		taskList = append(taskList, task)
	}
	var total int
	rows, _, err = db.Query(`SELECT 
	SUM(tmp.total) AS total
FROM(
	SELECT 
		COUNT(1) AS total
	FROM 
  		tmm.invite_bonus
	WHERE 
		user_id = %d
	UNION 
	SELECT 
		COUNT(1) AS total
	FROM 
 		tmm.reading_logs
 	WHERE
		user_id = %d
	UNION 
	SELECT 
		COUNT(1) AS total
	FROM 
		tmm.device_share_tasks AS sha
	INNER JOIN tmm.devices AS dev ON (dev.id = sha.device_id)
	WHERE 
		dev.user_id = %d
) AS tmp
`,req.Id,req.Id,req.Id)
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
