package general_task

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strconv"
	"fmt"
	"strings"
)

type GeneralTask struct {
	common.GeneralTask
	Completed uint64 `json:"completed"`
	Rejected  uint64 `json:"rejected"`
	Waiting   uint64 `json:"waiting"`
}

func GeneralTaskListHandler(c *gin.Context) {
	var req admin.Pages
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var offset int
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Page > 1 {
		offset = (req.Page - 1) * req.Limit
	}

	var where []string
	status, err := strconv.Atoi(c.DefaultQuery(`online_status`, `0`))
	if CheckErr(err, c) {
		return
	}

	id, err := strconv.Atoi(c.DefaultQuery(`id`, `0`))
	if CheckErr(err, c) {
		return
	}

	if status != 0 {
		where = append(where, fmt.Sprintf(`AND online_status = %d`, status))
	}
	if id > 0 {
		where = append(where, fmt.Sprintf(`AND id = %d`, id))
	}

	query := `
SELECT 
	gt.id,
	gt.creator,
	gt.title,
	gt.summary,
	gt.image,
	gt.details,
	gt.points,
	gt.points_left,
	gt.bonus,
	gt.online_status,
	gt.inserted_at,
	gt.completed,
	COUNT(IF(dgt.status = 0,0,NULL)) AS waiting,
	COUNT(IF(dgt.status = -1,0,NULL)) AS rejected
FROM 
	tmm.general_tasks AS gt
LEFT JOIN tmm.device_general_tasks AS dgt ON (dgt.task_id = gt.id)
WHERE 
    1 = 1 %s
GROUP BY gt.id 
ORDER BY gt.id 
LIMIT %d OFFSET %d

`

	db := Service.Db
	rows, _, err := db.Query(query, strings.Join(where, ` `), req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var list []*GeneralTask
	for _, row := range rows {
		task := &GeneralTask{}
		task.Id = row.Uint64(0)
		task.Creator = row.Uint64(1)
		task.Title = row.Str(2)
		task.Summary = row.Str(3)
		task.Image = row.Str(4)
		task.Details = row.Str(5)
		task.Points = decimal.NewFromFloat(row.Float(6))
		task.PointsLeft = decimal.NewFromFloat(row.Float(7))
		task.Bonus = decimal.NewFromFloat(row.Float(8))
		task.OnlineStatus = int8(row.Int(9))
		task.InsertedAt = row.Str(10)
		task.Completed = row.Uint64(11)
		task.Waiting = row.Uint64(12)
		task.Rejected = row.Uint64(13)
		list = append(list, task)
	}

	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.general_tasks  WHERE 1 = 1 %s`, strings.Join(where, `  `))
	if CheckErr(err, c) {
		return
	}

	var total int
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"total": total,
			"data":  list,
		},
	})
}
