package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

type AppsRequest struct {
	Page     uint            `json:"page" form:"page"`
	PageSize uint            `json:"page_size" form:"page_size"`
	Idfa     string          `json:"idfa" form:"idfa"`
    Imei     string          `json:"imei" form:"imei"`
    Mac      string          `json:"mac" form:"mac"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	MineOnly bool            `json:"mine_only" form:"mine_only"`
}

func AppsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}

	device := common.DeviceRequest{
		Idfa:     req.Idfa,
        Imei:     req.Imei,
        Mac:      req.Mac,
	}
	deviceId := device.DeviceId()
    if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	db := Service.Db

	onlineStatusConstrain := fmt.Sprintf("AND a.points_left>0 AND a.online_status=1 AND NOT EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.task_id=a.id AND dat.device_id='%s' AND dat.status=1 LIMIT 1)", db.Escape(deviceId))
	orderBy := "a.bonus DESC, a.id DESC"
	if req.MineOnly {
		onlineStatusConstrain = fmt.Sprintf("AND a.creator = %d", user.Id)
		orderBy = "a.id DESC"
	}

	query := `SELECT
    a.id,
    a.platform,
    a.name,
    a.bundle_id,
    a.store_id,
    a.bonus,
    a.points,
    a.points_left,
    a.downloads,
    a.inserted_at,
    a.updated_at,
    a.creator,
    IFNULL(asi.id, 0),
    a.online_status,
    a.download_url
FROM tmm.app_tasks AS a
LEFT JOIN tmm.app_scheme_ids AS asi ON (asi.bundle_id = a.bundle_id)
WHERE a.platform='%s' %s
ORDER BY %s LIMIT %d, %d`
	rows, _, err := db.Query(query, db.Escape(req.Platform), onlineStatusConstrain, orderBy, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var tasks []common.AppTask
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(5))
		points, _ := decimal.NewFromString(row.Str(6))
		pointsLeft, _ := decimal.NewFromString(row.Str(7))
		creator := row.Uint64(11)
		task := common.AppTask{
			Id:             row.Uint64(0),
			Platform:       row.Str(1),
			Name:           row.Str(2),
			BundleId:       row.Str(3),
			StoreId:        row.Uint64(4),
			Bonus:          bonus,
			Points:         points,
			PointsLeft:     pointsLeft,
			InsertedAt:     row.ForceLocaltime(9).Format(time.RFC3339),
			UpdatedAt:      row.ForceLocaltime(10).Format(time.RFC3339),
			SchemeId:       row.Uint64(12),
            DownloadUrl:    row.Str(14),
		}
		if creator == user.Id {
			task.Downloads = row.Uint(8)
			task.Creator = creator
			task.OnlineStatus = int8(row.Int(13))
		}
		if task.StoreId == 0 {
			lookup, err := common.App{BundleId: task.BundleId}.LookUp(Service)
			if err == nil {
				task.StoreId = lookup.TrackId
				task.Icon = lookup.ArtworkUrl512
			}
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
