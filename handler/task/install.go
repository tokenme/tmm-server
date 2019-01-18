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
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
	BundleId string          `json:"bundle_id" form:"bundle_id" binding:"required"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	TaskId   uint64          `json:"task_id" form:"task_id" binding:"required"`
	Status   int8            `json:"status" form:"status"`
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
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if CheckWithCode(len(deviceId) == 0, NOTFOUND_ERROR, "invalid device", c) {
		return
	}
	rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "You have been finished the task", c) {
		return
	}
	if req.Status != -1 {
		rows, _, err = db.Query(`SELECT 1 FROM tmm.app_tasks AS appt WHERE id=%d AND platform='%s' AND bundle_id='%s' AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId, db.Escape(req.Platform), db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "Task not avaliable", c) {
			return
		}
	}
	var bonus decimal.Decimal
	if req.Status == 0 {
		_, _, err = db.Query(`INSERT IGNORE INTO tmm.device_app_tasks (device_id, task_id, bundle_id) VALUES ('%s', %d, '%s')`, db.Escape(deviceId), req.TaskId, db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
    } else if req.Status == 2 {
        rows, _, err := db.Query(`SELECT 1 FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d AND status=-1 LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) > 0, "You have been finished the task", c) {
			return
		}
        _, _, err = db.Query(`UPDATE tmm.device_app_tasks SET status = 2 WHERE device_id='%s' AND task_id = %d`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
	} else if Check(req.Status == 1, "请先升级", c) {
		return
		appTask := common.AppTask{Id: req.TaskId}
		bonus, err = appTask.Install(user, deviceId, Service, Config)
		if CheckErr(err, c) {
			return
		}
	} else if req.Status == -1 {
		appTask := common.AppTask{Id: req.TaskId}
		bonus, err = appTask.Uninstall(user, deviceId, Service, Config)
		if CheckErr(err, c) {
			return
		}
	}
	task := common.AppTask{
		Id:            req.TaskId,
		BundleId:      req.BundleId,
		InstallStatus: req.Status,
		Bonus:         bonus,
	}
	c.JSON(http.StatusOK, task)
}
