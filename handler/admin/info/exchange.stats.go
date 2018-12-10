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
	endTime = time.Now().Format("2006-01-02 ")
	if req.StartTime != "" {
		startTime = req.StartTime
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0,0,-7).Format("2006-01-02 ")
		when = append(when, fmt.Sprintf(` inserted_at >= '%s' `, db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		when = append(when, fmt.Sprintf(` inserted_at <= '%s' `, db.Escape(endTime)))
	}
	var top10 string
	if req.Top10 {
		top10 = " LIMIT 10"
	}

	query := `
SELECT 
	u.id AS id,
	u.mobile AS mobile,
	u.nickname AS nickname , 
	tmp.tmm_add - tmp.tmm_ AS tmm,
	tmp.points_add - tmp.points_ AS points,
	tmp.numbers AS numbers
FROM (
	SELECT  
		COUNT(1) AS numbers,
		SUM(IF (direction=1, er.points,0)) AS points_,
		SUM(IF (direction=1, er.tmm,0)) AS  tmm_add,
		SUM(IF (direction=-1, er.points,0)) AS points_add ,
		SUM(IF (direction=-1, er.tmm,0))   AS tmm_,
		er.user_id 
	FROM 
		tmm.exchange_records AS er
	WHERE 
		er.status = 1 AND %s
	GROUP BY 
		user_id
) AS tmp , 	
	ucoin.users AS u 
	WHERE tmp.user_id = u.id
	ORDER BY tmm DESC
%s
`

	rows, res, err := db.Query(query, strings.Join(when, " AND "), top10)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var info ExchangeStats
	for _, row := range rows {
		tmm, err := decimal.NewFromString(row.Str(res.Map(`tmm`)))
		if CheckErr(err, c) {
			return
		}
		points, err := decimal.NewFromString(row.Str(res.Map(`points`)))
		if CheckErr(err, c) {
			return
		}
		count := row.Int(res.Map(`numbers`))
		if req.Top10 {
			user := &Users{
				Tmm:           tmm,
				Point:         points,
				ExchangeCount: count,
			}
			user.Id = row.Uint64(res.Map(`id`))
			user.Nick = row.Str(res.Map(`nickname`))
			user.Mobile = row.Str(res.Map(`mobile`))
			info.Top10 = append(info.Top10, user)
		}
		info.ExchangeCount = info.ExchangeCount + count
	}
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Numbers = len(rows)
	info.Title = `交换Ucoin排行榜`
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
