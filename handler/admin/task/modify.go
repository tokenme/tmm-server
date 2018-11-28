package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func ModifyTaskHandler(c *gin.Context) {
	var (
		task      ShareTask
		query     string
		cidQuery  = `select cid from tmm.share_task_categories where task_id = %d`
		cidInsert = `Insert INTO tmm.share_task_categories (task_id,cid) VALUE(%d,%d)`
		db        = Service.Db
		cliList   = make(map[int]struct{})
	)
	if CheckErr(c.Bind(&task), c) {
		return
	}
	query = `update tmm.share_tasks set title='%s',summary='%s',link='%s',image='%s',points='%s',bonus='%s',max_viewers=%d where id = %d `
	_, _, err := db.Query(query, db.Escape(task.Title), db.Escape(task.Summary), db.Escape(task.Link),
		db.Escape(task.Image), db.Escape(task.Points), db.Escape(task.Bonus),
		task.MaxViewers, task.Id)
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
	c.JSON(http.StatusOK, gin.H{
		"code":http.StatusOK,
		"msg":"",
		"data":"",
	})
}
