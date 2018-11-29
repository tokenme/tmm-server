package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"fmt"
	"github.com/tokenme/tmm/handler/task"
	"strings"
	"github.com/tokenme/tmm/handler/admin"
)

func AddShareHandler(c *gin.Context) {
	var db = Service.Db
	var req task.ShareAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var query = `
	INSERT INTO tmm.share_tasks 
	(creator, title, summary, link,image, points, points_left, bonus, max_viewers) 
	VALUES 	(%d, '%s', '%s', '%s', '%s', %s, %s, %s, %d)
	`
	_, res, err := db.Query(query, 0, db.Escape(req.Title),
		db.Escape(req.Summary), db.Escape(req.Link), db.Escape(req.Image), db.Escape(req.Points.String()),
		db.Escape(req.Points.String()), db.Escape(req.Bonus.String()), req.MaxViewers)
	if CheckErr(err, c) {
		return
	}
	var insertArray []string
	for _, cid := range req.Cid {
		insertArray = append(insertArray, fmt.Sprintf(`(%d,%d,%d)`, res.InsertId(), cid, 1))
	}
	if len(insertArray) > 0 {
		cidInsert := `INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUE %s  ON 
		DUPLICATE KEY UPDATE cid=VALUES(cid),is_auto=VALUES(is_auto)`
		if _, _, err := db.Query(cidInsert, strings.Join(insertArray, `,`)); CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Data:    req,
		Code:    0,
	},
	)
}
