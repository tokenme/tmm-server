package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/qiniu"
	"github.com/mkideal/log"
	"net/http"
	"time"
)

type ShareAddRequest struct {
	Title         string          `json:"title" form:"title" binding:"required"`
	Summary       string          `json:"summary" form:"summary" binding:"required"`
	Link          string          `json:"link" form:"link" binding:"required"`
	Image         string          `json:"image" form:"image"`
	FileExtension string          `json:"image_extension" from:"image_extension"`
	Points        decimal.Decimal `json:"points" form:"points" binding:"required"`
	Bonus         decimal.Decimal `json:"bonus" form:"bonus" binding:"required"`
	MaxViewers    uint            `json:"max_viewers" form:"max_viewers" binding:"required"`
	Cid           []int           `json:"cid" form:"cid"`
}

func ShareAddHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req ShareAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	pointsPerTs, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := req.Points.Div(pointsPerTs)
	db := Service.Db
	query := `SELECT
d.id
FROM tmm.devices AS d
WHERE d.user_id = %d AND d.points >= %s
ORDER BY d.points DESC LIMIT 1`
	rows, _, err := db.Query(query, user.Id, req.Points)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	deviceId := rows[0].Str(0)
	_, ret, err := db.Query(`UPDATE tmm.devices AS d SET d.points=d.points-%s, d.consumed_ts = d.consumed_ts + %d WHERE id='%s' AND d.points>=%s AND d.user_id=%d`, req.Points.String(), ts.IntPart(), db.Escape(deviceId), req.Points.String(), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(ret.AffectedRows() == 0, NOT_ENOUGH_POINTS_ERROR, "not enough points in device", c) {
		return
	}
	if req.Image != "" && req.FileExtension == "webp" {
		newImage, _, err := qiniu.ConvertImage(req.Image, "jpeg", Config.Qiniu)
		if err != nil {
			log.Error(err.Error())
		} else {
			req.Image = newImage
		}
	}
	_, ret, err = db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, points, points_left, bonus, max_viewers) VALUES (%d, '%s', '%s', '%s', '%s', %s, %s, %s, %d)`, user.Id, db.Escape(req.Title), db.Escape(req.Summary), db.Escape(req.Link), db.Escape(req.Image), req.Points.String(), req.Points.String(), req.Bonus.String(), req.MaxViewers)
	if CheckErr(err, c) {
		return
	}
	now := time.Now().Format(time.RFC3339)
	task := common.ShareTask{
		Id:         ret.InsertId(),
		Title:      req.Title,
		Summary:    req.Summary,
		Link:       req.Link,
		Image:      req.Image,
		Points:     req.Points,
		PointsLeft: req.Points,
		Bonus:      req.Bonus,
		MaxViewers: req.MaxViewers,
		InsertedAt: now,
		UpdatedAt:  now,
		Creator:    user.Id,
	}
	task.ShareLink, _ = task.GetShareLink(deviceId, Config)
	c.JSON(http.StatusOK, task)
}
