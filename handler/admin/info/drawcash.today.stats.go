package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"net/http"
)

type WithDrawStats struct {
	WithdrawByPoint  string `json:"withdraw_by_point"`
	UserCountByPoint int    `json:"user_count_by_point"`
	WithdrawByUC     string `json:"withdraw_by_uc"`
	UserCountByUC    int    `json:"user_count_by_uc"`
	MaxDailWithdraw  int64  `json:"max_dail_withdraw"`
	TotalCny         string `json:"total_cny"`
	TotalUser        int    `json:"total_user"`
}

func GetTodayWithDrawStatsHandler(c *gin.Context) {

	date := c.DefaultQuery(`start_date`, time.Now().Format(`2006-01-02`))
	if tm, _ := time.Parse(`2006-01-02`, date); Check(tm.IsZero(), `日期错误`, c) {
		return
	}

	query := `
SELECT
	SUM(IF(tmp.types = 0,tmp.cny,0)) AS point_cny,
	COUNT(DISTINCT IF(tmp.types = 0,tmp.user_id,NULL)) AS point_count,
	SUM(IF(tmp.types = 1,tmp.cny,0)) AS Uc_cny,
	COUNT(DISTINCT IF(tmp.types = 1,tmp.user_id,NULL)) AS Uc_count,
	SUM(tmp.cny) AS  total_cny,
	COUNT(DISTINCT tmp.user_id) AS total_person 
FROM (
	SELECT
		user_id ,
		cny ,
		0 AS types
	FROM
		tmm.point_withdraws
	WHERE
		DATE(inserted_at) = '%s' AND verified != -1 
		
	UNION ALL

	SELECT
		user_id,
		cny,
		1 AS types
	FROM
		tmm.withdraw_txs
	WHERE
		DATE(inserted_at) = '%s' AND verified!=-1 AND tx_status = 1 
  UNION ALL 
	SELECT 
		user_id,
		cny,
		2 AS types
	FROM
		tmm.withdraw_logs
	WHERE 
		DATE(inserted_at) = '%s' 
) AS tmp
	`

	db := Service.Db
	rows, _, err := db.Query(query, db.Escape(date), db.Escape(date), db.Escape(date))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}

	row := rows[0]
	var Stats WithDrawStats
	Stats.MaxDailWithdraw = Config.MaxDailWithdraw
	Stats.WithdrawByPoint = fmt.Sprintf("%.2f", row.Float(0))
	Stats.UserCountByPoint = row.Int(1)
	Stats.WithdrawByUC = fmt.Sprintf("%.2f", row.Float(2))
	Stats.UserCountByUC = row.Int(3)
	Stats.TotalCny = fmt.Sprintf("%.2f", row.Float(4))
	Stats.TotalUser = row.Int(5)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    Stats,
	})
}
