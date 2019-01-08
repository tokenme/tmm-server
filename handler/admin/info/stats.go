package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"fmt"
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
FROM 
		tmm.point_withdraws
WHERE 
		inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
) AS points,
(
SELECT 
		COUNT(IF(inserted_at > DATE(NOW()),1,NULL)) AS today_number,
		COUNT(IF(inserted_at < DATE(NOW()),1,NULL)) AS yesterday_number,
		SUM(IF(inserted_at > DATE(NOW()),cny,0)) AS today_cny,
		SUM(IF(inserted_at < DATE(NOW()),cny,0)) AS yesterday_cny
FROM 
		tmm.withdraw_txs 
WHERE 
		inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)   AND tx_status = 1
) AS tmm,
( 
SELECT 
	SUM(IF(tmp.date < DATE(NOW()),tmp.times,0)) AS  yes_total,
	SUM(IF(tmp.date >= DATE(NOW()),tmp.times,0)) AS  today_total,
	SUM(IF(tmp.date < DATE(NOW()),tmp.points,0)) AS  yes_points,
	SUM(IF(tmp.date >= DATE(NOW()),tmp.points,0)) AS  today_points,
	COUNT( DISTINCT IF(tmp.date >= DATE(NOW()),tmp.user_id,NULL)) AS  today_users,
	COUNT( DISTINCT IF(tmp.date < DATE(NOW()),tmp.user_id,NULL)) AS  yes_users
FROM (
SELECT 		
			DATE(sha.inserted_at) AS date,
			dev.user_id AS user_id ,
			COUNT(1) AS times,
			SUM(sha.points) AS points
		FROM 
			tmm.device_share_tasks  AS sha 
		INNER JOIN 
			tmm.devices AS dev ON  (dev.id = sha.device_id)
		WHERE 
			sha.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
		GROUP BY 
			date,dev.user_id 
	UNION ALL 
		SELECT 
			DATE(app.inserted_at) AS date ,
			dev.user_id AS user_id,
			COUNT(1) AS times,
			SUM(IF(app.status = 1,app.points,0)) AS points
		FROM 
			tmm.device_app_tasks  AS app 
		INNER JOIN 
			tmm.devices AS dev ON  (dev.id = app.device_id)
		WHERE 
			app.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
		GROUP BY date,dev.user_id 
	UNION ALL 
		SELECT 
			DATE(inserted_at) AS date,
			user_id AS user_id ,
			COUNT(1) AS times,
			SUM(point)
		FROM 
			reading_logs 
		WHERE 
			inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
		GROUP BY 
			date,user_id 
			) AS tmp 
) AS device,
(
SELECT 
	   SUM(IF(inserted_at < DATE(NOW()),tmm,0)) AS yesterday_tmm, 
	   SUM(IF(inserted_at > DATE(NOW()),tmm,0)) AS today_tmm  
FROM 
	   tmm.exchange_records
WHERE 
	   direction = 1 AND inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND status = 1
)  AS exchanges,
(
SELECT 
	SUM(IF(inv.inserted_at > DATE(NOW()),inv.bonus,0)) AS today_bonus,
	SUM(IF(inv.inserted_at < DATE(NOW()),inv.bonus,0)) AS yesterday_bonus
FROM
	tmm.invite_bonus AS inv
WHERE
	inv.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) 
) AS inv_bonus ,
(
SELECT 
	COUNT(IF(u.created < DATE(NOW()),0,NULL)) AS yesterday_New_users,
	COUNT(IF(u.created > DATE(NOW()),0,NULL)) AS today_New_users,
	COUNT(IF(u.created < DATE(NOW()) AND inv.parent_id > 0,u.id,NULL)) AS yesterday_invite,
	COUNT(IF(u.created > DATE(NOW()) AND inv.parent_id > 0,0,NULL)) AS today_invite,  
	COUNT(distinct IF(u.created < DATE(NOW()) AND inv.parent_id > 0 AND
		EXISTS(
		SELECT
		1
		FROM tmm.devices AS dev
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id)
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id )
		LEFT JOIN tmm.reading_logs AS reading ON (reading.user_id = dev.user_id)
		LEFT JOIN tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id)
		WHERE dev.user_id = u.id AND ( 
		sha.task_id > 0 OR 
		app.task_id > 0 OR 
		reading.point > 0 OR
		daily.user_id > 0)
		LIMIT 1
		),u.id,NULL))  AS yesterday_active,
	  COUNT(distinct IF(u.created > DATE(NOW()) AND inv.parent_id > 0 AND
		EXISTS(
		SELECT
		1
		FROM tmm.devices AS dev
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id)
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id )
		LEFT JOIN tmm.reading_logs AS reading ON (reading.user_id = dev.user_id)
		LEFT JOIN tmm.daily_bonus_logs AS daily ON (daily.user_id = dev.user_id)
		WHERE dev.user_id = u.id AND ( 
		sha.task_id > 0 OR 
		app.task_id > 0 OR 
		reading.point > 0 OR
		daily.user_id > 0)
		LIMIT 1
		),u.id,NULL))  AS today_active
FROM 
	ucoin.users  AS u 
LEFT JOIN invite_codes AS inv ON inv.user_id = u.id
WHERE
	u.created > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
) AS users,
(
SELECT 
	COUNT(DISTINCT IF(tmp.date < DATE(NOW()),tmp.user_id,NULL)) AS yesterday_active,
	COUNT(DISTINCT IF(tmp.date >= DATE(NOW()),tmp.user_id,NULL)) AS today_active
FROM (
 SELECT  
	dev.user_id AS user_id ,
	DATE(sha.inserted_at) AS date 
	FROM 
	  tmm.device_share_tasks  AS sha 
	INNER JOIN tmm.devices AS dev ON dev.id = sha.device_id
	WHERE sha.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	GROUP BY user_id,date
	UNION ALL
	SELECT 
		dev.user_id AS user_id ,
		DATE(app.inserted_at) AS date 
	FROM 
		tmm.device_app_tasks  AS app 
	INNER JOIN tmm.devices AS dev ON dev.id = app.device_id
	WHERE app.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	GROUP BY user_id,date
	UNION ALL 
	SELECT 
		user_id  AS user_id ,
		DATE(reading.inserted_at) AS date 
	FROM 
		tmm.reading_logs AS reading 
	WHERE reading.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) 
	GROUP BY user_id,date
	UNION ALL 
	SELECT
		user_id   AS user_id ,
		DATE(reading.updated_at) AS date 
	FROM 
		tmm.reading_logs AS reading 
	WHERE reading.updated_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	UNION ALL 
	SELECT 
		user_id  AS user_id ,
		DATE(updated_on) AS date 
	FROM 
		daily_bonus_logs
	WHERE updated_on >= DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) 
	GROUP BY user_id,date
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
