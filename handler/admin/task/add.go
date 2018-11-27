package task

import (
	"github.com/gin-gonic/gin"
	"time"
	"github.com/shopspring/decimal"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"net/http"
)

func AddShareHandler(c *gin.Context) {
	var req ShareAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var (
		db        = Service.Db
		cidInsert = `Insert INTO tmm.share_task_categories (task_id,cid) VALUE(%d,%d)`
	)

	_, ret, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, points, points_left, bonus, max_viewers) VALUES (%d, '%s', '%s', '%s', '%s', %s, %s, %s, %d)`, 0, db.Escape(req.Title), db.Escape(req.Summary), db.Escape(req.Link), db.Escape(req.Image), db.Escape(req.Points), db.Escape(req.Points), db.Escape(req.Bonus), req.MaxViewers)
	if CheckErr(err, c) {
		return
	}
	for _, cid := range req.Cid {
		_, _, err := db.Query(cidInsert, ret.InsertId(), cid)
		if CheckErr(err, c) {
			return
		}
	}
	now := time.Now().Format(time.RFC3339)
	boints, err := decimal.NewFromString(req.Points)
	if CheckErr(err, c) {
		return
	}
	bonus, err := decimal.NewFromString(req.Bonus)
	if CheckErr(err, c) {
		return
	}
	task_ := common.ShareTask{
		Id:         ret.InsertId(),
		Title:      req.Title,
		Summary:    req.Summary,
		Link:       req.Link,
		Image:      req.Image,
		Points:     boints,
		PointsLeft: boints,
		Bonus:      bonus,
		MaxViewers: req.MaxViewers,
		InsertedAt: now,
		UpdatedAt:  now,
		Creator:    0,
		Cid:        req.Cid,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": task_})
}
