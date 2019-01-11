package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"strconv"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func GetTaskHandler(c *gin.Context) {
	taskid, err := strconv.Atoi(c.Query(`taskid`))
	if CheckErr(err, c) {
		return
	}

	var query = `SELECT 
	creator,
	title,
	summary,
	link,
	image,
	points,
	points_left,
	bonus,
	max_viewers,
	viewers,
	online_status,
	inserted_at,
	updated_at 
	FROM tmm.share_tasks WHERE id = %d`

	var db = Service.Db
	rows, res, err := db.Query(query, taskid)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}

	row := rows[0]
	points, err := decimal.NewFromString(row.Str(res.Map(`points`)))
	if CheckErr(err, c) {
		return
	}
	pointsLeft, err := decimal.NewFromString(row.Str(res.Map(`points_left`)))
	if CheckErr(err, c) {
		return
	}
	bonus, err := decimal.NewFromString(row.Str(res.Map(`bonus`)))
	if CheckErr(err, c) {
		return
	}
	var task common.ShareTask
	task.Id = uint64(taskid)
	task.Creator = row.Uint64(res.Map(`creator`))
	task.Title = row.Str(res.Map(`title`))
	task.Summary = row.Str(res.Map(`summary`))
	task.Link = row.Str(res.Map(`link`))
	task.Image = row.Str(res.Map(`image`))
	task.MaxViewers = row.Uint(res.Map(`max_viewers`))
	task.OnlineStatus = int8(row.Int(res.Map(`online_status`)))
	task.InsertedAt = row.Str(res.Map(`inserted_at`))
	task.UpdatedAt = row.Str(res.Map(`updated_at`))
	task.Points = points
	task.PointsLeft = pointsLeft
	task.Bonus = bonus

	var cidQuery = `SELECT cid FROM tmm.share_task_categories WHERE task_id = %d`
	rows, _, err = db.Query(cidQuery, taskid)
	if CheckErr(err, c) {
		return
	}

	var cidList = []int{}
	for _, row := range rows {
		cidList = append(cidList, row.Int(0))
	}
	task.Cid = cidList

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    task,
	})
}
