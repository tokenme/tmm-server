package general_task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
)

func EditGeneralTaskHandler(c *gin.Context) {
	var task common.GeneralTask
	if CheckErr(c.Bind(&task), c) {
		return
	}

	db := Service.Db
	var fieldList []string
	if task.Title != "" {
		fieldList = append(fieldList, fmt.Sprintf(`title = '%s'`, db.Escape(task.Title)))
	}
	if task.Summary != "" {
		fieldList = append(fieldList, fmt.Sprintf(`summary = '%s' `, db.Escape(task.Summary)))
	}
	if task.Image != "" {
		fieldList = append(fieldList, fmt.Sprintf(`image = '%s' `, db.Escape(task.Image)))
	}
	if task.Details != "" {
		fieldList = append(fieldList, fmt.Sprintf(`details = '%s' `, db.Escape(task.Details)))
	}
	if task.Points.String() != "0" {
		fieldList = append(fieldList, fmt.Sprintf(`points = '%s' `, task.Points.String()))
	}
	if task.PointsLeft.String() != "0" {
		fieldList = append(fieldList, fmt.Sprintf(`points_left = '%s' `, task.PointsLeft.String()))
	}
	if task.OnlineStatus != 0 {
		fieldList = append(fieldList, fmt.Sprintf(`online_status = %d `, task.OnlineStatus))
	}
	if task.Bonus.String() != "0" {
		fieldList = append(fieldList, fmt.Sprintf(" bonus = '%s' ", task.Bonus.String()))
	}

	if len(fieldList) > 0 {
		_, _, err := db.Query(`UPDATE tmm.general_tasks SET %s   WHERE id = %d`, strings.Join(fieldList, `,`), task.Id)
		if CheckErr(err, c) {
			return
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
	})

}
