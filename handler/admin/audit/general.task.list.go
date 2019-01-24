package audit

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"github.com/tokenme/tmm/common"
	"github.com/shopspring/decimal"
	"strings"
)

type GeneralTask struct {
	common.GeneralTask
	DeviceId           string   `json:"device_id"`
	UserId             uint64   `json:"user_id"`
	Nick               string   `json:"nick"`
	Mobile             string   `json:"mobile"`
	Blocked            int      `json:"blocked"`
	Avatar             string   `json:"avatar"`
	CertificateImages  []string `json:"certificate_images,omitempty"`
}

func AuditGeneralTaskListHandler(c *gin.Context) {
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
	dg.device_id,
	dg.task_id,
	dg.status,
	dg.cert_info,
	dg.cert_images,
	dg.cert_comments,
	dg.inserted_at,
	u.mobile,
	u.id,
	IFNULL(wx.nick,u.nickname),
	IFNULL(wx.avatar,u.avatar),
	gt.title,
	gt.details,
	gt.image,
	gt.points,
	gt.points_left,
	gt.bonus,
	IF(us.user_id > 0,IF(us.blocked = us.block_whitelist,0,1),0) AS blocked

FROM 
	tmm.device_general_tasks AS dg
INNER JOIN tmm.devices AS dev ON (dev.id = dg.device_id)
INNER JOIN ucoin.users AS u ON (u.id = dev.user_id)
INNER JOIN general_tasks AS gt ON (gt.id = dg.task_id)
LEFT  JOIN tmm.wx AS wx ON(wx.user_id = dev.user_id)
LEFT  JOIN tmm.user_settings AS us ON (us.user_id = dev.user_id)
WHERE dg.status = %d
ORDER BY dg.inserted_at DESC 
LIMIT %d OFFSET %d
`

	db := Service.Db
	rows, _, err := db.Query(query, req.CertificateStatus, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	var data []*GeneralTask

	for _, row := range rows {
		task := &GeneralTask{}
		task.DeviceId = row.Str(0)
		task.Id = row.Uint64(1)
		task.CertificateStatus = int8(row.Int(2))
		task.CertificateInfo = row.Str(3)
		task.CertificateImages = strings.Split(row.Str(4), `,`)
		task.CertificateComment = row.Str(5)
		task.InsertedAt = row.Str(6)
		task.Mobile = row.Str(7)
		task.UserId = row.Uint64(8)
		task.Nick = row.Str(9)
		task.Avatar = row.Str(10)
		task.Title = row.Str(11)
		task.Details = row.Str(12)
		task.Image = row.Str(13)
		task.Points = decimal.NewFromFloat(row.Float(14))
		task.PointsLeft = decimal.NewFromFloat(row.Float(15))
		task.Bonus = decimal.NewFromFloat(row.Float(16))
		task.Blocked = row.Int(17)
		data = append(data, task)
	}

	var total int
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.device_general_tasks WHERE status = %d`, req.CertificateStatus)
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
