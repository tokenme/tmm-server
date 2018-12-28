package account_controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"fmt"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
)

func UserInfoHandler(c *gin.Context) {
	db := Service.Db
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	query := `
SELECT 
	u.id AS id , 
	u.mobile AS mobile ,
	wx.nick AS nick,
	uc.cny AS uc_cny,
	point.cny AS point_cny,
	IFNULL(uc.cny,0)+IFNULL(point.cny,0) AS cny,
	SUM(dev.points) AS points,
	inv.direct AS direct,
	inv.indirect AS indirect,
	inv.active AS active,
	bonus.inv_bonus AS inv_bonus,
	IFNULL(sha.points,0) AS sha_points,
	reading.point AS reading_point,
	u.wallet_addr AS addr
FROM 
	devices AS dev
INNER JOIN  
	ucoin.users AS u  ON (u.id = dev.user_id)
LEFT JOIN 
	tmm.wx AS wx ON (wx.user_id =dev.user_id),
(
	SELECT 
		SUM(cny) AS cny
	FROM 
		tmm.withdraw_txs
	WHERE 
		user_id = %d AND tx_status = 1
) AS uc,
(
	SELECT 
		SUM(cny) AS cny
	FROM 
		tmm.point_withdraws 
		WHERE user_id = %d
) AS point,
(
	SELECT
		COUNT(distinct IF(inv.parent_id = %d,inv.user_id,NULL)) AS direct,
		COUNT(distinct IF(inv.grand_id = %d,inv.user_id,NULL)) AS indirect,
		COUNT(distinct IF(dev.updated_at > DATE_SUB(NOW(),INTERVAL 3 DAY) AND (inv.parent_id = %d OR inv.grand_id = %d),inv.user_id,NULL)) AS active
	FROM 
		tmm.invite_codes  AS inv 
	LEFT JOIN tmm.devices AS dev ON (dev.user_id = inv.user_id)
) AS inv,
(
	SELECT 
		SUM(sha.points) AS points 
	FROM(
		SELECT
			IFNULL(SUM(sha.points),0) AS points
		FROM
			tmm.device_share_tasks AS sha
		INNER JOIN
			tmm.devices AS dev ON  (dev.id = sha.device_id)
		WHERE
			 dev.user_id = %d UNION ALL
   		SELECT
			IFNULL(SUM(app.points),0) AS points
		FROM
			tmm.device_app_tasks AS app
		INNER JOIN
			tmm.devices AS dev ON  (dev.id = app.device_id)
		WHERE
			 dev.user_id = %d AND app.status = 1
		 ) AS sha
) AS sha,
(
	SELECT
		IFNULL(SUM(bonus),0) AS inv_bonus
	FROM 
		tmm.invite_bonus 
	WHERE 
		user_id = %d
) AS bonus,
(
	SELECT 
		IFNULL(SUM(point),0) AS point
	FROM
		tmm.reading_logs
	WHERE user_id = %d
) AS reading
WHERE 
	u.id = %d
	`
	if id <= 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: `错误参数`,
		})
		return
	}
	rows, res, err := db.Query(query, id, id, id, id, id, id, id, id, id, id, id)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
		})
		return
	}
	row := rows[0]
	point, err := decimal.NewFromString(row.Str(res.Map(`points`)))
	if CheckErr(err, c) {
		return
	}

	tokenABI, err := utils.NewToken(Config.TMMTokenAddress, Service.Geth)
	if CheckErr(err, c) {
		return
	}
	_, _, decimals, _, _, _, _, _, balance, err := utils.TokenMeta(tokenABI, row.Str(res.Map(`addr`)))
	balanceDecimal, err := decimal.NewFromString(balance.String())
	tmm := balanceDecimal.Div(decimal.New(1, int32(decimals)))
	user := &admin.Users{
		Point:           point.Ceil(),
		DrawCash:        fmt.Sprintf("%.2f", row.Float(res.Map(`cny`))),
		DrawCashByUc:    fmt.Sprintf("%.2f", row.Float(res.Map(`uc_cny`))),
		DrawCashByPoint: fmt.Sprintf("%.2f", row.Float(res.Map(`point_cny`))),
		DirectFriends:   row.Int(res.Map(`direct`)),
		IndirectFriends: row.Int(res.Map(`indirect`)),
		ActiveFriends:   row.Int(res.Map(`active`)),
		Tmm:             tmm.Ceil(),
		PointByShare:    int(row.Float(res.Map(`sha_points`))),
		PointByReading:  int(row.Float(res.Map(`reading_point`))),
		PointByInvite:   int(row.Float(res.Map(`inv_bonus`))),
	}
	user.ChildrenNumber = user.DirectFriends + user.IndirectFriends
	user.TotalMakePoint = user.PointByShare + user.PointByReading + user.PointByInvite
	user.Id = row.Uint64(res.Map(`id`))
	user.Mobile = row.Str(res.Map(`mobile`))
	user.Nick = row.Str(res.Map(`nick`))
	rows, _, err = db.Query("SELECT id,platform FROM tmm.devices WHERE user_id = %d", id)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0 || err != nil, admin.Not_Found, c) {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
		})
		return
	}
	for _,row:=range rows{
		device := &common.Device{}
		device.Id = row.Str(0)
		device.Platform = row.Str(1)
		user.DeviceList = append(user.DeviceList,device)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    user,
	})
}
