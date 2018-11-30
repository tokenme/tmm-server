package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

type AppAddRequest struct {
	Name        string          `json:"name" form:"name" binding:"required"`
	BundleId    string          `json:"bundle_id" form:"bundle_id" binding:"required"`
	Points      decimal.Decimal `json:"points" form:"points" binding:"required"`
	Bonus       decimal.Decimal `json:"bonus" form:"bonus" binding:"required"`
    DownloadUrl string          `json:"download_url" from:"download_url" binding:"required"`
    Icon        string          `json:"icon" from:"icon" binding:"required"`
	Platform    common.Platform `json:"platform" form:"platform" binding:"required"`
}

func AppAddHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppAddRequest
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
	_, ret, err := db.Query(`UPDATE tmm.devices AS d SET d.points = d.points - %s, d.consumed_ts = d.consumed_ts + %d WHERE id='%s' AND d.points >= %s AND d.user_id=%d`, req.Points.String(), ts.IntPart(), db.Escape(deviceId), req.Points.String(), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(ret.AffectedRows() == 0, NOT_ENOUGH_POINTS_ERROR, "not enough points in device", c) {
		return
	}
	_, ret, err = db.Query(`INSERT INTO tmm.app_tasks (creator, platform, name, bundle_id, download_url, icon, points, points_left, bonus) VALUES (%d, '%s', '%s', '%s', '%s', '%s', %s, %s, %s)`, user.Id, db.Escape(req.Platform), db.Escape(req.Name), db.Escape(req.BundleId), db.Escape(req.DownloadUrl), db.Escape(req.Icon), req.Points.String(), req.Points.String(), req.Bonus.String())
	if CheckErr(err, c) {
		return
	}
	_, schemeRet, err := db.Query(`INSERT IGNORE INTO tmm.app_scheme_ids (bundle_id) VALUES ('%s')`, db.Escape(req.BundleId))
	if CheckErr(err, c) {
		return
	}
	schemeId := schemeRet.InsertId()
	if schemeId == 0 {
		rows, _, err := db.Query(`SELECT id FROM tmm.app_scheme_ids WHERE bundle_id='%s' LIMIT 1`, db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			schemeId = rows[0].Uint64(0)
		}
	}
	now := time.Now().Format(time.RFC3339)
	task := common.AppTask{
		Id:         ret.InsertId(),
		Platform:   req.Platform,
		Name:       req.Name,
		SchemeId:   schemeId,
		BundleId:   req.BundleId,
		Points:     req.Points,
		PointsLeft: req.Points,
		Bonus:      req.Bonus,
		InsertedAt: now,
		UpdatedAt:  now,
		Creator:    user.Id,
	}
	lookup, err := common.App{BundleId: task.BundleId}.LookUp(Service)
	if err == nil {
		task.StoreId = lookup.TrackId
		task.Icon = lookup.ArtworkUrl512
	}
	c.JSON(http.StatusOK, task)
}
