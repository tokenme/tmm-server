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

func StatsByWithDrawCash(c *gin.Context) {

	key := GetCacheKey(`drawcash`)

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
	points.today_count AS p_today_count,
	points.today_users AS p_today_users,
	points.today_cny   AS p_today_cny,
	points.yesterday_count AS p_yesterday_count,
	points.yesterday_users AS p_yesterday_users,
	points.yesterday_cny AS p_yesterday_cny,
	uc.today_count 	    AS u_today_count,
	uc.today_users 	    AS u_today_users,
	uc.today_cny        AS u_today_cny,
	uc.yesterday_count  AS u_yesterday_count,
	uc.yesterday_users  AS u_yesterday_users,
	uc.yesterday_cny    AS u_yesterday_cny,
	unmet.yesterday_cny AS um_yesterday_cny,
	unmet.today_cny     AS um_today_cny,
	IF((SELECT 1 FROM tmm.stats_data_logs WHERE record_on = '%s' 
	AND( 
	points_cash != 0 
	OR points_users != 0 
	OR points_count != 0 
	OR uc_cash  != 0 
	OR uc_users != 0 
	OR uc_count != 0
	OR unmet_cash != 0) 
	) = 1,1,0) AS IsEXISTS
FROM
(
    SELECT
        COUNT(IF(inserted_at > DATE(NOW()),1,NULL)) AS today_count,
        COUNT(DISTINCT IF(inserted_at > DATE(NOW()),user_id,NULL)) AS today_users,
 		SUM(IF(inserted_at > DATE(NOW()),cny,0)) 	AS today_cny,
        COUNT(IF(inserted_at < DATE(NOW()),1,NULL)) AS yesterday_count,
		COUNT(DISTINCT IF(inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_users,
        SUM(IF(inserted_at < DATE(NOW()),cny,0)) AS yesterday_cny
    FROM 
		tmm.point_withdraws
    WHERE 
		inserted_at > '%s' AND verified!=-1
) AS points,
(
    SELECT
		COUNT(IF(inserted_at > DATE(NOW()),1,NULL)) AS today_count,
        COUNT(DISTINCT IF(inserted_at > DATE(NOW()),user_id,NULL)) AS today_users,
 		SUM(IF(inserted_at > DATE(NOW()),cny,0)) 	AS today_cny,
        COUNT(IF(inserted_at < DATE(NOW()),1,NULL)) AS yesterday_count,
		COUNT(DISTINCT IF(inserted_at < DATE(NOW()),user_id,NULL)) AS yesterday_users,
        SUM(IF(inserted_at < DATE(NOW()),cny,0)) AS yesterday_cny
    FROM tmm.withdraw_txs
    WHERE inserted_at > '%s' AND tx_status=1 AND verified!=-1
) AS uc,
(
	SELECT 
		SUM(IF(inserted_at > DATE(NOW()),cny,0)) AS today_cny,
		SUM(IF(inserted_at < DATE(NOW()),cny,0) ) AS yesterday_cny
	FROM tmm.withdraw_logs
	WHERE inserted_at > '%s'
) AS unmet
`

	db := Service.Db
	yesterday := time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)
	rows, res, err := db.Query(query, yesterday, yesterday, yesterday, yesterday)
	if CheckErr(err, c) {
		return
	}

	var list StatsList
	var yStats, tStats StatsData

	if len(rows) > 0 {
		row := rows[0]
		yStats.UserCountByUcDrawCash = row.Uint64(res.Map(`u_yesterday_users`))
		yStats.UserCountByPointDrawCash = row.Uint64(res.Map(`p_yesterday_users`))
		yStats.Cash = fmt.Sprintf("%.2f", row.Float(res.Map(`u_yesterday_cny`))+row.Float(res.Map(`p_yesterday_cny`)))
		tStats.UserCountByUcDrawCash = row.Uint64(res.Map(`u_today_users`))
		tStats.UserCountByPointDrawCash = row.Uint64(res.Map(`p_today_users`))
		tStats.Cash = fmt.Sprintf("%.2f", row.Float(res.Map(`u_today_cny`))+row.Float(res.Map(`p_today_cny`)))

		//INSERT yesterday stats
		if !row.Bool(res.Map(`IsEXISTS`)) {

			_,_,err=db.Query(`INSERT INTO tmm.stats_data_logs(record_on , points_cash , points_users , points_count , uc_cash , uc_users , uc_count , unmet_cash)
			VALUES('%s',%s,%d,%d,%s,%d,%d,%s)
			ON DUPLICATE KEY UPDATE points_cash = VALUES(points_cash),points_users = VALUES(points_users),points_count = VALUES(points_count), 
			uc_cash = VALUES(uc_cash), uc_users = VALUES(uc_users) , uc_count = VALUES(uc_count), unmet_cash =VALUES(unmet_cash)`,
				yesterday, row.Str(res.Map(`p_yesterday_cny`)), row.Int(res.Map(`p_yesterday_users`)), row.Int(res.Map(`p_yesterday_count`)),
				row.Str(res.Map(`u_yesterday_cny`)), row.Int(res.Map(`u_yesterday_users`)), row.Int(res.Map(`u_yesterday_count`)), row.Str(res.Map(`um_yesterday_cny`)))
			fmt.Println(err)
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
