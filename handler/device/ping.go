package device

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/probab/dst"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"math"
	"net/http"
)

const (
	DEVICE_APP_PING_CACHE_KEY     = "APPPing:%s-%s"
	DEVICE_APP_PING_CAP_CACHE_KEY = "APPPingCap:%s-%s"
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
	if Check(pingRequest.Device.IsEmulator || pingRequest.Device.IsJailBrojen, "invalid request", c) {
		return
	}
	if Check(pingRequest.Ts > 70, "invalid request", c) {
		return
	}

	ping(Service, pingRequest)
	c.JSON(http.StatusOK, APIResponse{Msg: "OK"})
}

func ping(service *common.Service, pingRequest common.PingRequest) {
	device := pingRequest.Device
	deviceId := device.DeviceId()
	appId := device.AppId()

	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	lastPingKey := utils.Sha1(fmt.Sprintf(DEVICE_APP_PING_CACHE_KEY, deviceId, appId))
	lastPingTs, _ := redis.Int64(redisConn.Do("GET", lastPingKey))
	nowTs := time.Now().Unix()
	redisConn.Do("SETEX", lastPingKey, 60, nowTs)
	if nowTs-lastPingTs < 60 {
		log.Warn("Too fast Ping %d, Device:%s, AppId:%s", nowTs-lastPingTs, deviceId, appId)
		return
	}
	pingCapKey := utils.Sha1(fmt.Sprintf(DEVICE_APP_PING_CAP_CACHE_KEY, deviceId, appId))
	pingCap, _ := redis.Int64(redisConn.Do("GET", pingCapKey))
	if pingCap > 120 {
		log.Warn("PingCap: %d, Device:%s, AppId:%s", pingCap, deviceId, appId)
		return
	}
	redisConn.Do("SETEX", lastPingKey, 24*60, pingCap+pingRequest.Ts)

	db := service.Db
	query := `UPDATE
                tmm.devices AS d,
                tmm.device_apps AS da,
                tmm.apps AS a
            SET
                d.total_ts=d.total_ts+%d,
                da.total_ts=da.total_ts+%d,
                a.total_ts=a.total_ts+%d,
                d.tmp_ts=d.tmp_ts+%d,
                d.lastping_at=NOW(),
                da.lastping_at=NOW(),
                a.lastping_at=NOW(),
                a.is_active=1
            WHERE
                da.device_id=d.id
            AND d.id='%s'
            AND da.app_id=a.id AND a.id='%s'
            AND a.platform='%s'`
	_, _, err := db.Query(query, pingRequest.Ts, pingRequest.Ts, pingRequest.Ts, pingRequest.Ts, db.Escape(deviceId), db.Escape(appId), device.Platform)
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
		raven.CaptureError(err, nil)
		log.Error(err.Error())
		return
	}
	floatPoints, _ := d.Points.Float64()
	growthFactor, _ := d.GrowthFactor.Float64()
	growthRate := Config.GrowthRate * growthFactor / (1 + math.Log1p(floatPoints)) * tmpTs
	fn := dst.LogisticPDF(0, 2000.33)
	growthPoints := decimal.NewFromFloat(fn(floatPoints) * 10000 * growthRate)
	log.Warn("points:%s, growthRate: %f, fn:%f, growthPoint:%s", d.Points.String(), growthRate, fn(floatPoints), growthPoints.String())
	_, _, err = db.Query(`UPDATE tmm.devices SET points=points + %s, tmp_ts=0 WHERE id='%s' LIMIT 1`, growthPoints.String(), deviceId)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Error(err.Error())
	}
}
