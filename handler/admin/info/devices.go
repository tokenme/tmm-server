package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const devicesKey = `info-devices`

func DevicesHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	context, err := redisConn.Do(`GET`, devicesKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := DevicesInfo{}
		if json.Unmarshal(context.([]byte), &info) == nil {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    info,
			})
			return
		}
	}
	query := `SELECT 
	COUNT(0) AS total_devices,
	COUNT(IF(platform='ios',0,NULL)) AS ios, 
	COUNT(IF(platform='android',0,NULL)) AS android 
	FROM tmm.devices`
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	totalDevices := row.Int(res.Map(`total_devices`))
	totalIosDevices := row.Int(res.Map(`ios`))
	totalAndroidDevices := row.Int(res.Map(`android`))
	info := DevicesInfo{}
	info.TotalDevices = totalDevices
	info.TotalIosDevices = totalIosDevices
	info.TotalAndroidDevices = totalAndroidDevices

	query = `SELECT COUNT(*) AS total ,
	COUNT(IF(uc.created > DATE_SUB(NOW(),INTERVAL 1 DAY),0,NULL))  AS _day,
	COUNT(IF(uc.created > DATE_SUB(NOW(),INTERVAL 1 WEEK),0,NULL))  AS _week,
	COUNT(IF(uc.created > DATE_SUB(NOW(),INTERVAL 1 MONTH),0,NULL))  AS _month,
	SUM(inv.invites) AS invite_user 
	FROM ucoin.users AS uc 
	INNER JOIN tmm.top_invites_users AS inv ON (inv.mobile = uc.mobile) 
	WHERE uc.active =1 `
	rows, res, err = db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row = rows[0]
	info.UserDownloadChannel.InviteDownload = row.Int(res.Map(`invite_user`))
	info.UserDownloadChannel.NormalDownload = row.Int(res.Map(`total`)) - row.Int(res.Map(`invite_user`))
	info.NewUser.Day = row.Int(res.Map(`_day`))
	info.NewUser.Week = row.Int(res.Map(`_week`))
	info.NewUser.Month = row.Int(res.Map(`_month`))
	data, err := json.Marshal(&info)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do("SET", devicesKey, data, "EX", KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
