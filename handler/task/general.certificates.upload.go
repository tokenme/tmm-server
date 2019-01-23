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

type GeneralCertificatesUploadRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
	TaskId   uint64          `json:"task_id" form:"task_id" binding:"required"`
    Images   string          `json:"images" form:"images" binding:"required"`
}

func GeneralCertificatesUploadHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req GeneralCertificatesUploadRequest
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
    rows, _, err := db.Query(`SELECT 1 FROM tmm.device_general_tasks WHERE device_id='%s' AND task_id=%d AND status=1 LIMIT 1`, db.Escape(deviceId), req.TaskId)
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
    rows, _, err = db.Query(`SELECT 1 FROM tmm.general_tasks WHERE id=%d AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId)
    if CheckErr(err, c) {
        return
    }
    if Check(len(rows) == 0, "Task not avaliable", c) {
        return
    }
    _, _, err = db.Query(`INSERT INTO tmm.device_general_tasks (device_id, task_id, cert_images, status) VALUES ('%s', %d, '%s', 0) ON DUPLICATE KEY UPDATE cert_images=VALUES(cert_images), status=VALUES(status)`, db.Escape(deviceId), req.TaskId, db.Escape(req.Images))
    if CheckErr(err, c) {
        return
    }
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
