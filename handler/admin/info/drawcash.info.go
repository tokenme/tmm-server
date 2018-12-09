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

func DrawCashInfoHandler(c *gin.Context) {
	db := Service.Db
	var req InfoRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var ptwhen []string
	var txwhen []string
	var startTime, endTime string
	var top10 string
	endTime = time.Now().Format("2006-01-02")
	fmt.Println(req)
	if req.StartTime != "" {
		startTime = req.StartTime
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  >= '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  >= '%s' ", db.Escape(startTime)))
	} else {
		startTime = time.Now().AddDate(0, 0, -7).Format("2006-01-02 ")
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at  >= '%s' ", db.Escape(startTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at  >= '%s' ", db.Escape(startTime)))
	}

	if req.EndTime != "" {
		endTime = req.EndTime
		ptwhen = append(ptwhen, fmt.Sprintf(" AND pw.inserted_at <= '%s' ", db.Escape(endTime)))
		txwhen = append(txwhen, fmt.Sprintf(" AND tx.inserted_at <= '%s' ", db.Escape(endTime)))
	}
	fmt.Println(req.Top10)
	if req.Top10 {
		top10 = "LIMIT 10"
	}
	query := `SELECT 
	u.id AS id,
	u.mobile AS mobile,
	u.nickname AS nickname , 
	u.wx_nick AS wx_nick,
	tmp.cny AS cny
FROM (
 SELECT 
 user_id, 
 SUM(cny) AS cny
FROM(
	SELECT
            tx.user_id, SUM( tx.cny ) AS cny
        FROM
            tmm.withdraw_txs AS tx
		WHERE
			tx.tx_status = 1 %s
        GROUP BY
            tx.user_id UNION ALL
        SELECT
            pw.user_id, SUM( pw.cny ) AS cny
        FROM
            tmm.point_withdraws AS pw
        WHERE 
			pw.cny > 0 %s
		GROUP BY pw.user_id
				) AS tmp
		GROUP BY user_id
) AS tmp,ucoin.users AS u 
WHERE user_id = u.id
GROUP BY id 
ORDER BY cny DESC 
%s`
	rows, res, err := db.Query(query, strings.Join(txwhen, ""), strings.Join(ptwhen, ""), db.Escape(top10))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var info DrawCashInfo
	for _, row := range rows {
		cny, err := decimal.NewFromString(row.Str(res.Map(`cny`)))
		if CheckErr(err, c) {
			return
		}

		if req.Top10 {
			User := &User{
				Id:       row.Int(res.Map(`id`)),
				Mobile:   row.Str(res.Map(`mobile`)),
				Nick:     row.Str(res.Map(`nickname`)),
				WxNick:   row.Str(res.Map(`wx_nick`)),
				DrawCash: cny,
			}
			info.Top10 = append(info.Top10, User)
		}
		info.Money = info.Money.Add(cny)

	}
	info.CurrentTime = fmt.Sprintf("%s-%s", startTime, endTime)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
