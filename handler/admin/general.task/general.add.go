package general_task

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

func AddGeneralTaskHandler(c *gin.Context) {
	var task common.GeneralTask
	if CheckErr(c.Bind(&task), c) {
		return
	}

	query := `INSERT INTO tmm.general_tasks(creator,title,summary,image,details,points,points_left,bonus,online_status)
	VALUES(0,'%s','%s','%s','%s','%s','%s','%s',-1)`

	var db = Service.Db
	_, _, err := db.Query(query, db.Escape(task.Title), db.Escape(task.Summary), db.Escape(task.Image), db.Escape(task.Details),
		task.Points.String(), task.Points.String(), task.Bonus.String())
	if CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Data:    nil,
		Code:    0,
	})
}
