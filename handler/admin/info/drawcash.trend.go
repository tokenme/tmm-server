package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"time"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
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
	db := Service.Db
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
	query := `
SELECT
	tmp._value AS _value,
	tmp.date AS date,
	tmp.cash AS cash
FROM(
%s
) AS tmp
GROUP BY tmp.date
ORDER BY tmp.date 
`
	var data Data
	var title string
	seriesName := "金额"
	yaxisName := "金额"
	switch req.Type {
	case DrawCashByUc:
		title = "UC提现金额趋势图"
		query = fmt.Sprintf(query, fmt.Sprintf(`
SELECT  
	SUM(cny) AS _value,
	DATE_FORMAT(inserted_at,'%s') AS date,
	0 AS cash
FROM 
	tmm.withdraw_txs 
WHERE 
	tx_status = 1 AND inserted_at > '%s'  AND inserted_at < '%s' 
GROUP BY date`, TimeFormat, db.Escape(startTime), db.Escape(endTime)))
	case UcPerson:
		title = "UC提现人数趋势图"
		yaxisName = "人数"
		seriesName = "人数"
		query = fmt.Sprintf(query, fmt.Sprintf(`
SELECT  
	COUNT(distinct user_id) AS _value,
	DATE_FORMAT(inserted_at,'%s') AS date,
	0 AS cash
FROM 
	tmm.withdraw_txs 
WHERE 
	tx_status = 1 AND inserted_at > '%s'  AND inserted_at < '%s' 
GROUP BY date`, TimeFormat, db.Escape(startTime), db.Escape(endTime)))
	case DrawCashByPoint:
		title = "积分提现金额趋势图"
		query = fmt.Sprintf(query, fmt.Sprintf(`
SELECT  
	SUM(cny) AS _value,
	DATE_FORMAT(inserted_at,'%s') AS date,
	0 AS cash
FROM 
	tmm.point_withdraws 
WHERE 
	 inserted_at > '%s' AND inserted_at < '%s' 
GROUP BY date`, TimeFormat, db.Escape(startTime), db.Escape(endTime)))

	case PointPerson:
		title = "积分提现人数趋势图"
		yaxisName = "人数"
		seriesName = "人数"
		query = fmt.Sprintf(query, fmt.Sprintf(`
SELECT  
	COUNT(distinct user_id) AS _value,
	DATE_FORMAT(inserted_at,'%s') AS date,
	0 AS cash
FROM 
	tmm.point_withdraws 
WHERE 
	 inserted_at > '%s' AND inserted_at < '%s'
GROUP BY date`, TimeFormat, db.Escape(startTime), db.Escape(endTime)))
	case TotalDrawCash_:
		title = "提现总金额趋势图"
		query = fmt.Sprintf(query, fmt.Sprintf(`
SELECT 
	IFNULL(tmp._value,0) AS _value,
	IFNULL(tmp.date,'%s') AS date,
	beforecash.cash AS cash
FROM
(
	SELECT
		SUM(tmp._value) AS cash
	FROM (
		SELECT
			SUM(cny) AS _value
		FROM
			tmm.point_withdraws
		WHERE
	 		inserted_at < '%s' 
	UNION ALL
		SELECT
			SUM(cny)  AS _value
		FROM
			tmm.withdraw_txs
		WHERE
			tx_status = 1 AND inserted_at < '%s'  
		) AS tmp
) AS beforecash
LEFT JOIN(
SELECT
	SUM(tmp._value) AS _value,
	tmp.date AS date
FROM(
	SELECT
		SUM(cny) AS _value,
		DATE_FORMAT(inserted_at,'%s') AS date
	FROM
		tmm.point_withdraws
	WHERE
		inserted_at > '%s' AND inserted_at < '%s' 
	GROUP BY date
	UNION ALL
	SELECT
		SUM(cny)  AS _value,
		DATE_FORMAT(inserted_at,'%s') AS date
	FROM
		tmm.withdraw_txs
	WHERE
		tx_status = 1 AND inserted_at > '%s' AND inserted_at < '%s' 
	GROUP BY 
		date
	) AS tmp
	GROUP BY 
		date
) AS tmp  ON 1 = 1
`, db.Escape(endTime), db.Escape(startTime), db.Escape(startTime),
			TimeFormat, db.Escape(startTime), db.Escape(endTime),
			TimeFormat, db.Escape(startTime), db.Escape(endTime)))
	}
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var indexName, valueList []string
	format := "%.2f"
	if req.Type == UcPerson || req.Type == PointPerson {
		format = "%.0f"
	}
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)
	var cash float64
	if req.Type == TotalDrawCash_{
		cash = rows[0].Float(2)
	}
	dataMap := make(map[string]float64)
	for _, row := range rows {
		if req.Type == TotalDrawCash_ {
			cash += row.Float(0)
			dataMap[row.Str(1)] = cash
		} else {
			dataMap[row.Str(1)] = row.Float(0)
		}
	}
	if req.Type == TotalDrawCash_{
		cash = rows[0].Float(2)
	}
	for {
		if tm.Equal(end) {
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf(format, value))
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				if req.Type == TotalDrawCash_ {
					valueList = append(valueList, fmt.Sprintf(format, cash))
				} else {
					valueList = append(valueList, fmt.Sprintf(format, 0.0))
				}
			}
			break
		}
		if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			valueList = append(valueList, fmt.Sprintf(format, value))
			cash = value
			tm = tm.AddDate(0,0,1)
		} else {
			indexName = append(indexName, tm.Format(`2006-01-02`))
			if req.Type == TotalDrawCash_ {
				valueList = append(valueList, fmt.Sprintf(format, cash))
			} else {
				valueList = append(valueList, fmt.Sprintf(format, 0.0))
			}
			tm = tm.AddDate(0,0,1)
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
