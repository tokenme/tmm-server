package redeem

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

type PointsWithdrawListRequest struct {
	DeviceId string `json:"device_id" from:"device_id"`
	Page     uint   `json:"page" form:"page"`
	PageSize uint   `json:"page_size" form:"page_size"`
}

func PointsWithdrawListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req PointsWithdrawListRequest
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
	var where string
	if req.DeviceId != "" {
		where = fmt.Sprintf(" AND device_id='%s'", db.Escape(req.DeviceId))
	}
	query := `SELECT
        trade_num, points, cny, inserted_at
        FROM tmm.point_withdraws
        WHERE user_id=%d%s ORDER BY inserted_at DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, where, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var records []common.TMMWithdrawRecord
	for _, row := range rows {
		points, _ := decimal.NewFromString(row.Str(1))
		cny, _ := decimal.NewFromString(row.Str(2))
		record := common.TMMWithdrawRecord{
			Tx:         row.Str(0),
			TMM:        points,
			Cash:       cny,
			InsertedAt: row.ForceLocaltime(3).Format(time.RFC3339),
		}
		records = append(records, record)
	}
	c.JSON(http.StatusOK, records)
	return
}
