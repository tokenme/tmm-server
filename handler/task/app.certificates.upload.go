package task

import (
	//"github.com/davecgh/go-spew/spew"
	//"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type AppCertificatesUploadRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
	BundleId string          `json:"bundle_id" form:"bundle_id" binding:"required"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	TaskId   uint64          `json:"task_id" form:"task_id" binding:"required"`
    Images   string          `json:"images" form:"images" binding:"required"`
}

func AppCertificatesUploadHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppCertificatesUploadRequest
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
    if CheckWithCode(len(req.Images) == 0, NOTFOUND_ERROR, "no images", c) {
        return
    }
    rows, _, err := db.Query(`SELECT 1 FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d AND status=-1 LIMIT 1`, db.Escape(deviceId), req.TaskId)
    if CheckErr(err, c) {
        return
    }
    if Check(len(rows) > 0, "You have been finished the task", c) {
        return
    }
    rows, _, err = db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "Wrong device id", c) {
		return
	}
    rows, _, err = db.Query(`SELECT 1 FROM tmm.app_tasks AS appt WHERE id=%d AND platform='%s' AND bundle_id='%s' AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId, db.Escape(req.Platform), db.Escape(req.BundleId))
    if CheckErr(err, c) {
        return
    }
    if Check(len(rows) == 0, "Task not avaliable", c) {
        return
    }
    _, _, err = db.Query(`INSERT INTO tmm.device_app_task_certificates (device_id, task_id, bundle_id, images, status) VALUES ('%s', %d, '%s', '%s', 1) ON DUPLICATE KEY UPDATE images=VALUES(images), status=VALUES(status)`, db.Escape(deviceId), req.TaskId, db.Escape(req.BundleId), db.Escape(req.Images))
    if CheckErr(err, c) {
        return
    }
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
