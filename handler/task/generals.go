package task

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

type GeneralsRequest struct {
	Page     uint            `json:"page" form:"page"`
	PageSize uint            `json:"page_size" form:"page_size"`
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
    TaskId   uint64          `json:"task_id" form:"task_id"`
	MineOnly bool            `json:"mine_only" form:"mine_only"`
}

func GeneralsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req GeneralsRequest
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
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	db := Service.Db

	onlineStatusConstrain := "AND g.points_left>0 AND g.online_status=1"
	//onlineStatusConstrain := fmt.Sprintf("AND a.points_left>0 AND a.online_status=1 AND NOT EXISTS (SELECT 1 FROM tmm.device_general_tasks AS dat WHERE dat.task_id=a.id AND dat.device_id='%s' AND dat.status=1 LIMIT 1)", db.Escape(deviceId))
	orderBy := "g.bonus DESC, g.id DESC"
	if req.MineOnly {
		onlineStatusConstrain = fmt.Sprintf("AND g.creator = %d", user.Id)
		orderBy = "g.id DESC"
	}

    if req.TaskId > 0 {
        onlineStatusConstrain += fmt.Sprintf(" AND g.id = %d ", req.TaskId)
    }

	query := `SELECT
    g.id,
    g.creator,
    g.title,
    g.summary,
    g.image,
    g.details,
    g.points,
    g.points_left,
    g.bonus,
    g.online_status,
    g.inserted_at,
    g.updated_at,
    dgt.status,
    dgt.cert_images,
    dgt.cert_comments,
    dgt.cert_info
FROM tmm.general_tasks AS g
LEFT JOIN tmm.device_general_tasks AS dgt ON (dgt.task_id=g.id AND dgt.device_id='%s')
WHERE 1=1 %s
ORDER BY %s LIMIT %d, %d`
	rows, _, err := db.Query(query, db.Escape(deviceId), onlineStatusConstrain, orderBy, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var tasks []common.GeneralTask
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(8))
		points, _ := decimal.NewFromString(row.Str(6))
		pointsLeft, _ := decimal.NewFromString(row.Str(7))
		creator := row.Uint64(1)
		task := common.GeneralTask{
			Id:                 row.Uint64(0),
			Title:              row.Str(2),
			Summary:            row.Str(3),
			Image:              row.Str(4),
			Details:            row.Str(5),
			Bonus:              bonus,
			Points:             points,
			PointsLeft:         pointsLeft,
			OnlineStatus:       int8(row.Int(9)),
			InsertedAt:         row.ForceLocaltime(10).Format(time.RFC3339),
			UpdatedAt:          row.ForceLocaltime(11).Format(time.RFC3339),
            CertificateStatus:  int8(row.Int(12)),
            CertificateImages:  row.Str(13),
            CertificateComment: row.Str(14),
            CertificateInfo:    row.Str(15),
		}
		if creator == user.Id {
			task.Creator = creator
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
