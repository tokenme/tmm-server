package general_task

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strconv"
)

func GetGeneralTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `1`))
	if CheckErr(err, c) {
		return
	}

	query := `SELECT
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
	id = %d
`

	db := Service.Db
	rows, _, err := db.Query(query, id)
	if CheckErr(err, c) {
		return
	}

	var task common.GeneralTask
	if len(rows) > 0 {
		row := rows[0]
		task.Id = uint64(id)
		task.Title = row.Str(0)
		task.Summary = row.Str(1)
		task.Image = row.Str(2)
		task.Details = row.Str(3)
		task.Points = decimal.NewFromFloat(row.Float(4))
		task.PointsLeft = decimal.NewFromFloat(row.Float(5))
		task.Bonus = decimal.NewFromFloat(row.Float(6))
		task.OnlineStatus = int8(row.Int(7))
		task.InsertedAt = row.Str(8)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    task,
	})
}
