package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
)

func InvestWithdrawHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	itemId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT points, device_id FROM tmm.good_invests WHERE good_id=%d AND user_id=%d AND redeem_status=0 LIMIT 1`, itemId, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}

	points, _ := decimal.NewFromString(rows[0].Str(0))
	deviceId := rows[0].Str(1)
	tspoints, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := points.Div(tspoints)
	_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.good_invests AS gi SET d.points=d.points+gi.points, d.total_ts=d.total_ts+%d, gi.redeem_status=2 WHERE d.id=gi.device_id AND gi.redeem_status=0 AND gi.good_id=%d AND gi.user_id=%d AND d.id='%s'`, ts.IntPart(), itemId, user.Id, db.Escape(deviceId))
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
