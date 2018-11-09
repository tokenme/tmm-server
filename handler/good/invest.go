package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ykt"
	"net/http"
)

type InvestRequest struct {
	Idfa   string          `json:"idfa" form:"idfa"`
	Imei   string          `json:"imei" form:"imei"`
	Mac    string          `json:"mac" form:"mac"`
	GoodId uint64          `json:"good_id" form:"good_id" binding:"required"`
	Points decimal.Decimal `json:"points" form:"points" binding:"required"`
}

func InvestHandler(c *gin.Context) {
	var req InvestRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if CheckWithCode(len(deviceId) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}

	yktReq := ykt.GoodInfoRequest{
		Id:  req.GoodId,
		Uid: user.Id,
	}
	res, err := yktReq.Run()
	if CheckErr(err, c) {
		return
	}
	good := res.Data.Data

	points, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := req.Points.Div(points)
	db := Service.Db
	_, ret, err := db.Query(`UPDATE tmm.devices SET points=points-%s, consumed_ts=consumed_ts+%d WHERE id='%s' AND user_id=%d AND points>=%s`, req.Points.String(), ts.IntPart(), db.Escape(deviceId), user.Id, req.Points.String())
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(ret.AffectedRows() == 0, NOT_ENOUGH_POINTS_ERROR, "not enough points", c) {
		return
	}
	_, _, err = db.Query(`INSERT INTO tmm.good_invests (good_id, user_id, device_id, points) VALUES (%d, %d, '%s', %s) ON DUPLICATE KEY UPDATE points=points+VALUES(points)`, req.GoodId, user.Id, db.Escape(deviceId), req.Points.String())
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`INSERT IGNORE INTO tmm.goods (id, name, pic) VALUES (%d, '%s', '%s')`, req.GoodId, db.Escape(good.Name), db.Escape(good.Pic))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
