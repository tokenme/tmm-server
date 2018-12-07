package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
)

const InviteKey = `info-invite`

func InviteHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	context, err := redisConn.Do(`GET`, InviteKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := InviteInfo{}
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
	COUNT(*) AS total ,
	COUNT(IF(inv.task_id = 0 ,0,NULL)) AS invite,
	SUM(IF(inv.task_id = 0,inv.bonus,0)) AS cost
	FROM ucoin.users AS uc 
	INNER JOIN tmm.invite_bonus AS inv ON (inv.from_user_id = uc.id  ) 
	WHERE uc.active =1 `
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	TotalUser := row.Float(res.Map(`total`))
	invite := row.Int(res.Map(`invite`))
	totalCost, err := decimal.NewFromString(row.Str(res.Map(`cost`)))
	if CheckErr(err, c) {
		return
	}
	info := InviteInfo{}
	info.TotalCost = totalCost
	info.TotalInvite = invite
	info.InviteProportionRate = float64(invite) / TotalUser
	//分享页面打开率和成功率 等腾讯那边的数据
	data, err := json.Marshal(&info)
	if CheckErr(err, c) {
		return
	}
	redisConn.Do("SET", InviteKey, data, "EX", KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
