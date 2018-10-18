package bonus

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type DailyCommitRequest struct {
	Idfa string `json:"idfa" form:"idfa"`
}

func DailyCommitHandler(c *gin.Context) {
	var req DailyCommitRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	deviceRequest := common.DeviceRequest{
		Platform: common.IOS,
		Idfa:     req.Idfa,
	}
	deviceId := deviceRequest.DeviceId()

	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	_, ret, err := db.Query(`INSERT INTO tmm.daily_bonus_logs (user_id, updated_on, days) VALUES (%d, NOW(), 1) ON DUPLICATE KEY UPDATE days=IF(updated_on=DATE_SUB(NOW(), INTERVAL 1 DAY), days+1, IF(updated_on=DATE(NOW()), days, 1)), updated_on=VALUES(updated_on)`, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(ret.AffectedRows() == 0, DAILY_BONUS_COMMITTED_ERROR, "Already checked in", c) {
		return
	}
	rows, _, err := db.Query(`SELECT days FROM tmm.daily_bonus_logs WHERE user_id=%d LIMIT 1`, user.Id)
	if CheckErr(err, c) {
		return
	}
	days := rows[0].Int64(0)
	points := decimal.New(days, 0)
	pointsPerTs, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := points.Div(pointsPerTs)
	_, _, err = db.Query(`UPDATE tmm.devices SET points=points+%s, total_ts=total_ts+%s WHERE id='%s' AND user_id=%d`, points.String(), ts.String(), db.Escape(deviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"days": days, "points": points})
}
