package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

const statsKey = `info-stats-stats`
func StatsHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	context, err := redis.Bytes(redisConn.Do(`GET`, statsKey))
	if context != nil && err ==nil{
		var data  interface{}
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
		points.today_cny+tmm.today_cny AS today_cny,
		points.yesterday_cny+tmm.yesterday_cny AS yesterday_cny,
		device.today_points+inv.today_bonus AS today_point_supply,
		device.yes_points+inv.yesterday_bonus AS yesterday_point_supply,
	    device.today_users AS today_users,
		device.yes_users AS yesterday_users,
		device.today_total AS today_task_number,
		device.yes_total AS yesterday_task_number,
		exchanges.today_tmm AS today_tmm_supply,
		exchanges.yesterday_tmm AS yesterday_tmm_supply,
		inv.today_invite AS today_invite,
		inv.yesterday_invite AS yesterday_invite,
		inv.today_active AS today_active,
		inv.yesterday_active AS yesterday_active
FROM 
(
SELECT 
		COUNT(IF(inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS today_number,
		COUNT(IF(inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS yesterday_number,
		SUM(IF(inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),cny,0)) AS today_cny,
		SUM(IF(inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),cny,0)) AS yesterday_cny
FROM 
		tmm.point_withdraws
WHERE 
		inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY)   
) AS points,
(
SELECT 
		COUNT(IF(inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS today_number,
		COUNT(IF(inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS yesterday_number,
		SUM(IF(inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),cny,0)) AS today_cny,
		SUM(IF(inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),cny,0)) AS yesterday_cny
FROM 
		tmm.withdraw_txs 
WHERE 
		inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY)   AND tx_status = 1
) AS tmm,
( 
SELECT 
		SUM(tmp.today_total) AS today_total,
		SUM(tmp.today_points) AS today_points,
		SUM(tmp.yes_points) AS yes_points,
		SUM(tmp.yes_total) AS yes_total,
		SUM(tmp.today_users) AS today_users,
		SUM(tmp.yes_users) AS yes_users
FROM(
	SELECT 
			SUM(IF(sha.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),sha.points,0)) AS today_points,
			SUM(IF(sha.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),sha.points,0)) AS yes_points,
			COUNT(IF(sha.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),0,NULL)) AS yes_total,
			COUNT(IF(sha.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),0,NULL)) AS today_total,
			COUNT(distinct IF(sha.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),sha.device_id,0)) AS today_users,
			COUNT(distinct IF(sha.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),sha.device_id,0)) AS yes_users

	FROM 
			tmm.device_share_tasks  AS sha
	WHERE 
			sha.inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY)
	UNION ALL

	SELECT 
		    SUM(IF(app.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),app.points,0)) AS today_points,
			SUM(IF(app.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),app.points,0)) AS yes_points,
			COUNT(IF(app.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),0,NULL)) AS yes_total,
			COUNT(IF(app.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),0,NULL)) AS today_total,
			0 AS today_users,
			0 AS yes_users
	FROM
   		    tmm.device_app_tasks AS app 
	WHERE 
  		    app.inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY)
	) AS tmp 
) AS device,
(
SELECT 
	   SUM(IF(inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d '),tmm,0)) AS yesterday_tmm, 
	   SUM(IF(inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d '),tmm,0)) AS today_tmm  
FROM 
	   tmm.exchange_records
WHERE 
	   direction = 1 AND inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY)
)  AS exchanges,
(
SELECT  
	   COUNT(IF(inv.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d') AND dev.updated_at > DATE_SUB(NOW(),INTERVAL 2 DAY),1,NULL))  AS yesterday_active, 
	   COUNT(IF(inv.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d') AND dev.updated_at > DATE_SUB(NOW(),INTERVAL 2 DAY),1,NULL))  AS today_active,
	   COUNT(IF(inv.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS yesterday_invite,
	   COUNT(IF(inv.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),1,NULL)) AS today_invite,
	   SUM(IF(inv.inserted_at > DATE_FORMAT(NOW(),'%Y-%m-%d'),inv.bonus,0)) AS today_bonus,
	   SUM(IF(inv.inserted_at < DATE_FORMAT(NOW(),'%Y-%m-%d'),inv.bonus,0)) AS yesterday_bonus
FROM 
	   tmm.invite_bonus AS inv  
INNER JOIN tmm.devices AS dev ON (dev.user_id = inv.from_user_id)
WHERE 
	   inv.task_id = 0 AND inv.inserted_at > DATE_SUB(DATE_FORMAT(NOW(),'%Y-%m-%d'),INTERVAL 1 DAY) 
) AS inv`

	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `没有找到数据`, c) {
		return
	}
	row := rows[0]
	var yesterdayStats, todayStats StatsData

	yesterdayCash, err := decimal.NewFromString(row.Str(res.Map(`yesterday_cny`)))
	if CheckErr(err, c) {
		return
	}
	yesterdayUcSupply, err := decimal.NewFromString(row.Str(res.Map(`yesterday_tmm_supply`)))
	if CheckErr(err, c) {
		return
	}
	yesterdayPointSupply, err := decimal.NewFromString(row.Str(res.Map(`yesterday_point_supply`)))
	if CheckErr(err, c) {
		return
	}
	yesterdayStats.PointExchangeNumber = row.Int(res.Map(`point_yesterday_number`))
	yesterdayStats.UcoinExchangeNumber = row.Int(res.Map(`tmm_yesterday_number`))
	yesterdayStats.TotalTaskUser = row.Int(res.Map(`yesterday_users`))
	yesterdayStats.TotalFinishTask = row.Int(res.Map(`yesterday_task_number`))
	yesterdayStats.InviteNumber = row.Int(res.Map(`yesterday_invite`))
	yesterdayStats.Active = row.Int(res.Map(`yesterday_active`))
	yesterdayStats.Cash = yesterdayCash
	yesterdayStats.PointSupply = yesterdayPointSupply
	yesterdayStats.UcSupply = yesterdayUcSupply
	todayCash, err := decimal.NewFromString(row.Str(res.Map(`today_cny`)))
	if CheckErr(err, c) {
		return
	}
	todayUcSupply, err := decimal.NewFromString(row.Str(res.Map(`today_tmm_supply`)))
	if CheckErr(err, c) {
		return
	}
	todayPointSupply, err := decimal.NewFromString(row.Str(res.Map(`today_point_supply`)))
	if CheckErr(err, c) {
		return
	}
	todayStats.PointExchangeNumber = row.Int(res.Map(`point_today_number`))
	todayStats.UcoinExchangeNumber = row.Int(res.Map(`tmm_today_number`))
	todayStats.TotalTaskUser = row.Int(res.Map(`today_users`))
	todayStats.TotalFinishTask = row.Int(res.Map(`today_task_number`))
	todayStats.InviteNumber = row.Int(res.Map(`today_invite`))
	todayStats.Active = row.Int(res.Map(`today_active`))
	todayStats.Cash = todayCash
	todayStats.UcSupply = todayUcSupply
	todayStats.PointSupply = todayPointSupply
	var statsList StatsList
	statsList.Yesterday = yesterdayStats
	statsList.Today = todayStats
	bytes,err:=json.Marshal(&statsList)
	if CheckErr(err,c){
		return
	}
	redisConn.Do(`SET`, statsKey, bytes, `EX`, 60)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: statsList,
	})
}
