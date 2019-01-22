package stats

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
)

func StatsByPointSupply(c *gin.Context) {
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
        SUM(IF(tmp.record_on < DATE(NOW()), 1, 0)) AS  yes_total,
        SUM(IF(tmp.record_on >= DATE(NOW()), 1, 0)) AS  today_total,
        SUM(IF(tmp.record_on < DATE(NOW()), tmp.points, 0))+IFNULL(ib.yesterday_bonus,0) AS  yes_points,
        SUM(IF(tmp.record_on >= DATE(NOW()), tmp.points, 0))+IFNULL(ib.today_bonus,0) AS  today_points,
        COUNT( DISTINCT IF(tmp.record_on >= DATE(NOW()),tmp.user_id, NULL)) AS  today_users,
        COUNT( DISTINCT IF(tmp.record_on < DATE(NOW()),tmp.user_id, NULL)) AS  yes_users,
		IF(
		(SELECT 1 FROM tmm.stats_data_logs WHERE record_on = '%s' AND (points_supply != 0 OR task_count != 0 OR task_users_count != 0) 
		) =1,1,0) AS ISEXISTS
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
    ) AS tmp,
	(
		SELECT
       		SUM(IF(ic.inserted_at > DATE(NOW()), ic.bonus, 0)) AS today_bonus,
			SUM(IF(ic.inserted_at < DATE(NOW()), ic.bonus, 0)) AS yesterday_bonus
    	FROM 
			tmm.invite_bonus AS ic
    	WHERE ic.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	) AS ib
`

	db := Service.Db
	yesterday := time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)
	rows, res, err := db.Query(query, yesterday)
	if CheckErr(err, c) {
		return
	}

	var list StatsList
	var yStats, tStats StatsData

	if len(rows) > 0 {
		row := rows[0]
		yStats.PointSupply = fmt.Sprintf("%.2f", row.Float(res.Map(`yes_points`)))
		yStats.TotalFinishTask = row.Uint64(res.Map(`yes_total`))
		yStats.TotalTaskUser = row.Uint64(res.Map(`yes_users`))
		tStats.PointSupply = fmt.Sprintf("%.2f", row.Float(res.Map(`today_points`)))
		tStats.TotalFinishTask = row.Uint64(res.Map(`today_total`))
		tStats.TotalTaskUser = row.Uint64(res.Map(`today_users`))
		if !row.Bool(res.Map(`ISEXISTS`)) {
			if _, _, err := db.Query(`
		INSERT INTO tmm.stats_data_logs(record_on,points_supply,task_count,task_users_count)
		VALUES('%s',%s,%d,%d)
		ON DUPLICATE KEY UPDATE points_supply = VALUES(points_supply), task_count = VALUES(task_count) , task_users_count = VALUES(task_users_count)`,
				yesterday, db.Escape(yStats.PointSupply), yStats.TotalFinishTask, yStats.TotalTaskUser); CheckErr(err, c) {
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
