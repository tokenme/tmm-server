package orderbook

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func RateHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}

	db := Service.Db
	query := `SELECT 0, price FROM
(SELECT price FROM tmm.orderbook_trades ORDER BY id DESC LIMIT 1) AS t1
UNION
SELECT 1, price FROM
(SELECT price FROM tmm.orderbook_trades WHERE inserted_at<=DATE_SUB(NOW(), INTERVAL 1 DAY) ORDER BY inserted_at DESC LIMIT 1) AS t2`
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	var (
		currentPrice decimal.Decimal
		prevPrice    decimal.Decimal
		changeRate   decimal.Decimal
	)
	for _, row := range rows {
		id := row.Uint(0)
		if id == 0 {
			currentPrice, _ = decimal.NewFromString(row.Str(1))
		} else {
			prevPrice, _ = decimal.NewFromString(row.Str(1))
		}
	}
	if prevPrice.GreaterThan(decimal.Zero) {
		changeRate = currentPrice.Sub(prevPrice).Div(prevPrice)
	}
	c.JSON(http.StatusOK, gin.H{"price": currentPrice, "change_rate": changeRate})
}
