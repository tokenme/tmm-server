package info

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"time"
)

/*
6. UC提现金额/UC人数趋势图
7.积分提现金/积分人数趋势图
8. 提现总金额趋势图
*/

const (
	DrawCashByUc = iota
	UcPerson
	DrawCashByPoint
	PointPerson
	TotalDrawCash_
)

func DrawCashTrendHandler(c *gin.Context) {
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	startTime := time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
	endTime := time.Now().Format(`2006-01-02`)
	if req.StartTime != "" {
		startTime = req.StartTime
	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}

	s, _ := time.Parse(`2006-01-02`, startTime)
	e, _ := time.Parse(`2006-01-02`, endTime)
	if s.Unix() > e.Unix() {
		c.JSON(http.StatusOK, admin.Response{
			Code:    1,
			Message: "起始日期不能超过结束日期",
			Data:    nil,
		})
		return
	}

	var data Data
	var title,query string
	seriesName := "金额"
	yaxisName := "金额"
	db := Service.Db
	switch req.Type {
	case DrawCashByUc:
		title = "UC提现金额"
		query = fmt.Sprintf(`
SELECT
	SUM(cny) AS _value,
	DATE(inserted_at) AS record_on,
	0 AS cash
FROM tmm.withdraw_txs
WHERE tx_status = 1 AND inserted_at > '%s'  AND inserted_at <  DATE_ADD('%s', INTERVAL 1 DAY)
GROUP BY record_on`, db.Escape(startTime), db.Escape(endTime))
	case UcPerson:
		title = "UC提现人数"
		yaxisName = "人数"
		seriesName = "人数"
		query = fmt.Sprintf(`
SELECT
	COUNT(distinct user_id) AS _value,
	DATE(inserted_at) AS record_on,
	0 AS cash
FROM tmm.withdraw_txs
WHERE tx_status = 1 AND inserted_at > '%s'  AND inserted_at <  DATE_ADD('%s', INTERVAL 1 DAY)
GROUP BY record_on`, db.Escape(startTime), db.Escape(endTime))
	case DrawCashByPoint:
		title = "积分提现金额"
		query = fmt.Sprintf(`
SELECT
	SUM(cny) AS _value,
	DATE(inserted_at) AS record_on,
	0 AS cash
FROM tmm.point_withdraws
WHERE inserted_at > '%s' AND inserted_at <  DATE_ADD('%s', INTERVAL 1 DAY) AND verified!=-1
GROUP BY record_on`, db.Escape(startTime), db.Escape(endTime))

	case PointPerson:
		title = "积分提现人数"
		yaxisName = "人数"
		seriesName = "人数"
		query = fmt.Sprintf(`
SELECT
	COUNT(distinct user_id) AS _value,
	DATE(inserted_at) AS record_on,
	0 AS cash
FROM tmm.point_withdraws
WHERE inserted_at > '%s' AND inserted_at < DATE_ADD('%s', INTERVAL 1 DAY) AND verified!=-1
GROUP BY record_on`, db.Escape(startTime), db.Escape(endTime))
	case TotalDrawCash_:
		title = "提现总金额"
		query = fmt.Sprintf(`
SELECT
	IFNULL(t._value, 0) AS _value,
	IFNULL(t.record_on, '%s') AS record_on,
	beforecash.cash AS cash
FROM
(
	SELECT SUM(tmp._value) AS cash
	FROM (
		SELECT SUM(cny) AS _value
		FROM tmm.point_withdraws 
		WHERE inserted_at<'%s' AND verified!=-1
	UNION ALL
		SELECT SUM(cny)  AS _value
		FROM tmm.withdraw_txs
		WHERE verified!=-1 AND tx_status = 1  AND inserted_at<'%s' 
	) AS tmp
) AS beforecash
LEFT JOIN(
    SELECT SUM(cny) AS _value, record_on
    FROM(
        SELECT cny, DATE(inserted_at) AS record_on
        FROM tmm.point_withdraws
        WHERE inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY) AND verified!=-1
        UNION ALL
        SELECT cny, DATE(inserted_at) AS record_on
        FROM tmm.withdraw_txs
        WHERE tx_status=1 AND inserted_at BETWEEN '%s' AND DATE_ADD('%s', INTERVAL 1 DAY)
    ) AS tmp
) AS t ON (1=1)
`, db.Escape(endTime), db.Escape(startTime), db.Escape(startTime), db.Escape(startTime), db.Escape(endTime), db.Escape(startTime), db.Escape(endTime))
	}

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName, valueList []string
	if len(rows) > 0 {

		cash := rows[0].Float(2)

		dataMap := make(map[string]float64)
		for _, row := range rows {
			switch req.Type {
			case TotalDrawCash_:
				cash += row.Float(0)
				dataMap[row.Str(1)] = cash
			default:
				dataMap[row.Str(1)] = row.Float(0)
			}
		}

		format := "%.2f"
		cash = rows[0].Float(2)
		if req.Type == UcPerson || req.Type == PointPerson {
			format = "%.0f"
		}

		for {
			if s.Equal(e) {
				if value, ok := dataMap[s.Format(`2006-01-02`)]; ok {
					indexName = append(indexName, s.Format(`2006-01-02`))
					valueList = append(valueList, fmt.Sprintf(format, value))
				} else {
					indexName = append(indexName, s.Format(`2006-01-02`))
					if req.Type == TotalDrawCash_ {
						valueList = append(valueList, fmt.Sprintf(format, cash))
					} else {
						valueList = append(valueList, fmt.Sprintf(format, 0.0))
					}
				}
				break
			}
			if value, ok := dataMap[s.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, s.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, value))
				cash = value
				s = s.AddDate(0, 0, 1)
			} else {
				indexName = append(indexName, s.Format(`2006-01-02`))
				if req.Type == TotalDrawCash_ {
					valueList = append(valueList, fmt.Sprintf(format, cash))
				} else {
					valueList = append(valueList, fmt.Sprintf(format, 0.0))
				}
				s = s.AddDate(0, 0, 1)
			}
		}
	}
	data.Title.Text = title
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = yaxisName
	data.Series = append(data.Series, Series{Data: valueList, Name: seriesName, Type: "line"})
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})
}
