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
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
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

	db := Service.Db
	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	if Check(payload.Duration <= 5, "read too fast", c) {
		_, _, err = db.Query(`INSERT INTO tmm.reading_logs (user_id, task_id, ts, point) VALUES (%d, %d, %d ,0) ON DUPLICATE KEY UPDATE ts=ts+VALUES(ts)`, user.Id, payload.TaskId, payload.Duration)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	if Check(payload.Duration >= 480, "read too slow", c) {
		_, _, err = db.Query(`INSERT INTO tmm.reading_logs (user_id, task_id, ts, point) VALUES (%d, %d, %d ,0) ON DUPLICATE KEY UPDATE ts=ts+VALUES(ts)`, user.Id, payload.TaskId, payload.Duration)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	rows, _, err := db.Query(`SELECT st.id, rl.ts FROM tmm.share_tasks AS st LEFT JOIN tmm.reading_logs AS rl ON (rl.task_id=st.id AND rl.user_id=%d) WHERE st.id=%d LIMIT 1`, user.Id, payload.TaskId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}

	readTs := rows[0].Int64(1)
	if Check(readTs+payload.Duration >= 480, "read too much", c) {
		return
	}

	maxPoints := decimal.New(1, 2)
	if payload.Points.GreaterThan(maxPoints) {
		payload.Points = maxPoints
	}
	pointsPerTs, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	ts := payload.Points.Div(pointsPerTs)
	_, _, err = db.Query(`UPDATE tmm.devices SET total_ts=total_ts+%d, points=points+%s WHERE id='%s' AND user_id=%d`, ts.IntPart(), payload.Points.String(), db.Escape(deviceId), user.Id)
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`INSERT INTO tmm.reading_logs (user_id, task_id, ts, point) VALUES (%d, %d, %d ,%s) ON DUPLICATE KEY UPDATE ts=ts+VALUES(ts),point=point+VALUES(point)`, user.Id, payload.TaskId, payload.Duration, payload.Points.String())
	if err != nil {
		log.Error(err.Error())
	}
	c.JSON(http.StatusOK, payload)
}
