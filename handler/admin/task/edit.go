package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"net/http"
)

func ExitTaskHandler(c *gin.Context) {
	var (
		task      common.ShareTask
		query     string
		cidQuery  = `select cid from tmm.share_task_categories task_id = %d`
		cidInsert = `Insert INTO tmm.share_task_categories (task_id,cid) VALUE(%d,%d)`
		db        = Service.Db
		cliList   = make(map[int]struct{})
	)
	if CheckErr(c.Bind(&task), c) {
		return
	}
	query = `update tmm.share_task set creator = %d, title='%s',summary='%s',link='%s',image='%s'
	,points='%s',points_left='%s',bonus='%s',max_viewers=%d,viewers=%d,online_status=%d where id = %d `
	_, _, err := db.Query(query, task.Creator, db.Escape(task.Title), db.Escape(task.Link), db.Escape(task.Image),
		db.Escape(task.Points.String()), db.Escape(task.PointsLeft.String()), db.Escape(task.Bonus.String()),
		task.MaxViewers, task.Viewers, task.OnlineStatus,task.Id)
	if CheckErr(err, c) {
		return
	}
	rows, _, err := db.Query(cidQuery, task.Id)
	if CheckErr(err, c) {
		return
	}
	for _, row := range rows {
		cliList[row.Int(0)] = struct{}{}
	}
	for _, cli := range task.Cid {
		if _, ok := cliList[cli]; !ok {
			_, _, err := db.Query(cidInsert, task.Id, cli)
			if CheckErr(err, c) {
				return
			}
		}
	}
	c.JSON(http.StatusOK,APIResponse{Msg:""})
}
