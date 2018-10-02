package orderbook

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"math"
	"net/http"
	"strconv"
	"time"
)

const DefaultHours time.Duration = 24

func MarketGraphHandler(c *gin.Context) {
	reqHours, err := strconv.ParseInt(c.Param("hours"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	hours := time.Duration(reqHours) * time.Hour
	if reqHours == 0 {
		hours = DefaultHours
	}
	db := Service.Db
	var (
		interval string
		query    string
		week     = time.Duration(24 * 7 * time.Hour)
	)
	if hours <= week {
		interval = fmt.Sprintf("DATE_SUB(NOW(), INTERVAL %d HOUR)", hours)
		query = `SELECT COUNT(*) AS trades, SUM(quantity) AS quantity, AVG(price) AS price, MAX(price) AS high, MIN(price) AS low, DATE(updated_at) AS d, HOUR(updated_at) AS h  FROM tmm.orderbook_trades WHERE tx_status=1 AND side=1 AND updated_at>=%s GROUP BY d, h ORDER BY d, h`
	} else {
		interval = fmt.Sprintf("DATE_SUB(NOW(), INTERVAL %d DAY)", uint(math.Ceil(hours.Hours()/24.0)))
		query = `SELECT COUNT(*) AS trades, SUM(quantity) AS quantity, AVG(price) AS price, MAX(price) AS high, MIN(price) AS low, DATE(updated_at) AS d  FROM tmm.orderbook_trades WHERE tx_status=1 AND side=1 AND updated_at>=%s GROUP BY d ORDER BY d`
	}

	rows, _, err := db.Query(query, interval)
	if CheckErr(err, c) {
		return
	}
	var graph []common.MarketGraph
	mp := make(map[string]common.MarketGraph, 24)
	for _, row := range rows {
		quantity, _ := decimal.NewFromString(row.Str(1))
		price, _ := decimal.NewFromString(row.Str(2))
		low, _ := decimal.NewFromString(row.Str(3))
		high, _ := decimal.NewFromString(row.Str(4))
		var at string
		if len(row) == 7 {
			at = fmt.Sprintf("%sT%d:00:00Z00:00 ", row.Str(5), row.Uint(6))
		} else {
			at = row.ForceLocaltime(5).Format(time.RFC3339)
		}

		d := common.MarketGraph{
			Trades:   row.Uint64(0),
			Quantity: quantity,
			Price:    price,
			Low:      low,
			High:     high,
			At:       at,
		}
		mp[at] = d
	}
	now := time.Now()
	loc := now.Location()
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	step := time.Duration(1 * time.Hour)
	dateFormat := "2006-01-02T15:00:00Z07:00"
	if hours > week {
		step = time.Duration(24 * time.Hour)
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		dateFormat = "2006-01-02T00:00:00Z07:00"
	}
	var h = hours
	for h >= 0 {
		at := endTime.Add(-1 * h).Format(dateFormat)
		h -= step
		if d, found := mp[at]; found {
			graph = append(graph, d)
			continue
		}
		graph = append(graph, common.MarketGraph{
			At: at,
		})
	}
	c.JSON(http.StatusOK, graph)
}
