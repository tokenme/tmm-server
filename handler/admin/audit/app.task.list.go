package audit

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
	"fmt"
)

type AppTask struct {
	common.AppTask
	Images     []string `json:"images,omitempty"`
	Status     int      `json:"status,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	DeviceId   string   `json:"device_id,omitempty"`
	Points     string   `json:"points,omitempty"`
	PointsLeft string   `json:"points_left,omitempty"`
	Bonus      string   `json:"bonus,omitempty"`
	UserId     uint64     `json:"user_id,omitempty"`
	Nick       string   `json:"nick,omitempty"`
	Mobile     string   `json:"mobile,omitempty"`
	Avatar     string   `json:"avatar,omitempty"`
	Blocked    int      `json:"blocked"`
}

type Request struct {
	admin.Pages
	Status int `form:"status"`
}

func AuditAppTaskListHandler(c *gin.Context) {
	var req Request
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var offset int
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Page > 1 {
		offset = (req.Page - 1) * req.Limit
	}

	query := `
SELECT 
	dc.device_id,
	dc.task_id,
	dc.bundle_id,
	dc.inserted_at,
	dc.comment,
	dc.images,
	dc.status,
	at.name,
	at.icon,
	at.details,
	at.bonus,
	at.download_url,
	at.points,
	at.points_left,
	u.id,
	u.mobile,
	IFNULL(wx.avatar,u.avatar),
	IFNULL(wx.nick,u.nickname),
	IF(us.user_id > 0,IF(us.blocked = us.block_whitelist,0,1),0) AS blocked
FROM 
	tmm.device_app_task_certificates AS dc
INNER JOIN tmm.app_tasks AS at ON (at.id = dc.task_id)
INNER JOIN tmm.devices AS dev ON (dev.id = dc.device_id)
INNER JOIN ucoin.users AS u ON (u.id = dev.user_id)
LEFT  JOIN tmm.wx AS wx ON (wx.user_id = dev.user_id)
LEFT  JOIN tmm.user_settings AS us ON (us.user_id = dev.user_id)
WHERE dc.status = %d
ORDER BY dc.inserted_at DESC 
LIMIT %d OFFSET %d
`

	db := Service.Db
	rows, _, err := db.Query(query, req.Status, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var data []*AppTask

	for _, row := range rows {
		appTask := &AppTask{}
		appTask.DeviceId = row.Str(0)
		appTask.Id = row.Uint64(1)
		appTask.BundleId = row.Str(2)
		appTask.InsertedAt = row.Str(3)
		appTask.Comment = row.Str(4)
		appTask.Images = strings.Split(row.Str(5), `,`)
		appTask.Status = row.Int(6)
		appTask.Name = row.Str(7)
		appTask.Icon = row.Str(8)
		appTask.Details = row.Str(9)
		appTask.Bonus = fmt.Sprintf("%.2f", row.Float(10))
		appTask.DownloadUrl = row.Str(11)
		appTask.Points = fmt.Sprintf("%.2f", row.Float(12))
		appTask.PointsLeft = fmt.Sprintf("%.2f", row.Float(13))
		appTask.UserId = row.Uint64(14)
		appTask.Mobile = row.Str(15)
		appTask.Avatar = row.Str(16)
		appTask.Nick = row.Str(17)
		appTask.Blocked = row.Int(18)
		data = append(data, appTask)
	}

	var total int
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.device_app_task_certificates WHERE status = %d`, req.Status)
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			`total`: total,
			`data`:  data,
		},
	})
}
