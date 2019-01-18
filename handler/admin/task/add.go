package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/tokenme/tmm/tools/videospider"
	"net/http"
	"strings"
)

func AddShareHandler(c *gin.Context) {
	var req admin.AddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var fieldList []string
	var valueList []string
	var db = Service.Db
	if req.Title != "" {
		fieldList = append(fieldList, `title`)
		valueList = append(valueList, fmt.Sprintf(`'%s'`, db.Escape(req.Title)))
	}
	if req.Summary != "" {
		fieldList = append(fieldList, `summary`)
		valueList = append(valueList, fmt.Sprintf(`'%s'`, db.Escape(req.Summary)))
	}
	if req.Image != "" {
		fieldList = append(fieldList, `image`)
		valueList = append(valueList, fmt.Sprintf(`'%s'`, db.Escape(req.Image)))
	}
	if req.Points.String() != "0" {
		fieldList = append(fieldList, `points,points_left`)
		valueList = append(valueList, fmt.Sprintf(`'%s','%s'`, req.Points.String(), req.Points.String()))
	}
	if req.Bonus.String() != "0" {
		fieldList = append(fieldList, `bonus`)
		valueList = append(valueList, fmt.Sprintf(`'%s'`, req.Bonus.String()))
	}
	if req.MaxViewers != 0 {
		fieldList = append(fieldList, `max_viewers`)
		valueList = append(valueList, fmt.Sprintf(`%d`, req.MaxViewers))
	}
	if req.Link != "" {
		Video := videospider.NewClient(Service, Config)
		if video, err := Video.Get(req.Link); err == nil {
			if Video.Save(video) == nil {
				c.JSON(http.StatusOK, admin.Response{
					Code:    0,
					Message: admin.API_OK,
					Data:    video,
				})
				return
			}
		}
		if Check(len(fieldList) == 0, `Error Link`, c) {
			return
		}
		fieldList = append(fieldList, `link`)
		valueList = append(valueList, fmt.Sprintf(`'%s'`, req.Link))
	}

	if Check(len(fieldList) == 0 || len(fieldList) != len(valueList), `Invalid param`, c) {
		return
	}
	fieldList = append(fieldList, `creator`)
	valueList = append(valueList, fmt.Sprintf(`%d`, 0))

	var query = `
	INSERT INTO tmm.share_tasks 
	(%s) 
	VALUES (%s)`


	_, res, err := db.Query(query, strings.Join(fieldList, `,`), strings.Join(valueList, `,`))
	if CheckErr(err, c) {
		return
	}

	var insertArray []string
	for _, cid := range req.Cid {
		insertArray = append(insertArray, fmt.Sprintf(`(%d,%d,%d)`, res.InsertId(), cid, 1))
	}
	if len(insertArray) > 0 {
		cidInsert := `INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUE %s  ON 
		DUPLICATE KEY UPDATE cid=VALUES(cid),is_auto=VALUES(is_auto)`
		if _, _, err := db.Query(cidInsert, strings.Join(insertArray, `,`)); CheckErr(err, c) {
			return
		}
	}

	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Data:    req,
		Code:    0,
	},
	)
}
