package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
	"github.com/tokenme/tmm/handler/admin"
)

func ModifyShareAppHandler(c *gin.Context) {
	var db = Service.Db
	var task common.AppTask
	if CheckErr(c.Bind(&task), c) {
		return
	}
	var query []string
	if task.Name != "" {
		query = append(query, fmt.Sprintf(`name = '%s'`, db.Escape(task.Name)))
	}
	if task.BundleId != "" {
		query = append(query, fmt.Sprintf(`bundle_id = '%s'`, db.Escape(task.BundleId)))
	}
	if task.DownloadUrl != "" {
		query = append(query, fmt.Sprintf(`download_url = '%s'`, db.Escape(task.DownloadUrl)))
	}
	if task.Size != 0 {
		query = append(query, fmt.Sprintf(`size = %d`, task.Size))
	}
	if task.Points.String() != "0" {
		query = append(query, fmt.Sprintf(`points= %s `, db.Escape(task.Points.String())))
		query = append(query, fmt.Sprintf(`points_left= %s `, db.Escape(task.Points.String())))
	}
	if task.Bonus.String() != "0" {
		query = append(query, fmt.Sprintf(`bonus= %s `, db.Escape(task.Bonus.String())))
	}
	if task.Icon != "" {
		query = append(query, fmt.Sprintf(`icon = '%s'`, db.Escape(task.Icon)))
	}
	if task.OnlineStatus != 0 {
		if task.OnlineStatus != -1 && task.OnlineStatus != 1 {
			task.OnlineStatus = 1
		}
		query = append(query, fmt.Sprintf(`online_status = %d`, task.OnlineStatus))
	}
	if task.Details != "" {
		query = append(query, fmt.Sprintf(`details = '%s'`, db.Escape(task.Details)))
	}
	if Check(task.Id == 0, `Id Not Can Empty `, c) {
		return
	}
	if len(query) > 0 {
		_, _, err := db.Query(`UPDATE tmm.app_tasks SET %s WHERE id = %d`, strings.Join(query, `,`), task.Id)
		if CheckErr(err, c) {
			return
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	})
}
