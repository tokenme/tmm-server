package orderbook

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/orderbook"
	"net/http"
	"strconv"
	"time"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

func OrdersHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	sideInt, _ := strconv.ParseUint(c.Param("side"), 10, 64)
	side := orderbook.Side(sideInt)
	if Check(side != orderbook.Ask && side != orderbook.Bid, "invalid request", c) {
		return
	}
	page, _ := strconv.ParseUint(c.Param("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(c.Param("pageSize"), 10, 64)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 || pageSize > DEFAULT_PAGE_SIZE {
		pageSize = DEFAULT_PAGE_SIZE
	}

	db := Service.Db
	query := `SELECT
    id,
    side,
    quantity,
    price,
    deal_quantity,
    deal_eth,
    online_status,
    inserted_at,
    updated_at
FROM tmm.orderbooks
WHERE user_id=%d AND side=%d ORDER BY id DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, side, (page-1)*pageSize, pageSize)
	if CheckErr(err, c) {
		return
	}
	var orders []common.Order
	for _, row := range rows {
		quantity, _ := decimal.NewFromString(row.Str(2))
		price, _ := decimal.NewFromString(row.Str(3))
		dealQuantity, _ := decimal.NewFromString(row.Str(4))
		dealEth, _ := decimal.NewFromString(row.Str(5))
		o := common.Order{
			TradeId:      row.Uint64(0),
			Side:         orderbook.Side(row.Uint(1)),
			Quantity:     quantity,
			Price:        price,
			DealQuantity: dealQuantity,
			DealEth:      dealEth,
			OnlineStatus: int8(row.Int(6)),
			InsertedAt:   row.ForceLocaltime(7).Format(time.RFC3339),
		}
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}
