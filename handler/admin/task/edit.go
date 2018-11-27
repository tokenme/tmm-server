package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"strconv"
	"github.com/shopspring/decimal"
	"net/http"
)

func EditTaskHandler(c *gin.Context) {
	var (
		taskid   int
		db       = Service.Db
		query    string
		cidQuery = `select cid from tmm.share_task_categories where task_id = %d`
		task     common.ShareTask
		cidList  []int
	)
	taskid, err := strconv.Atoi(c.Query(`taskid`))
	if CheckErr(err, c) {
		return
	}
	query = `select creator,title,summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,inserted_at,updated_at from tmm.share_tasks where id = %d`
	rows, result, err := db.Query(query, taskid)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `Not found `, c) {
		return
	}
	row := rows[0]
	points, err := decimal.NewFromString(row.Str(result.Map(`points`)))
	if CheckErr(err, c) {
		return
	}
	pointsLeft, err := decimal.NewFromString(row.Str(result.Map(`points_left`)))
	if CheckErr(err, c) {
		return
	}
	bonus, err := decimal.NewFromString(row.Str(result.Map(`bonus`)))
	if CheckErr(err, c) {
		return
	}
	task.Id = uint64(taskid)
	task.Creator = row.Uint64(result.Map(`creator`))
	task.Title = row.Str(result.Map(`title`))
	task.Summary = row.Str(result.Map(`summary`))
	task.Link = row.Str(result.Map(`link`))
	task.Image = row.Str(result.Map(`image`))
	task.MaxViewers = row.Uint(result.Map(`max_viewers`))
	task.OnlineStatus = int8(row.Int(result.Map(`online_status`)))
	task.InsertedAt = row.Str(result.Map(`inserted_at`))
	task.UpdatedAt = row.Str(result.Map(`updated_at`))
	task.Points = points
	task.PointsLeft = pointsLeft
	task.Bonus = bonus
	rows, _, err = db.Query(cidQuery, taskid)
	if CheckErr(err, c) {
		return
	}
	for _, row := range rows {
		cidList = append(cidList, row.Int(0))

	}
	task.Cid = cidList
	c.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": task,
	})
}
