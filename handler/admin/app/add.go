package app

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

func AddShareAppHandler(c *gin.Context) {
	var db = Service.Db
	var req admin.ShareAppRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	query := `INSERT INTO tmm.app_tasks (creator, platform, bundle_id, store_id, name, bonus, 
				download_url, icon, points, size, points_left) 
				VALUES (0, 'android', '%s', 0, '%s', %s, '%s', '%s', %s, %s, %s) `
	_, _, err := db.Query(query, db.Escape(req.BundleId), db.Escape(req.Title), req.Bonus.String(),
		db.Escape(req.Link), db.Escape(req.Image), req.Points.String(), db.Escape(req.Size), req.Points.String())
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Data:    nil,
		Code:    0,
	})
}