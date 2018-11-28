package task

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"github.com/tokenme/tmm/common"
)

func ModifyTaskHandler(c *gin.Context) {
	var (
		task      common.ShareTask
		query     string
		cidQuery  = `select cid from tmm.share_task_categories where task_id = %d`
		cidInsert = `Insert INTO tmm.share_task_categories (task_id,cid,is_auto) VALUE(%d,%d,%d)`
		delQuery  = `delete from share_task_categories where cid not in (%s) and task_id = %d`
		array     []string
		db        = Service.Db
		cliList   = make(map[int]struct{})
	)
	if CheckErr(c.Bind(&task), c) {
		return
	}
	if task.Title != "" && task.Summary != "" && task.Link != "" && task.Image != "" && task.MaxViewers != 0 {
		query = fmt.Sprintf(`update tmm.share_tasks set title = '%s' , summary='%s' ,link='%s' , image='%s' ,
		points='%s',bonus='%s',max_viewers=%d where id = %d`, task.Title, task.Summary, task.Link,
			task.Image, task.Points.String(), task.Bonus.String(), task.MaxViewers, task.Id)
	}
	if task.OnlineStatus != 0 {
		if Check(task.OnlineStatus != -1 && task.OnlineStatus != 1, `OnlineStatus only can be 1 or -1`, c) {
			return
		}
		query = fmt.Sprintf(`update tmm.share_tasks set online_status = %d where id = %d`, task.OnlineStatus, task.Id)
	}

	if query != "" {
		_, _, err := db.Query(query)
		if CheckErr(err, c) {
			return
		}
	}
	if task.Cid != nil {
		rows, _, err := db.Query(cidQuery, task.Id)
		if CheckErr(err, c) {
			return
		}
		for _, row := range rows {
			cliList[row.Int(0)] = struct{}{}
		}
		for _, cli := range task.Cid {
			array = append(array, strconv.Itoa(cli))
			if _, ok := cliList[cli]; ok {
				continue
			} else {
				_, _, err := db.Query(cidInsert, task.Id, cli, 1)
				if CheckErr(err, c) {
					return
				}
			}
		}
		in := strings.Join(array, `,`)
		if in != "" {
			_, _, err = db.Query(delQuery, in, task.Id)
		} else {
			_, _, err = db.Query(`delete from share_task_categories where task_id = %d`, task.Id)
		}
		if CheckErr(err, c) {
			return
		}

	}
	c.JSON(http.StatusOK, APIResponse{Msg: `ok`})
}
