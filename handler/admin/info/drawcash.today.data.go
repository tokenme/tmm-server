package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/tokenme/tmm/handler/admin/account_controller"
	"net/http"
	"strconv"
	"time"
	"github.com/tokenme/tmm/handler/admin/withdraw"
)

type Transaction struct {
	Id             uint   `json:"id"`
	Mobile         string `json:"mobile"`
	Nick           string `json:"nick"`
	AccountCreated string `json:"account_created"`
	Types          string `json:"types"`
	InsertedAt     string `json:"inserted_at"`
	Cny            string `json:"cny"`
	Status         string `json:"status"`
	Verified       string `json:"verified"`
}

const (
	WithDrawMoneyByPoint = iota
	WithDrawMoneyByUc
)

var DrawType = map[int]string{
	WithDrawMoneyByPoint: "积分提现",
	WithDrawMoneyByUc:    "UC提现",
}

func GetWithDrawDataHandler(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery(`page`, `1`))
	if CheckErr(err, c) {
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery(`limit`, `10`))
	if CheckErr(err, c) {
		return
	}
	date := c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`))
	if tm, _ := time.Parse(`2006-01-02`, date); Check(tm.IsZero(), `日期错误`, c) {
		return
	}

	var offset int
	if page < 0 || pageSize < 0 {
		offset = 0
		pageSize = 10
	} else {
		offset = pageSize * (page - 1)
	}

	query := `
SELECT 
	tmp.user_id,
	tmp.cny, 
	tmp.inserted_at,
	tmp.tx_status,
	tmp.created,
	tmp.mobile,
	tmp.types,
	tmp.nick,
	tmp.verified
FROM (
	SELECT 
		point.user_id AS user_id ,
		point.cny  AS cny,
		point.inserted_at AS inserted_at ,
		1 AS tx_status,
		u.created,
		u.mobile ,
		0 AS types,
		IFNULL(wx.nick,u.nickname) AS nick,
		point.verified AS verified
	FROM 
		tmm.point_withdraws AS point
	INNER JOIN 
		ucoin.users AS u  ON u.id = user_id
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = u.id 
	WHERE 
		DATE(point.inserted_at) = '%s'

	UNION ALL

	SELECT 
		uc.user_id AS user_id ,
		uc.cny AS cny, 
		uc.inserted_at AS inserted_at,
		uc.withdraw_status AS withdraw_status,
		u.created,
		u.mobile ,
		1 AS types,
		IFNULL(wx.nick,u.nickname) AS nick,
		uc.verified AS verified
	FROM 
		tmm.withdraw_txs AS uc
	INNER JOIN 
		ucoin.users AS u  ON u.id = user_id 	
	LEFT JOIN 
		tmm.wx AS wx ON wx.user_id = u.id
	WHERE 
		DATE(uc.inserted_at) = '%s'    
) AS tmp
ORDER BY tmp.inserted_at DESC 
LIMIT %d OFFSET %d

	`
	totalQuery := `
SELECT 
	SUM(tmp.number)
FROM (
	SELECT 
		COUNT(1) AS number
	FROM 
		tmm.point_withdraws 	
	WHERE 
		DATE(inserted_at) = '%s'

	UNION ALL

	SELECT 
		COUNT(1) AS number
	FROM 
		tmm.withdraw_txs
		
	WHERE 
		DATE(inserted_at) = '%s'  
) AS tmp`
	db := Service.Db
	rows, _, err := db.Query(query, db.Escape(date), db.Escape(date), pageSize, offset)
	if CheckErr(err, c) {
		return
	}

	var list []*Transaction
	for _, row := range rows {
		list = append(list, &Transaction{
			Id:             row.Uint(0),
			Cny:            fmt.Sprintf(`%.2f`, row.Float(1)),
			InsertedAt:     row.Str(2),
			Status:         account_controller.MsgMap[row.Int(3)],
			AccountCreated: row.Str(4),
			Mobile:         row.Str(5),
			Types:          DrawType[row.Int(6)],
			Nick:           row.Str(7),
			Verified:       withdraw.StatsMap[row.Int(8)],
		}, )
	}

	rows, _, err = db.Query(totalQuery, db.Escape(date), db.Escape(date))
	if CheckErr(err, c) {
		return
	}
	var total int
	if len(rows) > 0 {
		total = rows[0].Int(0)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"page":  page,
			"data":  list,
			"total": total,
		},
	})
}
