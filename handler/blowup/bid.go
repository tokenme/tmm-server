package blowup

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type BidRequest struct {
	Points    decimal.Decimal `json:"points" form:"points" required:"true"`
	SessionId uint64          `json:"session_id" form:"session_id" required:"true"`
	Idfa      string          `json:"idfa" form:"idfa"`
}

func BidHandler(c *gin.Context) {
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

	pointsTsRate, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := req.Points.Div(pointsTsRate)

	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.blowup_sessions WHERE id=%d LIMIT 1`, req.SessionId)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	_, ret, err := db.Query(`INSERT IGNORE INTO tmm.blowup_bids (user_id, session_id, points, ts) VALUES (%d, %d, %s, %s)`, user.Id, req.SessionId, req.Points.String(), ts.String())
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() > 0 {
		_, _, err := db.Query(`UPDATE tmm.devices AS d, tmm.blowup_bids AS bb SET d.points = d.points-bb.points, d.consumed_ts=d.consumed_ts+bb.ts WHERE d.id='%s' AND d.user_id=%d AND bb.user_id=d.user_id AND bb.session_id=%d AND d.points>=bb.points`, db.Escape(deviceId), user.Id, req.SessionId)
		if CheckErr(err, c) {
			_, _, err = db.Query(`DELETE FROM tmm.blowup_bids WHERE user_id=%d AND session_id=%d`, user.Id, req.SessionId)
			if err != nil {
				log.Error(err.Error())
			}
			return
		}
	}

	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
