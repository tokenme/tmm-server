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

func StatsByExchange(c *gin.Context) {
	key := GetCacheKey(`exchange`)

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
       SUM(IF(direction =  1 AND inserted_at < DATE(NOW()),tmm,0)) AS yesterday_tmm,
	   COUNT(DISTINCT IF(direction =  1 AND inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_tmm_users,
       SUM(IF(direction =  1 AND inserted_at > DATE(NOW()),tmm,0)) AS today_tmm,
       SUM(IF(direction = -1 AND inserted_at < DATE(NOW()),points,0)) AS yesterday_point,
       COUNT(DISTINCT IF(direction = -1 AND inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_point_users,
       SUM(IF(direction = -1 AND inserted_at > DATE(NOW()),points,0)) AS today_point,
	   IF(
	   (SELECT 1 FROM tmm.stats_data_logs WHERE record_on = '%s' AND ( points_ec_uc != 0 OR uc_ec_points != 0 )
		) = 1 ,1,0) AS ISEXISTS
    FROM tmm.exchange_records
    WHERE inserted_at > '%s' AND status=1
	`

	db := Service.Db
	yesterday := time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)

	rows, _, err := db.Query(query, yesterday, yesterday)
	if CheckErr(err, c) {
		return
	}

	var list StatsList
	var yStats, tStats StatsData

	if len(rows) > 0 {
		row := rows[0]
		yStats.UcSupply = fmt.Sprintf("%.2f", row.Float(0))
		tStats.UcSupply = fmt.Sprintf("%.2f", row.Float(2))

		if !row.Bool(6) {
			if _, _, err := db.Query(`INSERT INTO tmm.stats_data_logs
			(record_on, points_ec_uc ,points_ec_uc_users, uc_ec_points , uc_ec_points_users ) VALUES('%s',%s,%s,%s,%s)
			ON DUPLICATE KEY UPDATE points_ec_uc = VALUES(points_ec_uc), uc_ec_points = VALUES(uc_ec_points) , 
			points_ec_uc_users = VALUES(points_ec_uc_users),  uc_ec_points_users = VALUES(uc_ec_points_users)
			`, yesterday, row.Str(0), row.Str(1), row.Str(3), row.Str(4)); CheckErr(err, c) {
				return
			}
		}
	}

	/*
	   SUM(IF(direction =  1 AND inserted_at < DATE(NOW()),tmm,0)) AS yesterday_tmm,
	   COUNT(IF(direction =  1 AND inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_tmm_users,
       SUM(IF(direction =  1 AND inserted_at > DATE(NOW()),tmm,0)) AS today_tmm,
       SUM(IF(direction = -1 AND inserted_at < DATE(NOW()),points,0)) AS yesterday_point,
       COUNT(IF(direction = -1 AND inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_point_users,
       SUM(IF(direction = -1 AND inserted_at > DATE(NOW()),points,0)) AS today_point,
		ISEXISTS
	*/
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
