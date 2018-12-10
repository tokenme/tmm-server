package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/shopspring/decimal"
	"fmt"
	"strings"
	"time"
)

func InvestsStatsHandler(c *gin.Context) {
	db := Service.Db
	var req StatsRequest
	var startTime, endTime string
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var when []string
	endTime = time.Now().Format("2006-01-02 ")
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(`AND inv.inserted_at >= '%s'`, db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0,0,-7).Format("2006-01-02 ")
		when = append(when, fmt.Sprintf(`AND inv.inserted_at >= '%s'`,db.Escape(startTime)))
	}
	if req.EndTime != "" {
		endTime = req.EndTime
		when = append(when, fmt.Sprintf(`AND inv.inserted_at <= '%s'`,db.Escape(endTime)))
	} else {
		endTime = time.Now().String()
	}
	var top10 string
	if req.Top10{
		top10 = " LIMIT 10"
	}
	query := `
SELECT 
	g.id AS id,
	g.name AS title,
	SUM(inv.points)AS point
FROM tmm.goods AS g
INNER JOIN tmm.good_invests AS inv ON (inv.good_id = g.id )
WHERE inv.redeem_status = 0 %s
GROUP BY g.id 
ORDER BY point DESC %s`
	rows, _, err := db.Query(query, strings.Join(when, " "),top10)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var info InvestsStats
	for _, row := range rows {
		point, err := decimal.NewFromString(row.Str(2))
		if CheckErr(err, c) {
			return
		}
		if len(info.Top10) < 10 {
			info.Top10 = append(info.Top10, &Good{
				Id:    row.Int(0),
				Title: row.Str(1),
				Point: point,
			})
		}
		info.InvestsPoint = info.InvestsPoint.Add(point)
	}
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Numbers = len(rows)
	info.Title = "商品投资排行榜"
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
