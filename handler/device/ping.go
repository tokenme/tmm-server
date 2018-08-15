package device

import (
	"encoding/json"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"math"
	"net/http"
)

func PingHandler(c *gin.Context) {
	req := c.MustGet("Request").(APIRequest)
	secret := c.MustGet("Secret").(string)
	decrepted, err := utils.DesDecrypt(req.Payload, []byte(secret))
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	var pingRequest common.PingRequest
	err = json.Unmarshal(decrepted, &pingRequest)
	if CheckErr(err, c) {
		raven.CaptureError(err, nil)
		return
	}
	//if Check(pingRequest.Device.IsEmulator || pingRequest.Device.IsJailBrojen, "invalid request", c) {
	//  return
	//}
	ping(Service, pingRequest)
	c.JSON(http.StatusOK, APIResponse{Msg: "OK"})
}

func ping(service *common.Service, pingRequest common.PingRequest) {
	device := pingRequest.Device
	db := service.Db
	deviceId := device.DeviceId()
	appId := device.AppId()
	_, _, err := db.Query(`UPDATE tmm.device_apps SET total_ts=total_ts+%d, lastping_at=NOW() WHERE device_id='%s' AND app_id='%s' LIMIT 1`, pingRequest.Ts, db.Escape(deviceId), db.Escape(appId))
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
	}
	_, _, err = db.Query(`UPDATE tmm.devices SET total_ts=total_ts+%d, tmp_ts=tmp_ts+%d, lastping_at=NOW() WHERE id='%s' LIMIT 1`, pingRequest.Ts, pingRequest.Ts, db.Escape(deviceId))
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
	}
	_, _, err = db.Query(`UPDATE tmm.apps SET total_ts=total_ts+%d, is_active=1, lastping_at=NOW() WHERE id='%s' AND platform='%s' LIMIT 1`, pingRequest.Ts, db.Escape(appId), device.Platform)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
	}
	rows, _, err := db.Query(`SELECT d.tmp_ts, SUM(a.tmm), d.points FROM tmm.devices AS d INNER JOIN tmm.device_apps AS da ON (da.device_id=d.id) INNER JOIN tmm.apps AS a ON (a.id=da.app_id) WHERE d.id='%s' LIMIT 1`, db.Escape(deviceId))
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
		return
	}
	if len(rows) == 0 {
		return
	}
	row := rows[0]
	tmpTs := row.ForceFloat(0)
	if tmpTs < float64(Config.MinGrowthTS) {
		return
	}
	tmmBalance, _ := decimal.NewFromString(row.Str(1))
	points, _ := decimal.NewFromString(row.Str(2))
	d := common.Device{
		Id:         deviceId,
		TMMBalance: tmmBalance,
		Points:     points,
	}
	d.GrowthFactor, err = d.GetGrowthFactor(service)
	if err != nil {
		return
	}
	floatPoints, _ := d.Points.Float64()
	growthFactor, _ := d.GrowthFactor.Float64()
	growthRate := Config.GrowthRate * growthFactor / (1 + math.Log1p(floatPoints)) * tmpTs
	growthPoints := d.Points.Mul(decimal.NewFromFloat(1 + growthRate))
	_, _, err = db.Query(`UPDATE tmm.devices SET points=%s, tmp_ts=0 WHERE id='%s' LIMIT 1`, growthPoints.String(), deviceId)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
	}
}
