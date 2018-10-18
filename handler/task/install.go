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
	"strings"
)

type AppInstallRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	BundleId string          `json:"bundle_id" form:"bundle_id" binding:"required"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	TaskId   uint64          `json:"task_id" form:"task_id" binding:"required"`
	Status   int8            `json:"status" form:"status"`
}

func AppInstallHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req AppInstallRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	device := common.DeviceRequest{
		Idfa:     req.Idfa,
		Platform: req.Platform,
	}
	deviceId := device.DeviceId()
	rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	if req.Status != -1 {
		rows, _, err = db.Query(`SELECT 1 FROM tmm.app_tasks AS appt WHERE id=%d AND platform='%s' AND bundle_id='%s' AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId, db.Escape(req.Platform), db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
	}
	var bonus decimal.Decimal
	if req.Status == 0 {
		_, _, err = db.Query(`INSERT IGNORE INTO tmm.device_app_tasks (device_id, task_id, bundle_id) VALUES ('%s', %d, '%s')`, db.Escape(deviceId), req.TaskId, db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
	} else if req.Status == 1 {
		pointsPerTs, err := common.GetPointsPerTs(Service)
		if CheckErr(err, c) {
			return
		}
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = d.points + IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left),
    d.total_ts = d.total_ts + CEIL(IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) / %s),
    appt.points_left = IF(appt.points_left > appt.bonus, appt.points_left - appt.bonus, 0),
    appt.downloads = appt.downloads + 1,
    dat.points = IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left),
    dat.status = 1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status != 1`
		_, _, err = db.Query(query, pointsPerTs.String(), db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))

		query = `SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
WHERE ic.user_id = %d
ORDER BY d.lastping_at DESC LIMIT 1) AS t1
UNION
SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
WHERE ic.user_id = %d
ORDER BY d.lastping_at DESC LIMIT 1) AS t2`
		rows, _, err = db.Query(query, user.Id, user.Id)
		if CheckErr(err, c) {
			return
		}
		var (
			inviterBonus = bonus.Mul(decimal.NewFromFloat(Config.InviteBonusRate))
			deviceIds    []string
			insertLogs   []string
		)
		for _, row := range rows {
			deviceIds = append(deviceIds, fmt.Sprintf("'%s'", db.Escape(row.Str(0))))
			insertLogs = append(insertLogs, fmt.Sprintf("(%d, %d, %s, 2, %d)", row.Uint64(1), user.Id, inviterBonus.String(), req.TaskId))
		}
		if len(deviceIds) > 0 {
			_, ret, err := db.Query(`UPDATE tmm.devices SET points = points + %s WHERE id IN (%s)`, inviterBonus.String(), strings.Join(deviceIds, ","))
			if CheckErr(err, c) {
				return
			}
			if ret.AffectedRows() > 0 {
				_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, task_type, task_id) VALUES %s`, strings.Join(insertLogs, ","))
				if CheckErr(err, c) {
					return
				}
			}
		}

	} else if req.Status == -1 {
		rows, _, err := db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "not found", c) {
			return
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))
		pointsPerTs, err := common.GetPointsPerTs(Service)
		if CheckErr(err, c) {
			return
		}
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = IF(d.points > dat.points, d.points - dat.points, 0),
    d.consumed_ts = d.consumed_ts + CEIL(IF(d.points > dat.points, d.points - dat.points, 0) / %s),
    appt.points_left = appt.points_left + IF(d.points > dat.points, dat.points, 0),
    appt.downloads = appt.downloads - IF(d.points > dat.points, 1, 0),
    dat.status = -1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status = 1`
		_, _, err = db.Query(query, pointsPerTs, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}

		query = `SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
WHERE ic.user_id = %d
ORDER BY d.points DESC LIMIT 1) AS t1
UNION
SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
WHERE ic.user_id = %d
ORDER BY d.points DESC LIMIT 1) AS t2`
		rows, _, err = db.Query(query, user.Id, user.Id)
		if CheckErr(err, c) {
			return
		}
		var (
			deviceIds  []string
			inviterIds []string
		)
		for _, row := range rows {
			deviceIds = append(deviceIds, fmt.Sprintf("'%s'", db.Escape(row.Str(0))))
			inviterIds = append(inviterIds, fmt.Sprintf("%d", row.Uint64(1)))
		}
		_, ret, err := db.Query(`UPDATE tmm.devices AS d, tmm.invite_bonus AS ib SET d.points = IF(d.points>ib.bonus, d.points-ib.bonus, 0) WHERE ib.user_id IN (%s) AND d.id IN (%s) AND ib.task_type=2 AND ib.task_id=%d`, strings.Join(inviterIds, ","), strings.Join(deviceIds, ","), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if ret.AffectedRows() > 0 {
			_, _, err = db.Query(`DELETE FROM tmm.invite_bonus WHERE user_id IN (%s) AND ib.task_type=2 AND ib.task_id=%d`, strings.Join(inviterIds, ","), req.TaskId)
			if CheckErr(err, c) {
				return
			}
		}
	}
	task := common.AppTask{
		Id:            req.TaskId,
		BundleId:      req.BundleId,
		InstallStatus: req.Status,
		Bonus:         bonus,
	}
	c.JSON(http.StatusOK, task)
}
