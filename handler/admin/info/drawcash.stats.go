package info

import (
	. "github.com/tokenme/tmm/handler"
	"github.com/gin-gonic/gin"
	"fmt"
	"strings"
	"github.com/shopspring/decimal"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"time"
)

func DrawCashStatsHandler(c *gin.Context) {
	db := Service.Db
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var ptwhen []string
	var txwhen []string
	var startTime, endTime string
	var top10 string
	endTime = time.Now().Format("2006-01-02 15:04:05")
	if req.StartTime != "" {
		startTime = req.StartTime
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  >= '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  >= '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  >= '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  >= '%s' ", db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at <= '%s' ", db.Escape(endTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at <= '%s' ", db.Escape(endTime)))
	}
	if req.Top10 {
		top10 = "LIMIT 10"
	}
	query := `SELECT 
	us.id AS id ,
	wx.nick AS nickname , 
	IFNULL(tmp.cny,0) AS cny,
	us.mobile AS mobile
FROM (
 SELECT 
 user_id, 
 SUM(cny) AS cny
FROM(
	SELECT
            tx.user_id, 
			SUM( tx.cny ) AS cny
        FROM
            tmm.withdraw_txs AS tx
		WHERE
			tx.tx_status = 1 %s
        GROUP BY
            tx.user_id UNION ALL
        SELECT
            pw.user_id, 
			SUM( pw.cny ) AS cny
        FROM
            tmm.point_withdraws AS pw
        WHERE 
			pw.cny > 0 %s
		GROUP BY pw.user_id
				) AS tmp
		GROUP BY user_id
) AS tmp,ucoin.users AS us 
LEFT JOIN tmm.wx AS wx  ON (wx.user_id = us.id)
WHERE tmp.user_id = us.id AND NOT EXISTS  (SELECT 1 FROM user_settings AS us  WHERE us.blocked= 1 AND us.user_id=us.id AND us.block_whitelist=0  LIMIT 1)
GROUP BY us.id 
ORDER BY cny DESC 
%s`
	rows, res, err := db.Query(query, strings.Join(txwhen, ""), strings.Join(ptwhen, ""), db.Escape(top10))
	if CheckErr(err, c) {
		return
	}
	var info DrawCashStats
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Numbers = len(rows)
	info.Title = `提现排行榜`
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: admin.API_OK,
			Data:    info,
		})
		return
	}
	for _, row := range rows {
		cny, err := decimal.NewFromString(row.Str(res.Map(`cny`)))
		if CheckErr(err, c) {
			return
		}

		if req.Top10 {
			user := &admin.Users{
				DrawCash: fmt.Sprintf("%.2f", row.Float(res.Map(`cny`))),
			}
			user.Mobile = row.Str(res.Map(`mobile`))
			user.Id = row.Uint64(res.Map(`id`))
			user.Nick = row.Str(res.Map(`nickname`))
			info.Top10 = append(info.Top10, user)
		}
		info.Money = info.Money.Add(cny)
	}
	info.Money = info.Money.Ceil()
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)
	info.Numbers = len(rows)
	info.Title = `提现排行榜`
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
