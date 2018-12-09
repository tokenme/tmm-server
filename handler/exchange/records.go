package exchange

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/coins/eth/utils"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

const DEFAULT_PAGE_SIZE = 10

type RecordType = uint8

const (
	RedeemCdpRecordType RecordType = 1
)

type RecordsRequest struct {
	Page      uint                        `json:"page" form:"page"`
	PageSize  uint                        `json:"page_size" form:"page_size"`
	Direction common.TMMExchangeDirection `json:"direction" form:"direction"`
	Type      RecordType                  `json:"type" form:"type"`
}

func RecordsHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req RecordsRequest
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
	if req.Type == RedeemCdpRecordType {
		query := `SELECT
        co.device_id, co.points, co.grade, co.inserted_at
        FROM tmm.cdp_orders AS co
        INNER JOIN tmm.devices AS d ON (d.id=co.device_id)
        WHERE d.user_id=%d ORDER BY co.inserted_at DESC LIMIT %d, %d`
		rows, _, err := db.Query(query, user.Id, (req.Page-1)*req.PageSize, req.PageSize)
		if CheckErr(err, c) {
			return
		}
		var records []common.RedeemCdpRecord
		for _, row := range rows {
			points, _ := decimal.NewFromString(row.Str(1))
			record := common.RedeemCdpRecord{
				DeviceId:   row.Str(0),
				Points:     points,
				Grade:      row.Str(2),
				InsertedAt: row.ForceLocaltime(3).Format(time.RFC3339),
			}
			records = append(records, record)
		}
		c.JSON(http.StatusOK, records)
		return
	}
	pointsPerTs, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	query := `SELECT
        tx, status, device_id, tmm, points, direction, inserted_at
        FROM tmm.exchange_records
        WHERE user_id=%d AND direction=%d ORDER BY inserted_at DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, user.Id, req.Direction, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var records []common.ExchangeRecord
	for _, row := range rows {
		tmm, _ := decimal.NewFromString(row.Str(3))
		points, _ := decimal.NewFromString(row.Str(4))
		record := common.ExchangeRecord{
			Tx:         row.Str(0),
			Status:     common.ExchangeTxStatus(row.Uint(1)),
			DeviceId:   row.Str(2),
			Tmm:        tmm,
			Points:     points,
			Direction:  common.TMMExchangeDirection(row.Int(5)),
			InsertedAt: row.ForceLocaltime(6).Format(time.RFC3339),
		}
		if record.Status == common.ExchangeTxPending {
			receipt, err := utils.TransactionReceipt(Service.Geth, c, record.Tx)
			if err == nil {
				record.Status = common.ExchangeTxStatus(receipt.Status)
				if record.Status == common.ExchangeTxFailed && record.Direction == common.TMMExchangeIn {
					_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.exchange_records AS er SET d.points=d.points + er.points, d.total_ts = CEIL(d.total_ts + %s), er.status=0 WHERE d.id=er.device_id AND er.tx='%s'`, pointsPerTs.String(), db.Escape(record.Tx))
					if err != nil {
						log.Error(err.Error())
					}
				} else if record.Status == common.ExchangeTxFailed && record.Direction == common.TMMExchangeOut {
					_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.exchange_records AS er SET d.points=IF(d.points > er.points, d.points - er.points, 0), d.consumed_ts = CEIL(d.consumed_ts + %s), er.status=0 WHERE d.id=er.device_id AND er.tx='%s'`, pointsPerTs.String(), db.Escape(record.Tx))
					if err != nil {
						log.Error(err.Error())
					}
				} else {
					_, _, err := db.Query(`UPDATE tmm.exchange_records SET status=%d WHERE tx='%s'`, receipt.Status, db.Escape(record.Tx))
					if err != nil {
						log.Error(err.Error())
					}
				}
			}
		}
		records = append(records, record)
	}
	c.JSON(http.StatusOK, records)
}
