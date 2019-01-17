package app

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func AddShareAppHandler(c *gin.Context) {
	var db = Service.Db
	var req common.AppTask
	if CheckErr(c.Bind(&req), c) {
		return
	}
	query := `INSERT INTO tmm.app_tasks (creator, platform, bundle_id, store_id, name, bonus, 
				download_url, icon, points, size, points_left,details) 
				VALUES (0, 'android', '%s', 0, '%s', %s, '%s', '%s', %s, %d, %s,'%s') `
	_, _, err := db.Query(query, db.Escape(req.BundleId), db.Escape(req.Name), req.Bonus.String(),
		db.Escape(req.DownloadUrl), db.Escape(req.Icon), req.Points.String(), req.Size, req.Points.String(),req.Details)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Data:    nil,
		Code:    0,
	})
}