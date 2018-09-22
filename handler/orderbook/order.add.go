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
)

type OrderAddRequest struct {
	ProcessType orderbook.ProcessType `json:"process_type" form:"process_type" binding:"required"`
	Side        orderbook.Side        `json:"side" form:"side" binding:"required"`
	Quantity    decimal.Decimal       `json:"quantity" form:"quantity" binding:"required"`
	Price       decimal.Decimal       `json:"price" form:"price" binding:"required"`
}

func OrderAddHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req OrderAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	db := Service.Db
	query := `INSERT INTO tmm.orderbooks (user_id, side, process_type, quantity, price) VALUES (%d, %d, %d, %s, %s)`
	_, _, err := db.Query(query, user.Id, req.Side, req.ProcessType, req.Quantity, req.Price)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
