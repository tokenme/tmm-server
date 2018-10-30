package bonus

import (
	"encoding/json"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"time"
)

type ReadingRequest struct {
	APPKey  string `form:"key" json:"key" binding:"required"`
	Idfa    string `json:"idfa" form:"idfa"`
	Imei    string `json:"imei" form:"imei"`
	Mac     string `json:"mac" form:"mac"`
	Payload string `json:"payload" form:"payload"  binding:"required"`
}

type ReadingBonus struct {
	TaskId   uint64          `json:"task_id" form:"task_id"`
	Points   decimal.Decimal `json:"points" form:"points"`
	Duration int64           `json:"duration" form:"duration"`
	Ts       int64           `json:"ts" form:"ts"`
}

func ReadingHandler(c *gin.Context) {
	var req ReadingRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	secret := GetAppSecret(req.APPKey)
	if Check(secret == "", "empty secret", c) {
		log.Error("empty secret")
		return
	}
	decrepted, err := utils.DesDecrypt(req.Payload, []byte(secret))
	if CheckErr(err, c) {
		return
	}
	var payload ReadingBonus
	err = json.Unmarshal(decrepted, &payload)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	if Check(payload.Ts < time.Now().Add(-10*time.Minute).Unix() || payload.Ts > time.Now().Add(10*time.Minute).Unix(), "expired request", c) {
		return
	}
	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.share_tasks WHERE id=%d LIMIT 1`, payload.TaskId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	_, _, err = db.Query(`UPDATE tmm.devices SET total_ts=total_ts+%d, points=points+%s WHERE id='%s' AND user_id=%d`, payload.Duration, payload.Points.String(), db.Escape(deviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, payload)
}
