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
	id,
	creator,
	title,
	summary,
	image,
	details,
	points,
	points_left,
	bonus,
	online_status,
	inserted_at
FROM 
	tmm.general_tasks
WHERE 
    1 = 1 %s
LIMIT %d OFFSET %d 

`

	db := Service.Db
	rows, _, err := db.Query(query, strings.Join(where, ` `), req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var list []*common.GeneralTask
	for _, row := range rows {
		list = append(list, &common.GeneralTask{
			Id:           row.Uint64(0),
			Creator:      row.Uint64(1),
			Title:        row.Str(2),
			Summary:      row.Str(3),
			Image:        row.Str(4),
			Details:      row.Str(5),
			Points:       decimal.NewFromFloat(row.Float(6)),
			PointsLeft:   decimal.NewFromFloat(row.Float(7)),
			Bonus:        decimal.NewFromFloat(row.Float(8)),
			OnlineStatus: int8(row.Int(9)),
			InsertedAt:   row.Str(10),
		})
	}

	rows, _, err = db.Query(`SELECT COUNT(1) FROm tmm.general_tasks  WHERE 1 = 1 %s`, strings.Join(where, ` AND `))
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
			`total`: total,
			`data`:  list,
		},
	})
}
