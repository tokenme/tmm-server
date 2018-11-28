package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"fmt"
	"github.com/tokenme/tmm/handler/task"
	"strings"
)

func AddShareHandler(c *gin.Context) {
	var req task.ShareAddRequest
	var insertArray []string
	if CheckErr(c.Bind(&req), c) {
		fmt.Println(req)
		return
	}
	var (
		db        = Service.Db
		cidInsert = `INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUE %s  ON DUPLICATE KEY UPDATE cid=VALUES(cid),is_auto=VALUES(is_auto)`
	)

	_, ret, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, points, points_left, bonus, max_viewers) VALUES (%d, '%s', '%s', '%s', '%s', %s, %s, %s, %d)`, 0, db.Escape(req.Title), db.Escape(req.Summary), db.Escape(req.Link), db.Escape(req.Image), db.Escape(req.Points.String()), db.Escape(req.Points.String()), db.Escape(req.Bonus.String()), req.MaxViewers)
	if CheckErr(err, c) {
		return
	}
	for _, cid := range req.Cid {
		insertArray = append(insertArray, fmt.Sprintf(`(%d,%d,%d)`, ret.InsertId(), cid, 1))
	}
	if len(insertArray) > 0 {
		_, _, err := db.Query(cidInsert, strings.Join(insertArray, `,`))
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    req})
}
