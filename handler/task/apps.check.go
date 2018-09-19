package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type AppsCheckRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
}

func AppsCheckHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppsCheckRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	device := common.DeviceRequest{
		Idfa:     req.Idfa,
		Platform: req.Platform,
	}
	deviceId := device.DeviceId()
	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.user_devices AS ud WHERE user_id=%d AND device_id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	rows, _, err = db.Query(`SELECT dat.task_id, dat.bundle_id, dat.status, asi.id FROM tmm.device_app_tasks AS dat LEFT JOIN tmm.app_scheme_ids AS asi ON (asi.bundle_id=dat.bundle_id) WHERE dat.device_id='%s' AND dat.updated_at>DATE_SUB(NOW(), INTERVAL 7 DAY)`, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	var tasks []common.AppTask
	for _, row := range rows {
		task := common.AppTask{
			Id:       row.Uint64(0),
			BundleId: row.Str(1),
			Status:   row.Int(2),
			SchemeId: row.Uint64(3),
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
