package task

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
	"time"
)

func ShareImpHandler(c *gin.Context) {
	taskId, deviceId, err := common.DecryptShareTaskLink(c.Param("encryptedTaskId"), c.Param("encryptedDeviceId"), Config)
	if CheckErr(err, c) {
		return
	}
	if Check(taskId == 0 || deviceId == "", "not found", c) {
		return
	}
	isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")
	if !isWx {
		log.Info("DeviceId: %s, UA: %s", deviceId, c.Request.UserAgent())
	}
	if !isWx {
		c.Redirect(http.StatusFound, "https://ucoin.tianxi100.com/_.gif")
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
LEFT JOIN tmm.devices AS d ON (d.id=dst.device_id)
WHERE st.id=%d
  AND st.online_status=1
LIMIT 1`
	rows, _, err := db.Query(query, db.Escape(deviceId), taskId)
	if err != nil {
		c.Redirect(http.StatusFound, "https://ucoin.tianxi100.com/_.gif")
		return
	}
	if len(rows) == 0 {
		log.Error("Share Imp Not found")
		c.Redirect(http.StatusFound, "https://ucoin.tianxi100.com/_.gif")
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
	if strings.HasPrefix(task.Link, "https://tmm.tokenmama.io/article/show") {
		task.Link = strings.Replace(task.Link, "https://tmm.tokenmama.io/article/show", "https://static.tianxi100.com/article/show", -1)
	}

	trackId := c.Query("track_id")
	var openid string
	if len(trackId) > 0 {
		cryptOpenid, err := common.DecodeCryptOpenid([]byte(Config.YktApiSecret), trackId)
		if err != nil {
			log.Error("Decrypt track id error")
		} else {
			openid = cryptOpenid.Openid
		}
	}

	if openid == "" {
		c.Redirect(http.StatusFound, "https://ucoin.tianxi100.com/_.gif")
		return
	}

	userViewers := row.Uint(9)
	var (
		cookieFound = false
		ipFound     = false
		openidFound = false
	)
	if _, err := c.Cookie(task.CookieKey()); err == nil {
		log.Warn("Share Imp Cookie Found")
		cookieFound = true
	}

	ipInfo := ClientIP(c)
	ipKey := task.IpKey(ipInfo)
	openidKey := task.OpenidKey(openid)
	{
		redisConn := Service.Redis.Master.Get()
		defer redisConn.Close()
		cacheRet, err := redis.Strings(redisConn.Do("MGET", ipKey, openidKey))
		if err != nil {
			log.Error(err.Error())
		} else {
			if cacheRet[0] != "" {
				log.Warn("Share Imp IP Found: %s, time: %s", ipKey, cacheRet[0])
				ipFound = true
			}
			if cacheRet[1] != "" {
				log.Warn("Share Imp Openid Found: %s, time: %s", openidKey, cacheRet[1])
				openidFound = true
			}
		}

		if !cookieFound && !ipFound && !openidFound && (task.PointsLeft.GreaterThanOrEqual(bonus) && task.MaxViewers > userViewers) {
			log.Info("Share Imp Task: %d, Device: %s", taskId, deviceId)
			_, _, err := db.Query(`INSERT IGNORE INTO tmm.device_share_tasks (device_id, task_id) VALUES ('%s', %d)`, db.Escape(deviceId), taskId)
			if err == nil {
				pointsPerTs, err := common.GetPointsPerTs(Service)
				if err != nil {
					return
				}
				var bonusRate float64 = 1
				rows, _, err := db.Query(`SELECT ul.task_bonus_rate FROM tmm.user_settings AS us INNER JOIN tmm.user_levels AS ul ON (ul.id=us.level) INNER JOIN tmm.devices AS d ON (d.user_id=us.user_id) WHERE d.id='%s' LIMIT 1`, db.Escape(deviceId))
				if err != nil {
					log.Error(err.Error())
				} else if len(rows) > 0 {
					bonusRate = rows[0].ForceFloat(0) / 100
				}

				query := `UPDATE tmm.devices AS d, tmm.device_share_tasks AS dst, tmm.share_tasks AS st
            SET
                d.points = d.points + IF(st.points_left > st.bonus, st.bonus, st.points_left) * %.2f,
                d.total_ts = d.total_ts + CEIL(IF(st.points_left > st.bonus, st.bonus, st.points_left) / %s),
                dst.points = dst.points + IF(st.points_left > st.bonus, st.bonus, st.points_left) * %.2f,
                dst.viewers = dst.viewers + 1,
                st.points_left = IF(st.points_left > st.bonus, st.points_left - st.bonus, 0),
                st.viewers = st.viewers + 1
            WHERE
                d.id='%s'
            AND dst.device_id = d.id
            AND dst.task_id = %d
            AND st.id = dst.task_id`
				_, _, err = db.Query(query, bonusRate, pointsPerTs.String(), bonusRate, db.Escape(deviceId), task.Id)
				if err != nil {
					log.Error(err.Error())
				}
				c.SetCookie(task.CookieKey(), "1", 60*60*24*30, "/", Config.Domain, true, true)
				_, err = redisConn.Do("SETEX", ipKey, 600, time.Now().Format("2006-01-02 15:04:05"))
				if err != nil {
					log.Error(err.Error())
				}
				if openid != "" {
					_, err = redisConn.Do("SETEX", openidKey, 60*60*24*30, time.Now().Format("2006-01-02 15:04:05"))
					if err != nil {
						log.Error(err.Error())
					}
				}
				query = `SELECT t.id, t.inviter_id, t.user_id FROM
(SELECT id, inviter_id, user_id FROM
    (SELECT
    d.id,
    ic.parent_id AS inviter_id,
    ic.user_id
    FROM tmm.invite_codes AS ic
    LEFT JOIN tmm.devices AS d ON (d.user_id=ic.parent_id)
    LEFT JOIN tmm.devices AS d2 ON (d2.user_id=ic.user_id)
    WHERE d2.id='%s' AND ic.parent_id > 0
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
    WHERE d2.id='%s' AND ic.grand_id > 0
    ORDER BY d.lastping_at DESC LIMIT 1) AS t2
) AS t
LEFT JOIN tmm.wx AS wx ON (wx.user_id=t.inviter_id)
LEFT JOIN tmm.wx AS wx2 ON (wx2.user_id=t.user_id)
WHERE ISNULL(wx.open_id) OR ISNULL(wx2.open_id) OR wx.open_id!=wx2.open_id`
				rows, _, err = db.Query(query, db.Escape(deviceId), db.Escape(deviceId))
				if err != nil {
					log.Error(err.Error())
				}
				var (
					inviterBonus = bonus.Mul(decimal.NewFromFloat(bonusRate * Config.InviteBonusRate))
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
					_, ret, err := db.Query(`UPDATE tmm.devices SET points = points + %s, total_ts = total_ts + %d WHERE id IN (%s)`, inviterBonus.String(), inviterBonus.Div(pointsPerTs).IntPart(), strings.Join(deviceIds, ","))
					if err != nil {
						return
					}
					if ret.AffectedRows() > 0 {
						_, _, err = db.Query(`INSERT INTO tmm.invite_bonus (user_id, from_user_id, bonus, task_type, task_id) VALUES %s`, strings.Join(insertLogs, ","))
						if err != nil {
							return
						}
					}
				}
				trackSource := common.TrackFromUnknown
				if len(openid) > 0 {
					trackSource = common.TrackFromWechat
				}
				_, _, err = db.Query(`INSERT INTO tmm.share_task_stats (task_id, record_on, source) VALUES (%d, NOW(), %d) ON DUPLICATE KEY UPDATE uv=uv+1`, task.Id, trackSource)
				if err != nil {
					log.Error(err.Error())
				}
			} else {
				log.Error(err.Error())
			}
		}
	}
	c.Redirect(http.StatusFound, "https://ucoin.tianxi100.com/_.gif")
}
