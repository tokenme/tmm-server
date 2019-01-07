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
		inv.today_invite AS today_invite,
		inv.yesterday_invite AS yesterday_invite,
		inv.today_active AS today_active,
		inv.yesterday_active AS yesterday_active,
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
		COUNT(IF(tmp.today_total > 0,0,NULL)) AS today_users,
		COUNT(IF(tmp.yes_total > 0,0,NULL)) AS yes_users,
		SUM(tmp.today_points) AS today_points,
		SUM(tmp.yes_points) AS yes_points,
		SUM(tmp.yes_total) AS yes_total,
		SUM(tmp.today_total) AS today_total
FROM(
	SELECT 
	dev.user_id AS user_id,
	SUM(tmp.today_total) + IFNULL(reading.today_times,0) AS today_total,
	SUM(tmp.yes_total) + IFNULL(reading.yesterday_times,0) AS yes_total,
	SUM(tmp.yes_points) + IFNULL(reading.yesterday_point,0) AS yes_points,
	SUM(tmp.today_points) + IFNULL(reading.today_point,0) AS today_points
FROM
	tmm.devices AS dev
LEFT JOIN(
	SELECT 
		sha.device_id,
		SUM(IF(sha.inserted_at > DATE(NOW()),sha.points,0)) AS today_points,
		SUM(IF(sha.inserted_at < DATE(NOW()),sha.points,0)) AS yes_points,
		COUNT(IF(sha.inserted_at < DATE(NOW()),0,NULL)) AS yes_total,
		COUNT(IF(sha.inserted_at > DATE(NOW()),0,NULL)) AS today_total
	FROM 
		tmm.device_share_tasks  AS sha
	WHERE 
		sha.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	GROUP BY sha.device_id 
	UNION ALL
	SELECT 
		app.device_id,
		SUM(IF(app.inserted_at > DATE(NOW()),app.points,0)) AS today_points,
		SUM(IF(app.inserted_at < DATE(NOW()),app.points,0)) AS yes_points,
		COUNT(IF(app.inserted_at < DATE(NOW()),0,NULL)) AS yes_total,
		COUNT(IF(app.inserted_at > DATE(NOW()),0,NULL)) AS today_total
	FROM
		tmm.device_app_tasks AS app 
	WHERE 
		app.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND status = 1
	GROUP BY 
		app.device_id
	) AS tmp ON (tmp.device_id = dev.id)
	LEFT JOIN (
	SELECT
		user_id,
		SUM(IF(inserted_at > DATE(NOW()),point,0)) AS today_point,
		SUM(IF(inserted_at < DATE(NOW()),point,0)) AS yesterday_point,
		COUNT(IF(inserted_at > DATE(NOW()),0,NULL)) AS today_times,
		COUNT(IF(inserted_at < DATE(NOW()),0,NULL)) AS yesterday_times
	FROM
		tmm.reading_logs
	WHERE
		inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
	GROUP BY 
		user_id
	) AS reading ON (reading.user_id = dev.user_id)
	WHERE 1 = 1 AND ( IFNULL(tmp.today_total,0) > 0 OR IFNULL(tmp.yes_points,0) > 0 OR IFNULL(tmp.yes_total,0) > 0 OR IFNULL(tmp.today_points,0) > 0 )
	GROUP BY 
		dev.user_id

	) AS tmp 
) AS device,
(
SELECT 
	   SUM(IF(inserted_at < DATE(NOW()),tmm,0)) AS yesterday_tmm, 
	   SUM(IF(inserted_at > DATE(NOW()),tmm,0)) AS today_tmm  
FROM 
	   tmm.exchange_records
WHERE 
	   direction = 1 AND inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
)  AS exchanges,
(
SELECT
	  COUNT(distinct IF(inv.inserted_at < DATE(NOW()) AND
		EXISTS(
		SELECT
		1
		FROM tmm.devices AS dev
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id)
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id AND app.status = 1)
		LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id  )
		WHERE dev.user_id = inv.from_user_id AND ( 
		sha.inserted_at > inv.inserted_at OR 
		app.inserted_at > inv.inserted_at OR    
		reading.inserted_at > inv.inserted_at)
		LIMIT 1
		),inv.from_user_id,NULL))  AS yesterday_active,
	  COUNT(distinct IF(inv.inserted_at > DATE(NOW()) AND
		EXISTS(
		SELECT
		1
		FROM tmm.devices AS dev
		LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id)
		LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id AND app.status = 1)
		LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id  )
		WHERE dev.user_id = inv.from_user_id AND ( 
		sha.inserted_at > inv.inserted_at OR 
		app.inserted_at > inv.inserted_at OR    
		reading.inserted_at > inv.inserted_at)
		LIMIT 1
		),inv.from_user_id,NULL))  AS today_active,
	   COUNT(distinct IF(inv.inserted_at < DATE(NOW()) ,inv.from_user_id,NULL)) AS yesterday_invite,
	   COUNT(distinct IF(inv.inserted_at > DATE(NOW()) ,inv.from_user_id,NULL)) AS today_invite
FROM
	   tmm.invite_bonus AS inv
WHERE
	    inv.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) AND inv.task_type = 0
) AS inv,
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
	COUNT(IF(created < DATE(NOW()),0,NULL)) AS yesterday_New_users,
	COUNT(IF(created > DATE(NOW()),0,NULL)) AS today_New_users
FROM 
	ucoin.users 
WHERE
	created > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
) AS users,
(
SELECT
	COUNT(distinct IF(
	sha.inserted_at > DATE(NOW()) OR 
	app.inserted_at > DATE(NOW()) OR 
	reading.inserted_at > DATE(NOW()) OR 
	reading.updated_at > DATE(NOW())
	,dev.user_id,NULL)) AS today_active,
	COUNT(distinct IF(
	sha.inserted_at < DATE(NOW()) OR 
	app.inserted_at < DATE(NOW()) OR 
	reading.inserted_at < DATE(NOW()) OR 
	reading.updated_at < DATE(NOW())
	,dev.user_id,NULL)) AS yesterday_active
FROM 
	tmm.devices AS dev
LEFT JOIN tmm.device_share_tasks AS sha ON (sha.device_id = dev.id AND sha.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY))
LEFT JOIN tmm.device_app_tasks AS app ON (app.device_id = dev.id AND app.status = 1 AND app.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) )
LEFT JOIN reading_logs AS reading ON (reading.user_id = dev.user_id AND (reading.inserted_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY) OR reading.updated_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)))
WHERE
	dev.updated_at > DATE_SUB(DATE(NOW()),INTERVAL 1 DAY)
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
