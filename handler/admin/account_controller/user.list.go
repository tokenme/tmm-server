package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"strconv"
)

type withdrawType int
type pointType int

const (
	Any withdrawType = iota
	Point
	Ucoin
)
const (
	All pointType = iota
	invite
	reading
	share
)
const (
	TimeFormat = `%Y-%m-%d`
)

type SearchOptions struct {
	WithdrawType            withdrawType `form:"withdraw_type"`
	WithDrawNumberOfTimes   int          `form:"with_draw_number_of_times"`
	WithDrawAmount          int          `form:"with_draw_amount"`
	ExchangePointToUc       int          `form:"exchange_point_to_uc"`
	ExchangeUcToPoint       int          `form:"exchange_uc_to_point"`
	ExchangePointToUcAmount int          `form:"exchange_point_to_uc_amount"`
	ExchangeUcToPointAmount int          `form:"exchange_uc_to_point_amount"`
	MakePointType           pointType    `form:"make_point_type"`
	MakePointTimes          int          `form:"make_point_times"`
	MakePointDay            int          `form:"make_point_day"`
	OnlineBFNumber          int          `form:"online_bf_number"`
	OffLineBfNumber         int          `form:"off_line_bf_number"`
	StartDate               string       `form:"start_date"`
	EndDate                 string       `form:"end_date"`
	StartHours              string       `form:"start_hours"`
	EndHours                string       `form:"end_hours"`
	IsWhiteList             bool         `form:"is_white_list"`
	Id                      int          `form:"id"`
	Mobile                  string       `form:"mobile"`
}

func GetAccountList(c *gin.Context) {
	db := Service.Db
	var search SearchOptions
	if CheckErr(c.Bind(&search), c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, `0`))
	if CheckErr(err, c) {
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery(`limit`, `10`))
	if CheckErr(err, c) {
		return
	}
	var offset int
	if limit < 1 {
		limit = 10
	}
	if page > 0 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}
	query := `
SELECT
	u.id AS id,
	wx.nick AS nick,
	u.mobile AS mobile,
	cny.total AS cny_total,
	IFNULL(cny.cny,0)  AS  cny,
	point.total AS point_total,
	IFNULL(point.point,0) AS point,
	ex.exchange_total AS exchange_total,
	ex.point_to_tmm_times AS point_to_tmm_times,
	ex.tmm_to_point_times AS tmm_to_point_times,
	IFNULL(ex.point_to_tmm,0) AS point_to_tmm,
	ex.tmm_to_point AS tmm_to_point,
	IFNULL(inv.online,0) AS online,
	IFNULL(inv.offline,0) AS offline,
	IF(us_set.block_whitelist = NULL,0,IF(us_set.block_whitelist=us_set.block_whitelist,0,1)) AS blocked
FROM 
	ucoin.users AS u
LEFT JOIN tmm.wx AS wx ON (wx.user_id = u.id)
LEFT JOIN tmm.user_settings AS us_set ON (us_set.user_id = u.id)
LEFT JOIN (
SELECT 
	user_id,
	COUNT(1) AS exchange_total,
	COUNT(IF(direction = 1,0,NULL)) AS point_to_tmm_times,
	COUNT(IF(direction = -1,0,NULL)) AS tmm_to_point_times,
	SUM(IF(direction = 1 ,tmm,0)) AS point_to_tmm,
	SUM(IF(direction = -1 ,points,0)) AS tmm_to_point
FROM 
	tmm.exchange_records 
WHERE 
	status = 1 %s
GROUP BY 
	user_id 
) AS ex ON (ex.user_id = u.id)
LEFT JOIN (
SELECT 
	tmp.id AS id,
	SUM(tmp.online) AS online,
	SUM(tmp.offline) AS offline
FROM (
SELECT
	inv.parent_id AS id,
	COUNT(distinct IF(dev.lastping_at > DATE_SUB(NOW(),INTERVAL 1 DAY),inv.user_id,NULL))   AS online,
	COUNT(distinct IF(dev.lastping_at < DATE_SUB(NOW(),INTERVAL 1 DAY),inv.user_id,NULL))   AS offline
FROM 	
	invite_codes AS inv 
INNER JOIN tmm.devices AS dev ON (dev.user_id = inv.user_id)
GROUP BY  id UNION ALL
SELECT 
	inv.grand_id AS id,
	COUNT(distinct IF(dev.lastping_at > DATE_SUB(NOW(),INTERVAL 1 DAY),inv.user_id,NULL))   AS online,
	COUNT(distinct IF(dev.lastping_at < DATE_SUB(NOW(),INTERVAL 1 DAY),inv.user_id,NULL))   AS offline
FROM 	
	invite_codes AS inv 
INNER JOIN tmm.devices AS dev ON (dev.user_id = inv.user_id)
GROUP BY  id
) AS tmp 
GROUP BY tmp.id
) AS inv ON (inv.id = u.id)
%s 
WHERE 
    1 = 1 %s
GROUP BY 
	id
ORDER BY 
	id  ,cny DESC 
LIMIT %d OFFSET %d

`
	var leftJoin []string
	var when []string
	var where []string
	if search.StartDate != "" {
		when = append(when, fmt.Sprintf(` AND inserted_at > %s`, db.Escape(search.StartDate)))
	}
	if search.EndDate != "" {
		when = append(when, fmt.Sprintf(` AND inserted_at < %s`, db.Escape(search.EndDate)))
	}
	if search.StartHours != "" {
		when = append(when, fmt.Sprintf(` AND HOUR(inserted_at)  BETWEEN %s  AND %s `, db.Escape(search.StartHours), db.Escape(search.EndHours)))
	}
	switch search.WithdrawType {
	case Any:
		leftJoin = append(leftJoin, fmt.Sprintf(`
LEFT JOIN (
	SELECT 
	tmp.user_id AS user_id,
	SUM(tmp.cny) AS cny ,
	SUM(tmp.total) AS total
FROM (
	SELECT
		user_id, 
		IFNULL(SUM(cny),0) AS cny,
		COUNT(1) AS total
	FROM
		tmm.withdraw_txs
	WHERE
		tx_status = 1  %s
	GROUP BY
		user_id UNION ALL
	SELECT
		user_id, 
		IFNULL(SUM(cny),0) AS cny,
		COUNT(1) AS total 
	FROM
		tmm.point_withdraws 
	WHERE
		1 = 1  %s
	GROUP BY 
		user_id
	) AS tmp
	GROUP BY 
		user_id 
) AS cny ON (cny.user_id = u.id )`, strings.Join(when, ` `), strings.Join(when, ` `)))
	case Point:
		leftJoin = append(leftJoin, fmt.Sprintf(fmt.Sprintf(` 
LEFT JOIN (
	SELECT
		user_id AS user_id, 
		SUM( cny ) AS cny,
		COUNT(1) AS total
	FROM
		tmm.point_withdraws  
	WHERE
		1 = 1 %s
	GROUP BY  
		user_id 
) AS cny ON (cny.user_id = u.id )`, strings.Join(when, ` `))))
	case Ucoin:
		leftJoin = append(leftJoin, fmt.Sprintf(` 
LEFT JOIN (
	SELECT
		user_id, 
		SUM( cny ) AS cny,
		COUNT(1) AS total
	FROM
		tmm.withdraw_txs 
	WHERE 
		tx_status = 1 %s
	GROUP BY 
		user_id
) AS cny ON (cny.user_id = u.id )
 `, strings.Join(when, ` `)))
	default:
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: "未知参数",
		})
		return
	}
	switch search.MakePointType {
	case All:
		leftJoin = append(leftJoin, fmt.Sprintf(`
LEFT JOIN (
	SELECT 
		dev.user_id AS id,
		 IFNULL(SUM(tmp.points),0) + IFNULL(inv.bonus,0) + IFNULL(read_.point,0) AS point,
		 IFNULL(SUM(tmp.total),0) + IFNULL(inv.total,0) + IFNULL(read_.total,0) AS  total,
		COUNT(distinct tmp.date) AS _day
	FROM 
		tmm.devices AS dev 
	LEFT JOIN (
		SELECT 
			device_id, 
			SUM(points) AS points,
			COUNT(1) AS total,
			DATE_FORMAT(inserted_at,'%s') AS date
		FROM 
			tmm.device_share_tasks
		WHERE 
			points > 0 %s
		GROUP BY
			device_id,date UNION ALL
		SELECT 
			device_id, 
			SUM(points) AS points,
			COUNT(1) AS total,
			DATE_FORMAT(inserted_at,'%s') AS date
		FROM 
			tmm.device_app_tasks   
		WHERE
			status = 1 %s
		GROUP BY
			device_id,date   
		) AS tmp ON (tmp.device_id =dev.id)
	LEFT JOIN (
		SELECT 
			SUM(bonus) AS bonus,
			user_id AS user_id ,
			COUNT(1) AS total 
		FROM 
			tmm.invite_bonus  
		WHERE
			1 = 1 %s
		GROUP BY 
			user_id
	)AS inv ON (inv.user_id = dev.user_id ) 
	LEFT JOIN (
		SELECT
			SUM(point) AS point,
			user_id AS user_id,
			COUNT(1) AS total 
		FROM
			tmm.reading_logs 
		WHERE 
			1 = 1 %s
		GROUP BY
			user_id
	) AS read_ ON (read_.user_id = dev.user_id)
	GROUP BY 
		dev.user_id
) AS point ON (point.id = u.id)`, TimeFormat, strings.Join(when, ` `),
			TimeFormat, strings.Join(when, ` `), strings.Join(when, ` `),strings.Join(when,` `)))
	case invite:
		leftJoin = append(leftJoin, fmt.Sprintf(` 
LEFT JOIN (
	SELECT 
		SUM(bonus) AS point,
		user_id AS id ,
		COUNT(1) AS total,
		COUNT(distinct DATE_FORMAT(inserted_at,'%s')) AS _day
	FROM 
		tmm.invite_bonus 
	WHERE
		task_id = 0 %s
	GROUP BY 
		user_id
) AS point ON (point.id = u.id)`, TimeFormat, strings.Join(when, " ")))
	case reading:
		leftJoin = append(leftJoin, fmt.Sprintf(`
LEFT JOIN (
	SELECT 
		user_id AS id,
		SUM(point) AS point,
		COUNT(1)  AS total,
		COUNT(distinct DATE_FORMAT(inserted_at,'%s')) AS _day
	FROM 
		tmm.reading_logs 
	WHERE 
		1 = 1 %s
	GROUP BY 
	  user_id
) AS point ON (point.id = u.id)`, TimeFormat, strings.Join(when, " ")))
	case share:
		leftJoin = append(leftJoin, fmt.Sprintf(`
LEFT JOIN (
	SELECT 
		dev.user_id AS id,
		tmp.points AS point,
		tmp.total AS total,
		COUNT(distinct DATE_FORMAT(inserted_at,'%s')) AS _day
	FROM 
		tmm.devices AS dev  
	LEFT JOIN (
		SELECT 
			SUM(points) AS points,
			device_id  AS id ,
			COUNT(1) AS total
		FROM 
			tmm.device_share_tasks
		WHERE 	
    		1 = 1 %s
		GROUP BY 
			id
	) AS tmp ON (tmp.id = dev.id)
	GROUP BY 
	user_id
) AS point ON (point.id = u.id)`, TimeFormat, strings.Join(when, "")))
	default:
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: "未知参数",
		})
		return
	}
	if search.Id != 0 {
		where = append(where, fmt.Sprintf(`AND u.id = %d`, search.Id))
	}
	if search.Mobile != "" {
		where = append(where, fmt.Sprintf(` AND u.mobile = '%s'`, search.Mobile))
	}
	if search.WithDrawNumberOfTimes != 0 {
		where = append(where, fmt.Sprintf(" AND cny.total >= %d", search.WithDrawNumberOfTimes))
	}
	if search.WithDrawAmount != 0 {
		where = append(where, fmt.Sprintf(" AND cny >= %d", search.WithDrawAmount))
	}
	if search.ExchangePointToUc != 0 {
		where = append(where, fmt.Sprintf(` AND point_to_tmm_times >= %d `, search.ExchangePointToUc))
	}
	if search.ExchangePointToUcAmount != 0 {
		where = append(where, fmt.Sprintf(` AND point_to_tmm >= %d`, search.ExchangePointToUcAmount))
	}
	if search.ExchangeUcToPoint != 0 {
		where = append(where, fmt.Sprintf(` AND tmm_to_point_times >= %d`, search.ExchangeUcToPoint))
	}
	if search.ExchangeUcToPointAmount != 0 {
		where = append(where, fmt.Sprintf(` AND tmm_to_point >= %d`, search.ExchangeUcToPointAmount))
	}
	if search.MakePointTimes != 0 {
		where = append(where, fmt.Sprintf(` AND point.total >= %d`, search.MakePointTimes))
	}
	if search.MakePointDay != 0 {
		where = append(where, fmt.Sprintf(` AND IFNULL(point.point,1) / IFNULL(point._day,1) >= %d`, search.MakePointDay))
	}
	if search.OnlineBFNumber != 0 {
		where = append(where, fmt.Sprintf(` AND online >= %d`, search.OnlineBFNumber))
	}
	if search.OffLineBfNumber != 0 {
		where = append(where, fmt.Sprintf(` AND offline >= %d`, search.OffLineBfNumber))
	}
	if search.IsWhiteList {
		where = append(where, fmt.Sprintf(`  AND blocked = %d `, 1))
	}
	rows, res, err := db.Query(query, strings.Join(when, " "),
		strings.Join(leftJoin, " "), strings.Join(where, " "), limit, offset)
	if CheckErr(err, c) {
		return
	}
	var List []*admin.Users
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.Not_Found,
			Data: gin.H{
				"data":  List,
				"page":  page,
				"total": 100,
			},
		})
		return
	}

	for _, row := range rows {
		point, err := decimal.NewFromString(row.Str(res.Map(`point`)))
		if CheckErr(err, c) {
			return
		}
		pointToUcoin, err := decimal.NewFromString(row.Str(res.Map(`point_to_tmm`)))
		if CheckErr(err, c) {
			return
		}
		user := &admin.Users{
			Point:                point.Ceil(),
			DrawCash:             fmt.Sprintf("%.2f",row.Float(res.Map(`cny`))),
			ExchangeCount:        row.Int(res.Map(`point_to_tmm_times`)),
			OnlineBFNumber:       row.Int(res.Map(`online`)),
			OffLineBFNumber:      row.Int(res.Map(`offline`)),
			ExchangePointToUcoin: pointToUcoin.Ceil(),
			Blocked:              row.Int(res.Map(`blocked`)),
		}
		user.Nick = row.Str(res.Map(`nick`))
		user.Id = row.Uint64(res.Map(`id`))
		user.Mobile = row.Str(res.Map(`mobile`))
		List = append(List, user)
	}
	var total int
	rows, _, err = db.Query(`SELECT COUNT(id) FROM  ucoin.users`)
	if err == nil || len(rows) != 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"data":  List,
			"page":  page,
			"total": total,
		},
	})
}
