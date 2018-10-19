package blowup

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/blowup"
	"net/http"
)

type EscapeRequest struct {
	SessionId uint64 `json:"session_id" form:"session_id" required:"true"`
	Idfa      string `json:"idfa" form:"idfa"`
}

func EscapeHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req BidRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	device := common.DeviceRequest{
		Idfa:     req.Idfa,
		Platform: common.IOS,
	}
	deviceId := device.DeviceId()

	db := Service.Db
	_, ret, err := db.Query(`UPDATE tmm.devices AS d, tmm.blowup_bids AS bb, tmm.blowup_sessions AS bs SET bb.rate=NOW()-created_at, bb.escaped=1, d.points = d.points+bb.points * (1 + (NOW()-created_at)/100), d.total_ts=d.total_ts+bb.ts*(1 + (NOW()-created_at)/100) WHERE d.id='%s' AND d.user_id=%d AND bb.user_id=d.user_id AND bb.session_id=%d AND bs.id=bb.session_id AND bb.escaped=0 AND bb.inserted_at<=DATE_SUB(NOW(), INTERVAL 30 SECOND)`, db.Escape(deviceId), user.Id, req.SessionId)
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() == 0 {
		rows, _, err := db.Query(`SELECT 1 FROM tmm.blowup_bids WHERE user_id=%d AND session_id=%d AND inserted_at<=DATE_SUB(NOW(), INTERVAL 30 SECOND) LIMIT 1`, user.Id, req.SessionId)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, BLOWUP_ESCAPE_EARLY_ERROR, "escape early", c) {
			return
		} else if CheckWithCode(true, BLOWUP_ESCAPE_LATE_ERROR, "escape late", c) {
			return
		}
	}

	rows, _, err := db.Query(`SELECT points, rate/100 FROM tmm.blowup_bids WHERE user_id=%d AND session_id=%d LIMIT 1`, user.Id, req.SessionId)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	row := rows[0]
	points, _ := decimal.NewFromString(row.Str(0))
	rate, _ := decimal.NewFromString(row.Str(1))
	BlowupService.NewEscape(req.SessionId, points, rate)
	c.JSON(http.StatusOK, blowup.Event{
		Type:      blowup.EscapeEvent,
		SessionId: req.SessionId,
		Value:     points,
		Rate:      rate,
		Nick:      user.ShowName,
	})
}
