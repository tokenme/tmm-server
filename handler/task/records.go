package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

type RecordsRequest struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

func RecordsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req RecordsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}

	db := Service.Db
	query := `SELECT
    1 AS task_type,
    appt.name AS title,
    dat.points AS points,
    dat.updated_at AS updated_at,
    appt.bundle_id AS bundle_id,
    0 AS viewers
FROM tmm.devices AS d
LEFT JOIN tmm.device_app_tasks AS dat ON (dat.device_id = d.id)
INNER JOIN tmm.app_tasks AS appt ON (appt.id = dat.task_id)
WHERE d.user_id = %d AND dat.status = 1
UNION
SELECT
    2 AS task_type,
    st.title AS title,
    dst.points AS points,
    dst.updated_at AS updated_at,
    '' AS bundle_id,
    st.viewers AS viewers
FROM tmm.devices AS d
LEFT JOIN tmm.device_share_tasks AS dst ON (dst.device_id = d.id)
INNER JOIN tmm.share_tasks AS st ON (st.id = dst.task_id)
WHERE d.user_id = %d
ORDER BY updated_at DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, user.Id, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var records []common.TaskRecord
	for _, row := range rows {
		points, _ := decimal.NewFromString(row.Str(2))
		rec := common.TaskRecord{
			Type:      row.Uint(0),
			Title:     row.Str(1),
			Points:    points,
			UpdatedAt: row.ForceLocaltime(3).Format(time.RFC3339),
		}
		if rec.Type == common.AppTaskType {
			bundleId := row.Str(4)
			lookup, err := common.App{BundleId: bundleId}.LookUp(Service)
			if err == nil {
				rec.Image = lookup.ArtworkUrl512
			}
		} else {
			rec.Viewers = row.Uint(5)
		}
		records = append(records, rec)
	}
	c.JSON(http.StatusOK, records)
}
