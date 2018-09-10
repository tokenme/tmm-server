package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type AppInstallRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	BundleId string          `json:"bundle_id" form:"bundle_id" binding:"required"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	TaskId   uint64          `json:"task_id" form:"task_id" binding:"required"`
	Status   int             `json:"status" form:"status"`
}

func AppInstallHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppInstallRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	device := common.DeviceRequest{
		Idfa:     req.Idfa,
		Platform: req.Platform,
	}
	deviceId := device.DeviceId()
	rows, _, err := db.Query(`SELECT 1 FROM tmm.user_devices AS ud WHERE user_id=%d AND device_id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	if req.Status != -1 {
		rows, _, err = db.Query(`SELECT 1 FROM tmm.app_tasks AS appt WHERE id=%d AND platform='%s' AND bundle_id='%s' AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId, db.Escape(req.Platform), db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
	}
	var bonus decimal.Decimal
	if req.Status == 0 {
		_, _, err = db.Query(`INSERT IGNORE INTO tmm.device_app_tasks (device_id, task_id, bundle_id) VALUES ('%s', %d, '%s')`, db.Escape(deviceId), req.TaskId, db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
	} else if req.Status == 1 {
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = d.points + IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left),
    appt.points_left = IF(appt.points_left > appt.bonus, appt.points_left - appt.bonus, 0),
    appt.downloads = appt.downloads + 1,
    dat.points = IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left),
    dat.status = 1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status != 1`
		_, _, err = db.Query(query, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))
	} else if req.Status == -1 {
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = IF(d.points > dat.points, d.points - dat.points, 0),
    appt.points_left = appt.points_left + IF(d.points > dat.points, dat.points, 0),
    appt.downloads = appt.downloads - IF(d.points > dat.points, 1, 0),
    dat.status = -1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status = 1`
		_, _, err = db.Query(query, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
	}
	task := common.AppTask{
		Id:       req.TaskId,
		BundleId: req.BundleId,
		Status:   req.Status,
		Bonus:    bonus,
	}
	c.JSON(http.StatusOK, task)
}
