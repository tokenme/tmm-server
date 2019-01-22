package stats

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
)

func StatsByInvite(c *gin.Context) {

	key := GetCacheKey(`invite`)
	conn := Service.Redis.Master.Get()
	defer conn.Close()
	bytes, err := redis.Bytes(conn.Do(`GET`, key))
	if err == nil && bytes != nil {
		var list StatsList
		if !CheckErr(json.Unmarshal(bytes, &list), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    list,
			})
		}
		return
	}

	query := `
    SELECT
        COUNT(DISTINCT IF(record_on < DATE(NOW()), user_id, NULL)) AS yesterday_New_users,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()), user_id, NULL)) AS today_New_users,
        COUNT(DISTINCT IF(record_on < DATE(NOW()) AND parent_id > 0, user_id, NULL)) AS yesterday_invite,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()) AND parent_id > 0, user_id, NULL)) AS today_invite,
        COUNT(DISTINCT IF(record_on < DATE(NOW()) AND parent_id > 0 AND (st OR att OR rl OR dl), user_id, NULL)) AS yesterday_active,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()) AND parent_id > 0 AND (st OR att OR rl OR dl), user_id, NULL)) AS today_active,
		 IF(
	   (SELECT 1 FROM tmm.stats_data_logs WHERE record_on = '%s' AND ( new_users != 0 OR invite_count != 0 OR invite_active != 0 )
		) = 1 ,1,0) AS ISEXISTS
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
        LEFT JOIN invite_codes AS ic ON (ic.user_id = u.id)
        LEFT JOIN tmm.devices AS d ON (d.user_id=u.id)
        WHERE u.created > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
    ) AS tmp
	`

	db := Service.Db
	yesterday := time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)
	rows, _, err := db.Query(query, yesterday)
	if CheckErr(err, c) {
		return
	}

	var list StatsList
	var yStats, tStats StatsData

	if len(rows) > 0 {
		row := rows[0]
		yStats.NewUsers = row.Uint64(0)
		yStats.InviteNumber = row.Uint64(2)
		yStats.Active = row.Uint64(4)

		tStats.NewUsers = row.Uint64(1)
		tStats.InviteNumber = row.Uint64(3)
		tStats.Active = row.Uint64(5)

		if !row.Bool(6) {
			now:=time.Now().Format(`2006-01-02`)
			rows,_,err=db.Query(`
SELECT
    COUNT(DISTINCT user_id) AS users,
    record_on
FROM
(
SELECT d.user_id AS user_id, DATE(dst.inserted_at) AS record_on
FROM tmm.device_share_tasks AS dst
INNER JOIN tmm.devices AS d ON (d.id=dst.device_id)
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=d.user_id AND ic.parent_id>0)
WHERE dst.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT d.user_id AS user_id, DATE(dat.inserted_at) AS record_on
FROM tmm.device_app_tasks AS dat
INNER JOIN tmm.devices AS d ON (d.id=dat.device_id)
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=d.user_id AND ic.parent_id>0)
WHERE dat.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT rl.user_id AS user_id, DATE(rl.inserted_at) AS record_on
FROM tmm.reading_logs AS rl
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=rl.user_id AND ic.parent_id>0)
WHERE rl.inserted_at BETWEEN '%s' AND '%s'
UNION ALL
SELECT dbl.user_id AS user_id, DATE(dbl.updated_on) AS record_on
FROM tmm.daily_bonus_logs AS dbl
INNER JOIN tmm.invite_codes AS ic ON (ic.user_id=dbl.user_id AND ic.parent_id>0)
WHERE dbl.updated_on = '%s' 
) AS tmp
GROUP BY record_on`,yesterday,now,yesterday,now,yesterday,now,yesterday)
			active:=rows[0].Int(0)
			if _, _, err := db.Query(`INSERT INTO tmm.stats_data_logs(record_on, new_users , invite_count , invite_active ) VALUES('%s',%d,%d,%d)
			ON DUPLICATE KEY UPDATE new_users = VALUES(new_users), invite_count = VALUES(invite_count) , invite_active = VALUES(invite_active)
			`, yesterday, yStats.NewUsers, yStats.InviteNumber, active); CheckErr(err, c) {
				return
			}
		}
	}

	list.Today = tStats
	list.Yesterday = yStats

	if bytes, err := json.Marshal(list); err == nil {
		conn.Do(`SET`, key, bytes, `EX`, 60)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    list,
	})
}
