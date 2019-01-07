package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"github.com/tokenme/etherscan-api"
	"fmt"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type UserActiveStats struct {
	NewUser       int `json:"new_user"`
	LoginNumber   int `json:"login_number"`
	Active        int `json:"active"`
	Daily         int `json:"daily"`
	Read          int `json:"read"`
	Share         int `json:"share"`
	DownLoadApp   int `json:"down_load_app"`
	Invite        int `json:"invite"`
	DrawCash      int `json:"draw_cash"`
	DrawCashPoint int `json:"draw_cash_point"`
	Transaction   int `json:"transaction"`
}

const (
	UserActiveKey = "trend-active-%s-%s"
)

func UserActiveTrendHandler(c *gin.Context) {
	db := Service.Db
	conn := Service.Redis.Master.Get()
	defer conn.Close()
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	startDate := time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
	endDate := time.Now().Format(`2006-01-02`)
	if req.StartTime != "" {
		startDate = req.StartTime
	}
	if req.EndTime != "" {
		endDate = req.EndTime
	}
	bytes, err := redis.Bytes(conn.Do(`GET`, fmt.Sprintf(UserActiveKey, startDate, endDate)))
	if err == nil && bytes != nil {
		var stats UserActiveStats
		if !CheckErr(json.Unmarshal(bytes, &stats),c){
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    stats,
			})
		}
		return
	}
	query := `
SELECT 
	users.number AS users_number,
	devices._online_number AS _online_number,
	devices.active AS active,
	devices.download_app AS download_app,
	devices.share_article AS share_article,
	devices.daily_number AS daily_number,
	devices.reading_user AS reading_user,
	devices.invite_user AS invite_user,
	devices.draw_user AS draw_user,
	devices.draw_point_user AS draw_point_user
FROM(
SELECT 
	COUNT(1) AS number
FROM 
	ucoin.users 
WHERE created > '%s' AND created < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
) AS users ,
(
SELECT 
	COUNT(DISTINCT dev.user_id) AS _online_number,
	COUNT(DISTINCT IF(daily.user_id > 0 OR app.points > 0 OR sha.points > 0 OR reading.user_id > 0,dev.user_id,NULL)) AS active,
	COUNT(DISTINCT IF(app.points> 0,dev.user_id,NULL)) AS download_app,
	COUNT(DISTINCT IF(sha.points > 0,dev.user_id,NULL)) AS share_article,
	COUNT(DISTINCT daily.user_id ) AS daily_number,
	COUNT(DISTINCT reading.user_id ) AS reading_user,
	COUNT(DISTINCT invite.user_id) AS invite_user,
	COUNT(DISTINCT IF(draw_uc.cny > 0 OR draw_point.cny > 0,dev.user_id,NULL)) AS draw_user,
	COUNT(DISTINCT IF(draw_point.cny > 0,draw_point.user_id,NULL)) AS draw_point_user
FROM  tmm.devices AS dev
INNER JOIN ucoin.users  AS u ON u.id = dev.user_id

LEFT JOIN tmm.device_app_tasks AS app ON
(app.device_id = dev.id AND app.inserted_at > '%s' AND app.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE))

LEFT JOIN tmm.device_share_tasks  AS sha ON 
(sha.device_id = dev.id  AND sha.inserted_at > '%s' AND sha.inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE))

LEFT JOIN (
SELECT 
	user_id AS user_id
FROM tmm.reading_logs
WHERE inserted_at > '%s'  AND inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE) OR (updated_at > '%s' AND updated_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE))
GROUP BY user_id 
) AS reading ON reading.user_id = dev.user_id 

LEFT JOIN tmm.daily_bonus_logs AS daily 
ON (daily.user_id = dev.user_id  AND daily.updated_on > '%s' AND daily.updated_on < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE) )
LEFT JOIN tmm.invite_bonus AS invite 
ON (invite.user_id = dev.user_id AND invite.task_type = 0 AND invite.inserted_at > '%s' AND invite.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE))

LEFT JOIN tmm.withdraw_txs AS draw_uc 
ON (draw_uc.user_id = dev.user_id AND draw_uc.tx_status = 1 AND draw_uc.inserted_at > '%s' AND draw_uc.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE) )

LEFT JOIN tmm.point_withdraws AS draw_point 
ON (draw_point.user_id = dev.user_id AND draw_point.inserted_at > '%s' AND draw_point.inserted_at < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE))

WHERE 
dev.lastping_at > '%s' AND dev.lastping_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)
AND u.created > '%s' AND u.created < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)

) AS devices
	`
	rows, _, err := db.Query(query, db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
		db.Escape(startDate), db.Escape(endDate),
	)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}
	var stats UserActiveStats
	row := rows[0]
	stats.NewUser = row.Int(0)
	stats.LoginNumber = row.Int(1)
	stats.Active = row.Int(2)
	stats.DownLoadApp = row.Int(3)
	stats.Share = row.Int(4)
	stats.Daily = row.Int(5)
	stats.Read = row.Int(6)
	stats.Invite = row.Int(7)
	stats.DrawCash = row.Int(8)
	stats.DrawCashPoint = row.Int(9)
	query = `
SELECT 
	wallet_addr
FROM 
	ucoin.users
WHERE  created > '%s' AND created < DATE_ADD('%s', INTERVAL 60*23+59 MINUTE)

`
	rows, _, err = db.Query(query, startDate, endDate)
	if CheckErr(err, c) {
		return
	}
	addrMap := make(map[string]struct{})
	var count int
	for _, row := range rows {
		addrMap[row.Str(0)] = struct{}{}
	}

	tm, _ := time.Parse(`2006-01-02`, startDate)
	end, _ := time.Parse(`2006-01-02`, endDate)
	geth := Service.Geth
	client := etherscan.New(etherscan.Mainnet, Config.EtherscanAPIKey)
	newBlock, _ := geth.BlockByNumber(c, nil)
	blockTime := time.Unix(newBlock.Time().Int64(), 0)
	startBlock := int(newBlock.Number().Int64()) - int(blockTime.Sub(tm).Seconds()/12)
	address := Config.TMMTokenAddress
	txs, err := client.ERC20Transfers(&address, nil, &startBlock, nil, 0, 100000, false)
	if CheckErr(err, c) {
		return
	}
	for _, tx := range txs {
		txtime := tx.TimeStamp.Time()
		if txtime.After(tm) && txtime.Before(end.AddDate(0, 0, 1)) && tx.From != "0x12c9b5a2084decd1f73af37885fc1e0ced5d5ee8" && tx.From != "0x251e3c2de440c185912ea701a421d80bbe5ee8c9" {
			if _, find := addrMap[tx.From]; find {
				count++
			}
		}
		if txtime.After(end.AddDate(0, 0, 1)) {
			break
		}
	}
	stats.Transaction = count
	bytes,err =  json.Marshal(&stats)
	if CheckErr(err,c){
		return
	}
	_,err=conn.Do(`SET`,fmt.Sprintf(UserActiveKey, startDate, endDate),bytes,`EX`,60*60)
	fmt.Println(err)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    stats,
	})
}
