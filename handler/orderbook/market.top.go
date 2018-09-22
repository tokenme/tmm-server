package orderbook

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/orderbook"
	"net/http"
	"strconv"
)

func MarketTopHandler(c *gin.Context) {
	side, err := strconv.ParseUint(c.Param("side"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	var orderBy string
	if orderbook.Side(side) == orderbook.Ask {
		orderBy = "ORDER BY price ASC"
	} else {
		orderBy = "ORDER BY price DESC"
	}
	db := Service.Db
	query := `SELECT side, quantity, price, deal_quantity, deal_eth FROM tmm.orderbooks WHERE online_status=0 AND side=%d %s LIMIT 10`
	rows, _, err := db.Query(query, side, orderBy)
	if CheckErr(err, c) {
		return
	}
	var orders []common.Order
	for _, row := range rows {
		quantity, _ := decimal.NewFromString(row.Str(1))
		price, _ := decimal.NewFromString(row.Str(2))
		dealQuantity, _ := decimal.NewFromString(row.Str(3))
		dealEth, _ := decimal.NewFromString(row.Str(4))
		o := common.Order{
			Side:         orderbook.Side(row.Uint(0)),
			Quantity:     quantity,
			Price:        price,
			DealQuantity: dealQuantity,
			DealEth:      dealEth,
		}
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}
