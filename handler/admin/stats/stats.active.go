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

func StatsByWithActive(c *gin.Context) {

	key := GetCacheKey(`active`)

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
        COUNT(DISTINCT IF(record_on < DATE(NOW()), user_id, NULL)) AS yesterday_active,
        COUNT(DISTINCT IF(record_on >= DATE(NOW()), user_id, NULL)) AS today_active,
		 IF(
	   (SELECT 1 FROM tmm.stats_data_logs WHERE record_on = '%s' AND active_users != 0 
		) = 1 ,1,0) AS ISEXISTS
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
		yStats.AllActiveUsers = row.Uint64(0)
		tStats.AllActiveUsers = row.Uint64(1)

		if !row.Bool(2) {
			if _, _, err := db.Query(`INSERT INTO tmm.stats_data_logs(record_on, active_users  ) VALUES('%s',%d ) 
			ON DUPLICATE KEY UPDATE active_users = VALUES(active_users)
			`, yesterday, yStats.AllActiveUsers); CheckErr(err, c) {
				return
			}
		}
	}

	list.Today = tStats
	list.Yesterday = yStats

	if bytes, err := json.Marshal(list); err  == nil {
		conn.Do(`SET`, key, bytes, `EX`, 60)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    list,
	})
}
