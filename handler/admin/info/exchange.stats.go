package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"strings"
	"github.com/shopspring/decimal"
	"fmt"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func ExchangeStatsHandler(c *gin.Context) {
	db := Service.Db
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var when []string
	var startTime, endTime string
	endTime = time.Now().Format("2006-01-02 15:04:05")
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0,0,-7).Format("2006-01-02")
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		when = append(when, fmt.Sprintf(` inserted_at <  DATE_ADD('%s', INTERVAL 60*23+59 MINUTE) `, db.Escape(endTime)))
	}
	var top10 string
	if req.Top10 {
		top10 = " LIMIT 10"
	}

	query := `
SELECT 
	us.id AS id,
	wx.nick AS nickname, 
	tmp.tmm_add  AS tmm,
	us.mobile AS mobile
FROM(
	SELECT
		SUM(tmm) AS  tmm_add,
		er.user_id 
	FROM 
		tmm.exchange_records AS er
	WHERE 
		er.status = 1  AND direction=1 AND %s
	GROUP BY 
		user_id
) AS tmp ,ucoin.users AS us
LEFT JOIN tmm.wx AS wx ON (wx.user_id = us.id)
	WHERE tmp.user_id = us.id AND NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
	ORDER BY tmm DESC
%s
`

	rows, res, err := db.Query(query, strings.Join(when, " AND "), top10)
	if CheckErr(err, c) {
		return
	}
	var info ExchangeStats
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Numbers = len(rows)
	info.Title = `积分兑换UC数量排行榜`
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    info,
		})
		return
	}
	for _, row := range rows {
		tmm, err := decimal.NewFromString(row.Str(res.Map(`tmm`)))
		if CheckErr(err, c) {
			return
		}
		//points, err := decimal.NewFromString(row.Str(res.Map(`points`)))
		//if CheckErr(err, c) {
		//	return
		//}
		//count := row.Int(res.Map(`numbers`))
		if req.Top10 {
			user := &admin.User{
				Tmm:           tmm.Ceil(),
				//Point:         points.Ceil(),
				//ExchangeCount: count,
			}
			user.Mobile = row.Str(res.Map(`mobile`))
			user.Id = row.Uint64(res.Map(`id`))
			user.Nick = row.Str(res.Map(`nickname`))
			info.Top10 = append(info.Top10, user)
		}
		//info.ExchangeCount = info.ExchangeCount + count
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
