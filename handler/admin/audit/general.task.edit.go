package audit

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/tokenme/tmm/common"
)

func EditGeneralTaskHandler(c *gin.Context) {
	var task GeneralTask
	if CheckErr(c.Bind(&task), c) {
		return
	}

	if task.CertificateStatus == 1 {
		var user common.User
		user.Id = task.UserId
		_, err := task.CompleteTask(user, task.DeviceId, Service, Config)
		if CheckErr(err, c) {
			return
		}
	} else {
		query := `UPDATE tmm.device_general_tasks SET status = %d,cert_comments="%s" WHERE device_id = '%s' AND task_id = %d AND status = 0`
		db := Service.Db
		if _, _, err := db.Query(query, task.CertificateStatus, db.Escape(task.CertificateComment), db.Escape(task.DeviceId), task.Id); CheckErr(err, c) {
			return
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	},
	)

}
