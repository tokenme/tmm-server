package redeem

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

const DEFAULT_PAGE_SIZE = 10

type TMMWithdrawListRequest struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

func TMMWithdrawListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req TMMWithdrawListRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}
	db := Service.Db
	query := `SELECT
        tx, tx_status, tmm, cny, withdraw_status, inserted_at
        FROM tmm.withdraw_txs
        WHERE user_id=%d ORDER BY inserted_at DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var records []common.TMMWithdrawRecord
	for _, row := range rows {
		tmm, _ := decimal.NewFromString(row.Str(2))
		cny, _ := decimal.NewFromString(row.Str(3))
		record := common.TMMWithdrawRecord{
			Tx:             row.Str(0),
			TxStatus:       row.Uint(1),
			TMM:            tmm,
			Cash:           cny,
			WithdrawStatus: row.Uint(4),
			InsertedAt:     row.ForceLocaltime(5).Format(time.RFC3339),
		}
		records = append(records, record)
	}
	c.JSON(http.StatusOK, records)
	return
}
