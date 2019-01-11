package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strconv"
	"strings"
)

func ModifyTaskHandler(c *gin.Context) {
	var task common.ShareTask
	if CheckErr(c.Bind(&task), c) {
		return
	}

	var query []string
	var db = Service.Db
	if task.Title != "" {
		query = append(query, fmt.Sprintf(`title = '%s'`, db.Escape(task.Title)))
	}
	if task.Summary != "" {
		query = append(query, fmt.Sprintf(`summary = '%s'`, db.Escape(task.Summary)))
	}
	if task.Link != "" {
		query = append(query, fmt.Sprintf(`link = '%s' `, db.Escape(task.Link)))
	}
	if task.Image != "" {
		query = append(query, fmt.Sprintf(`image='%s'`, db.Escape(task.Image)))
	}
	if task.MaxViewers != 0 {
		query = append(query, fmt.Sprintf(`max_viewers=%d`, task.MaxViewers))
	}
	if task.Points.String() != "0" {
		query = append(query, fmt.Sprintf(`points='%s'`, db.Escape(task.Points.String())))
		query = append(query, fmt.Sprintf(`points_left= %s `, db.Escape(task.Points.String())))
	}
	if task.Bonus.String() != "0" {
		query = append(query, fmt.Sprintf(`bonus='%s'`, db.Escape(task.Bonus.String())))
	}
	if task.OnlineStatus == -1  || task.OnlineStatus == 1 {
		query = append(query, fmt.Sprintf(`online_status = %d`, task.OnlineStatus))
	}
	if Check(task.Id == 0, `Id Not Can Empty `, c) {
		return
	}

	if len(query) > 0 {
		if _, _, err := db.Query(`UPDATE tmm.share_tasks SET %s WHERE id = %d`, strings.Join(query, `,`), task.Id); CheckErr(err, c) {
			return
		}
	}

	if task.Cid != nil {
		var (
			lableArray  []string
			insertArray []string
		)

		for _, cid := range task.Cid {
			lableArray = append(lableArray, strconv.Itoa(cid))
			insertArray = append(insertArray, fmt.Sprintf(`(%d,%d,%d)`, task.Id, cid, 1))
		}

		if len(insertArray) > 0 {
			var cidInsert = `INSERT INTO tmm.share_task_categories 
			(task_id,cid,is_auto) VALUE %s  
			ON DUPLICATE KEY UPDATE cid=VALUES(cid),is_auto=VALUES(is_auto)`

			if _, _, err := db.Query(cidInsert, strings.Join(insertArray, `,`)); CheckErr(err, c) {
				return
			}
		}

		var delQuery = `DELETE FROM share_task_categories WHERE cid NOT IN (%s) AND task_id = %d`
		if len(lableArray) > 0 {
			if _, _, err := db.Query(delQuery, strings.Join(lableArray, `,`), task.Id); CheckErr(err, c) {
				return
			}
		} else {
			if _, _, err := db.Query(`DELETE FROM share_task_categories WHERE task_id = %d`, task.Id); CheckErr(err, c) {
				return
			}
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	})
}
