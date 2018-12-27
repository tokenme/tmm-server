package task

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
)

type AppInstallRequest struct {
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
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
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if CheckWithCode(len(deviceId) == 0, NOTFOUND_ERROR, "invalid device", c) {
		return
	}
	rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "You have been finished the task", c) {
		return
	}
	if req.Status != -1 {
		rows, _, err = db.Query(`SELECT 1 FROM tmm.app_tasks AS appt WHERE id=%d AND platform='%s' AND bundle_id='%s' AND online_status=1 AND points_left>=bonus LIMIT 1`, req.TaskId, db.Escape(req.Platform), db.Escape(req.BundleId))
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "Task not avaliable", c) {
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
		rows, _, err := db.Query(`SELECT 1 FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d AND status=-1 LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) > 0, "You have been finished the task", c) {
			return
		}

		pointsPerTs, err := common.GetPointsPerTs(Service)
		if CheckErr(err, c) {
			return
		}
		var bonusRate float64 = 1
		rows, _, err = db.Query(`SELECT ul.task_bonus_rate FROM tmm.user_settings AS us INNER JOIN tmm.user_levels AS ul ON (ul.id=us.level) INNER JOIN tmm.devices AS d ON (d.user_id=us.user_id) WHERE d.id='%s' LIMIT 1`, db.Escape(deviceId))
		if err != nil {
			log.Error(err.Error())
		} else if len(rows) > 0 {
			bonusRate = rows[0].ForceFloat(0) / 100
		}
		query := `UPDATE tmm.devices AS d, tmm.device_app_tasks AS dat, tmm.app_tasks AS appt
SET d.points = d.points + IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) * %.2f,
    d.total_ts = d.total_ts + CEIL(IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) / %s),
    appt.points_left = IF(appt.points_left > appt.bonus, appt.points_left - appt.bonus, 0),
    appt.downloads = appt.downloads + 1,
    dat.points = IF(appt.points_left > appt.bonus, appt.bonus, appt.points_left) * %.2f,
    dat.status = 1
WHERE
    d.id = '%s'
    AND appt.id = dat.task_id
    AND dat.device_id = d.id
    AND dat.task_id = %d
    AND dat.status != 1`
		_, _, err = db.Query(query, bonusRate, pointsPerTs.String(), bonusRate, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		rows, _, err = db.Query(`SELECT points FROM tmm.device_app_tasks WHERE device_id='%s' AND task_id=%d LIMIT 1`, db.Escape(deviceId), req.TaskId)
		if CheckErr(err, c) {
			return
		}
		if Check(len(rows) == 0, "Task not finished", c) {
			return
		}
		bonus, _ = decimal.NewFromString(rows[0].Str(0))

		query = `SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.parent_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
AND NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib WHERE ib.user_id=d.user_id AND ib.from_user_id=ic.user_id AND task_type=2 AND task_id=%d LIMIT 1)
ORDER BY d.lastping_at DESC LIMIT 1) AS t1
UNION
SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
AND NOT EXISTS (SELECT 1 FROM tmm.invite_bonus AS ib WHERE ib.user_id=d.user_id AND ib.from_user_id=ic.user_id AND task_type=2 AND task_id=%d LIMIT 1)
ORDER BY d.lastping_at DESC LIMIT 1) AS t2`
		rows, _, err = db.Query(query, user.Id, req.TaskId, user.Id, req.TaskId)
		if CheckErr(err, c) {
			return
		}
		var (
			inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * Config.InviteBonusRate))
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
		if Check(len(rows) == 0, "Task not finished", c) {
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
WHERE ic.user_id = %d AND d.user_id > 0
ORDER BY d.points DESC LIMIT 1) AS t1
UNION
SELECT id, user_id FROM
(SELECT
d.id,
d.user_id
FROM tmm.devices AS d
LEFT JOIN tmm.invite_codes AS ic ON (ic.grand_id=d.user_id)
WHERE ic.user_id = %d AND d.user_id > 0
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
		if len(deviceIds) > 0 {
			_, ret, err := db.Query(`UPDATE tmm.devices AS d, tmm.invite_bonus AS ib SET d.points = IF(d.points>ib.bonus, d.points-ib.bonus, 0) WHERE ib.user_id IN (%s) AND d.id IN (%s) AND ib.task_type=2 AND ib.task_id=%d`, strings.Join(inviterIds, ","), strings.Join(deviceIds, ","), req.TaskId)
			if CheckErr(err, c) {
				return
			}
			if ret.AffectedRows() > 0 {
				_, _, err = db.Query(`DELETE FROM tmm.invite_bonus WHERE user_id IN (%s) AND from_user_id=%d AND task_type=2 AND task_id=%d`, strings.Join(inviterIds, ","), user.Id, req.TaskId)
				if CheckErr(err, c) {
					return
				}
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
