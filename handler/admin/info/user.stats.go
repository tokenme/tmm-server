package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type UserStats struct {
	Severs int `json:"severs"`
	Month  int `json:"month"`
	Total  int `json:"total"`
}

func UserStatsHandler(c *gin.Context) {
	db := Service.Db

	query := `SELECT 
	COUNT(*) AS total , 
	COUNT(IF(created > date_sub(NOW(),interval 7 day),0,NULL)) AS serverDay,
	COUNT(IF(created > date_sub(NOW(),interval 1 MONTH),0,NULL)) AS _month 
FROM ucoin.users 
WHERE NOT EXISTS
	(SELECT 1 FROM user_settings AS us
	WHERE us.blocked= 1 AND us.user_id= id  AND us.block_whitelist=0  LIMIT 1)`
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}

	var stats UserStats
	row := rows[0]
	stats.Total = row.Int(0)
	stats.Severs = row.Int(1)
	stats.Month = row.Int(2)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    stats,
	})
}
