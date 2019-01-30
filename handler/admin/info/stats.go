package info

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

const statsKey = `info-stats-stats`

func StatsHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, statsKey))
	if context != nil && err == nil {
		var data interface{}
		if !CheckErr(json.Unmarshal(context, &data), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    data,
			})
			return
		}
		return
	}
	query := `
SELECT
    points.today_number AS point_today_number,
    points.yesterday_number AS point_yesterday_number,
    tmm.today_number AS tmm_today_number,
    tmm.yesterday_number AS tmm_yesterday_number,
    IFNULL(points.today_cny,0)+IFNULL(tmm.today_cny,0) AS today_cny,
    IFNULL(points.yesterday_cny,0)+IFNULL(tmm.yesterday_cny,0) AS yesterday_cny,
    IFNULL(device.today_points,0)+IFNULL(inv_bonus.today_bonus,0) AS today_point_supply,
    IFNULL(device.yes_points,0)+IFNULL(inv_bonus.yesterday_bonus,0) AS yesterday_point_supply,
    device.today_users AS today_users,
    device.yes_users AS yesterday_users,
    device.today_total AS today_task_number,
    device.yes_total AS yesterday_task_number,
    IFNULL(exchanges.today_tmm,0) AS today_tmm_supply,
    IFNULL(exchanges.yesterday_tmm,0) AS yesterday_tmm_supply,
    users.today_invite AS today_invite,
    users.yesterday_invite AS yesterday_invite,
    users.today_active AS today_active,
    users.yesterday_active AS yesterday_active,
    active.today_active AS today_all_active,
    active.yesterday_active AS yesterday_all_active,
    users.today_New_users AS today_New_users,
    users.yesterday_New_users AS yesterday_New_users
FROM
(
    SELECT
        COUNT(IF(inserted_at > DATE(NOW()),1,NULL)) AS today_number,
        COUNT(IF(inserted_at < DATE(NOW()),1,NULL)) AS yesterday_number,
        SUM(IF(inserted_at > DATE(NOW()),cny,0)) AS today_cny,
        SUM(IF(inserted_at < DATE(NOW()),cny,0)) AS yesterday_cny
    FROM tmm.point_withdraws 
    WHERE inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND verified!=-1 AND  trade_num !=  ""
) AS points,
(
    SELECT
        COUNT(IF(inserted_at > DATE(NOW()),1,NULL)) AS today_number,
        COUNT(IF(inserted_at < DATE(NOW()),1,NULL)) AS yesterday_number,
        SUM(IF(inserted_at > DATE(NOW()),cny,0)) AS today_cny,
        SUM(IF(inserted_at < DATE(NOW()),cny,0)) AS yesterday_cny
    FROM tmm.withdraw_txs
    WHERE inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND withdraw_status=1 AND verified!=-1 
) AS tmm,
(
    SELECT
        SUM(IF(tmp.record_on < DATE(NOW()), 1, 0)) AS  yes_total,
        SUM(IF(tmp.record_on >= DATE(NOW()), 1, 0)) AS  today_total,
        SUM(IF(tmp.record_on < DATE(NOW()), tmp.points, 0)) AS  yes_points,
        SUM(IF(tmp.record_on >= DATE(NOW()), tmp.points, 0)) AS  today_points,
        COUNT( DISTINCT IF(tmp.record_on >= DATE(NOW()),tmp.user_id, NULL)) AS  today_users,
        COUNT( DISTINCT IF(tmp.record_on < DATE(NOW()),tmp.user_id, NULL)) AS  yes_users
    FROM
    (
        SELECT DATE(dst.inserted_at) AS record_on, d.user_id AS user_id, dst.points AS points
        FROM tmm.device_share_tasks  AS dst
        INNER JOIN tmm.devices AS d ON  (d.id = dst.device_id)
        WHERE dst.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT DATE(app.inserted_at) AS record_on, d.user_id AS user_id, IF(app.status = 1, app.points, 0) AS points
        FROM tmm.device_app_tasks AS app
        INNER JOIN tmm.devices AS d ON  (d.id = app.device_id)
        WHERE app.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT DATE(rl.inserted_at) AS record_on, rl.user_id AS user_id, rl.point AS points
        FROM tmm.reading_logs AS rl
        WHERE rl.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
		UNION ALL 
		SELECT DATE(dgt.inserted_at) AS record_on, d.user_id AS user_id , dgt.points
		FROM tmm.device_general_tasks AS dgt 
		INNER JOIN tmm.devices AS d ON (d.id = dgt.device_id )
		WHERE dgt.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND status = 1
    ) AS tmp
) AS device,
(
    SELECT
       SUM(IF(inserted_at < DATE(NOW()),tmm,0)) AS yesterday_tmm,
       SUM(IF(inserted_at > DATE(NOW()),tmm,0)) AS today_tmm
    FROM tmm.exchange_records
    WHERE direction=1 AND inserted_at>DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND status=1
)  AS exchanges,
(
    SELECT
       SUM(IF(ic.inserted_at > DATE(NOW()), ic.bonus, 0)) AS today_bonus,
       SUM(IF(ic.inserted_at < DATE(NOW()), ic.bonus, 0)) AS yesterday_bonus
    FROM tmm.invite_bonus AS ic
    WHERE ic.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
) AS inv_bonus ,
(
    SELECT
        COUNT(DISTINCT IF(record_on < DATE(NOW()), user_id, NULL)) AS yesterday_New_users,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()), user_id, NULL)) AS today_New_users,
        COUNT(DISTINCT IF(record_on < DATE(NOW()) AND parent_id > 0, user_id, NULL)) AS yesterday_invite,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()) AND parent_id > 0, user_id, NULL)) AS today_invite,
        COUNT(DISTINCT IF(record_on < DATE(NOW()) AND parent_id > 0 AND (st OR att OR rl OR dl), user_id, NULL)) AS yesterday_active,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()) AND parent_id > 0 AND (st OR att OR rl OR dl), user_id, NULL)) AS today_active
    FROM
    (
        SELECT
            u.id AS user_id,
            ic.parent_id AS parent_id,
            DATE(u.created) AS record_on,
            EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id = d.id LIMIT 1) AS st,
            EXISTS (SELECT 1 FROM tmm.device_app_tasks AS app WHERE app.device_id = d.id LIMIT 1) AS att,
            EXISTS (SELECT 1 FROM tmm.reading_logs AS reading WHERE reading.user_id = u.id LIMIT 1) AS rl,
            EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS daily WHERE daily.user_id = u.id LIMIT 1) AS dl
        FROM ucoin.users AS u
        INNER JOIN invite_codes AS ic ON (ic.user_id = u.id)
        INNER JOIN tmm.devices AS d ON (d.user_id=u.id)
        WHERE u.created > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
    ) AS tmp
) AS users,
(
    SELECT
        COUNT(DISTINCT IF(record_on < DATE(NOW()), user_id, NULL)) AS yesterday_active,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()), user_id, NULL)) AS today_active
    FROM
    (
        SELECT d.user_id AS user_id, DATE(dst.inserted_at) AS record_on
        FROM tmm.device_share_tasks AS dst
        INNER JOIN tmm.devices AS d ON (d.id=dst.device_id)
        WHERE dst.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT d.user_id AS user_id, DATE(dat.inserted_at) AS record_on
        FROM tmm.device_app_tasks AS dat
        INNER JOIN tmm.devices AS d ON (d.id=dat.device_id)
        WHERE dat.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT rl.user_id AS user_id, DATE(rl.inserted_at) AS record_on
        FROM tmm.reading_logs AS rl
        WHERE rl.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT rl.user_id AS user_id, DATE(rl.updated_at) AS record_on
        FROM tmm.reading_logs AS rl
        WHERE rl.updated_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
        UNION ALL
        SELECT dbl.user_id AS user_id, DATE(dbl.updated_on) AS record_on
        FROM tmm.daily_bonus_logs AS dbl
        WHERE dbl.updated_on > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
    ) AS tmp
) AS active
`
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `没有找到数据`, c) {
		return
	}
	row := rows[0]
	var yesterdayStats, todayStats StatsData
	yesterdayStats.PointExchangeNumber = row.Int(res.Map(`point_yesterday_number`))
	yesterdayStats.UcoinExchangeNumber = row.Int(res.Map(`tmm_yesterday_number`))
	yesterdayStats.TotalTaskUser = row.Int(res.Map(`yesterday_users`))
	yesterdayStats.TotalFinishTask = row.Int(res.Map(`yesterday_task_number`))
	yesterdayStats.InviteNumber = row.Int(res.Map(`yesterday_invite`))
	yesterdayStats.Active = row.Int(res.Map(`yesterday_active`))
	yesterdayStats.Cash = fmt.Sprintf("%.2f", row.Float(res.Map(`yesterday_cny`)))
	yesterdayStats.PointSupply = fmt.Sprintf("%.1f", row.Float(res.Map(`yesterday_point_supply`)))
	yesterdayStats.UcSupply = fmt.Sprintf("%.1f", row.Float(res.Map(`yesterday_tmm_supply`)))
	yesterdayStats.AllActiveUsers = row.Int(res.Map(`yesterday_all_active`))
	yesterdayStats.NewUsers = row.Int(res.Map(`yesterday_New_users`))

	todayStats.PointExchangeNumber = row.Int(res.Map(`point_today_number`))
	todayStats.UcoinExchangeNumber = row.Int(res.Map(`tmm_today_number`))
	todayStats.TotalTaskUser = row.Int(res.Map(`today_users`))
	todayStats.TotalFinishTask = row.Int(res.Map(`today_task_number`))
	todayStats.InviteNumber = row.Int(res.Map(`today_invite`))
	todayStats.Active = row.Int(res.Map(`today_active`))
	todayStats.Cash = fmt.Sprintf("%.2f", row.Float(res.Map(`today_cny`)))
	todayStats.PointSupply = fmt.Sprintf("%.1f", row.Float(res.Map(`today_point_supply`)))
	todayStats.UcSupply = fmt.Sprintf("%.1f", row.Float(res.Map(`today_tmm_supply`)))
	todayStats.AllActiveUsers = row.Int(res.Map(`today_all_active`))
	todayStats.NewUsers = row.Int(res.Map(`today_New_users`))
	var statsList StatsList
	statsList.Yesterday = yesterdayStats
	statsList.Today = todayStats
	bytes, err := json.Marshal(&statsList)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do(`SET`, statsKey, bytes, `EX`, 60)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    statsList,
	})
}
