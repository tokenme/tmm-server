package account_controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"fmt"
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
	IFNULL(uc.cny+point.cny,0) AS cny,
	IFNULL(tmp.tmm-uc.tmm,0) AS tmm,
	SUM(dev.points) AS points,
	inv.direct AS direct,
	inv.indirect AS indirect,
	inv.online AS online,
	inv.active AS active,
	bonus.inv_bonus AS inv_bonus,
	sha.points AS sha_points,
	reading.point AS reading_point
FROM 
	devices AS dev
INNER JOIN  
	ucoin.users AS u  ON (u.id = dev.user_id)
LEFT JOIN 
	tmm.wx AS wx ON (wx.user_id =dev.user_id),
(
	SELECT 
		SUM(cny) AS cny,
		SUM(tmm) AS tmm
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
		SUM(IF(direction = 1,tmm,0))-SUM(IF(direction = -1,tmm,0)) AS tmm
	FROM 
  		tmm.exchange_records 
	WHERE
		status = 1 AND user_id = %d
) AS tmp,
(
	SELECT
		COUNT(IF(inv.parent_id = %d,0,NULL)) AS direct,
		COUNT(IF(inv.grand_id = %d,0,NULL)) AS indirect,
		COUNT(IF(dev.lastping_at > DATE_SUB(NOW(),INTERVAL 1 DAY) AND inv.parent_id = %d,1,NULL)) AS online,
		COUNT(IF(dev.lastping_at > DATE_SUB(NOW(),INTERVAL 3 DAY) AND inv.parent_id = %d,1,NULL)) AS active
	FROM 
		tmm.invite_codes  AS inv 
	INNER JOIN tmm.devices AS dev ON (dev.user_id = inv.user_id)
) AS inv,
(
	SELECT 
		SUM(sha.points) AS points 
	FROM 
		tmm.device_share_tasks AS sha  
	INNER JOIN 
		tmm.devices AS dev ON  (dev.id = sha.device_id)
	WHERE
		 dev.user_id = %d UNION ALL
    SELECT
		SUM(app.points) AS points 
	FROM 
		tmm.device_app_tasks AS app
	INNER JOIN 
		tmm.devices AS dev ON  (dev.id = app.device_id)
	WHERE
		 dev.user_id = %d
) AS sha,

(
	SELECT
		SUM(bonus) AS inv_bonus
	FROM 
		tmm.invite_bonus 
	WHERE 
		user_id = %d
) AS bonus,
(
	SELECT 
		SUM(point) AS point
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
	rows, res, err := db.Query(query, id, id, id, id, id, id, id, id, id, id, id, id)
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
	tmm, err := decimal.NewFromString(row.Str(res.Map(`tmm`)))
	if CheckErr(err, c) {
		return
	}
	user := &admin.Users{
		Point:           point,
		DrawCash:        fmt.Sprintf("%.2f", row.Float(res.Map(`cny`))),
		DrawCashByUc:    fmt.Sprintf("%.2f", row.Float(res.Map(`uc_cny`))),
		DrawCashByPoint: fmt.Sprintf("%.2f", row.Float(res.Map(`point_cny`))),
		Tmm:             tmm,
		DirectFriends:   row.Int(res.Map(`direct`)),
		IndirectFriends: row.Int(res.Map(`indirect`)),
		OnlineBFNumber:  row.Int(res.Map(`online`)),
		ActiveFriends:   row.Int(res.Map(`active`)),
		PointByShare:    int(row.Float(res.Map(`sha_points`))),
		PointByReading:  int(row.Float(res.Map(`reading_point`))),
		PointByInvite:   int(row.Float(res.Map(`inv_bonus`))),
	}
	user.TotalMakePoint = user.PointByShare + user.PointByReading + user.PointByInvite
	user.Id = row.Uint64(res.Map(`id`))
	user.Mobile = row.Str(res.Map(`mobile`))
	user.Nick = row.Str(res.Map(`nick`))

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    user,
	})
}
