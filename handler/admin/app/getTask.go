package app

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"strconv"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func GetAppTaskHandler(c *gin.Context) {
	var db = Service.Db
	taskid, err := strconv.Atoi(c.Query(`taskid`))
	if CheckErr(err, c) {
		return
	}
	query := `SELECT id, bundle_id, name, size, bonus, download_url, icon, points,details
				FROM tmm.app_tasks 
				WHERE id = %d  `
	rows, res, err := db.Query(query, taskid)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `Not found `, c) {
		return
	}
	row := rows[0]
	points, err := decimal.NewFromString(row.Str(res.Map(`points`)))
	if CheckErr(err, c) {
		return
	}
	bonus, err := decimal.NewFromString(row.Str(res.Map(`bonus`)))
	if CheckErr(err, c) {
		return
	}
	var task common.AppTask
	task.Id = uint64(taskid)
	task.BundleId = row.Str(res.Map(`bundle_id`))
	task.Name = row.Str(res.Map(`name`))
	task.Size = row.Uint(res.Map(`size`))
	task.Bonus = bonus
	task.DownloadUrl = row.Str(res.Map(`download_url`))
	task.Icon = row.Str(res.Map(`icon`))
	task.Details = row.Str(res.Map(`details`))
	task.Points = points
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    task,
	})
}