package task

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
)

func ShareHandler(c *gin.Context) {
	taskId, deviceId, err := common.DecryptShareTaskLink(c.Param("encryptedTaskId"), c.Param("encryptedDeviceId"), Config)
	if CheckErr(err, c) {
		return
	}
	if Check(taskId == 0 || deviceId == "", "not found", c) {
		return
	}

	db := Service.Db
	query := `SELECT
    st.id,
    st.title,
    st.summary,
    st.link,
    st.image,
    st.max_viewers,
    st.bonus,
    st.points,
    st.points_left,
    dst.viewers
FROM tmm.share_tasks AS st
LEFT JOIN tmm.device_share_tasks AS dst ON (dst.task_id=st.id AND dst.device_id='%s')
WHERE st.id=%d
LIMIT 1`
	rows, _, err := db.Query(query, db.Escape(deviceId), taskId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	row := rows[0]
	bonus, _ := decimal.NewFromString(row.Str(6))
	points, _ := decimal.NewFromString(row.Str(7))
	pointsLeft, _ := decimal.NewFromString(row.Str(8))
	task := common.ShareTask{
		Id:         row.Uint64(0),
		Title:      row.Str(1),
		Summary:    row.Str(2),
		Link:       row.Str(3),
		Image:      row.Str(4),
		MaxViewers: row.Uint(5),
		Bonus:      bonus,
		Points:     points,
		PointsLeft: pointsLeft,
	}
	task.IsWx = task.CheckIsWechat()

	userViewers := row.Uint(9)

	var (
		cookieFound = false
		ipFound     = false
	)
	if _, err := c.Cookie(task.CookieKey()); err == nil {
		cookieFound = true
	}
	ipInfo := ClientIP(c)
	ipKey := task.IpKey(ipInfo)
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	_, err = redisConn.Do("GET", ipKey)
	if err == nil {
		ipFound = true
	}
	if !cookieFound && !ipFound && (task.PointsLeft.GreaterThanOrEqual(bonus) && task.MaxViewers > userViewers) {
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.device_share_tasks (device_id, task_id) VALUES ('%s', %d)`, db.Escape(deviceId), taskId)
		if err == nil {
			pointsPerTs, err := common.GetPointsPerTs(Service)
			if CheckErr(err, c) {
				return
			}
			query := `UPDATE tmm.devices AS d, tmm.device_share_tasks AS dst, tmm.share_tasks AS st
            SET
                d.points = d.points + IF(st.points_left > st.bonus, st.bonus, st.points_left),
                d.total_ts = d.total_ts + CEIL(IF(st.points_left > st.bonus, st.bonus, st.points_left) / %s),
                dst.points = dst.points + IF(st.points_left > st.bonus, st.bonus, st.points_left),
                dst.viewers = dst.viewers + 1,
                st.points_left = IF(st.points_left > st.bonus, st.points_left - st.bonus, 0),
                st.viewers = st.viewers + 1
            WHERE
                d.id='%s'
            AND dst.device_id = d.id
            AND dst.task_id = %d
            AND st.id = dst.task_id`
			_, _, err = db.Query(query, pointsPerTs.String(), db.Escape(deviceId), task.Id)
			if err != nil {
				log.Error(query, pointsPerTs.String(), db.Escape(deviceId), task.Id)
				log.Error(err.Error())
			}
			c.SetCookie(task.CookieKey(), "1", 60*60*24*30, "/", Config.Domain, true, true)
			_, err = redisConn.Do("SETEX", ipKey, 600, true)
			if err != nil {
				log.Error(err.Error())
			}
			query = `SELECT id, inviter_id, user_id FROM
(SELECT
d.id,
ic.parent_id AS inviter_id,
ic.user_id
FROM tmm.invite_codes AS ic
LEFT JOIN tmm.devices AS d ON (d.user_id=ic.parent_id)
LEFT JOIN tmm.devices AS d2 ON (d2.user_id=ic.user_id)
WHERE d2.id='%s'
ORDER BY d.lastping_at DESC LIMIT 1) AS t1
UNION
SELECT id, inviter_id, user_id FROM
(SELECT
d.id,
ic.grand_id AS inviter_id,
ic.user_id
FROM tmm.invite_codes AS ic
LEFT JOIN tmm.devices AS d ON (d.user_id=ic.grand_id)
LEFT JOIN tmm.devices AS d2 ON (d2.user_id=ic.user_id)
WHERE d2.id='%s'
ORDER BY d.lastping_at DESC LIMIT 1) AS t2`
			rows, _, err = db.Query(query, db.Escape(deviceId), db.Escape(deviceId))
			if err != nil {
				log.Error(err.Error())
			}
			var (
				inviterBonus = bonus.Mul(decimal.NewFromFloat(Config.InviteBonusRate))
				deviceIds    []string
				insertLogs   []string
				userId       uint64
			)
			for _, row := range rows {
				deviceIds = append(deviceIds, fmt.Sprintf("'%s'", db.Escape(row.Str(0))))
				if userId == 0 {
					userId = row.Uint64(2)
				}
				insertLogs = append(insertLogs, fmt.Sprintf("(%d, %d, %s, 1, %d)", row.Uint64(1), userId, inviterBonus.String(), task.Id))
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
		} else {
			log.Error(err.Error())
		}
	}
	c.HTML(http.StatusOK, "share.tmpl", task)
}
