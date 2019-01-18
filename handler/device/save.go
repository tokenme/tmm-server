package device

import (
	"encoding/json"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	commonutils "github.com/tokenme/tmm/utils"
	"net/http"
	"strings"
)

func SaveHandler(c *gin.Context) {
	req := c.MustGet("Request").(APIRequest)
	secret := c.MustGet("Secret").(string)
	decrepted, err := commonutils.DesDecrypt(req.Payload, []byte(secret))
	if CheckErr(err, c) {
		log.Error(err.Error())
		raven.CaptureError(err, nil)
		return
	}
	var deviceRequest common.DeviceRequest
	err = json.Unmarshal(decrepted, &deviceRequest)
	if err != nil {
		deviceRequest = common.UnmarshalDeviceRequest(decrepted)
	}
	//if Check(deviceRequest.IsEmulator || deviceRequest.IsJailBrojen, "invalid request", c) {
	//	return
	//}
	err = saveDevice(Service, deviceRequest, c)
	if CheckErr(err, c) {
		return
	}
	err = saveApp(Service, deviceRequest)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "OK"})
}

func saveDevice(service *common.Service, deviceRequest common.DeviceRequest, c *gin.Context) error {
	deviceId := deviceRequest.DeviceId()
	if deviceId == "" {
		return nil
	}
	db := service.Db
	rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE id='%s' LIMIT 1`, db.Escape(deviceId))
	if err != nil {
		raven.CaptureError(err, nil)
		return err
	}
	if len(rows) > 0 {
		return nil
	}
	query := `INSERT INTO tmm.devices (id, platform, idfa, imei, mac, device_name, system_version, os_version, language, model, timezone, country, is_emulator, is_jailbrojen, is_tablet, lastping_at) VALUES ('%s', '%s', %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %d, %d, %d, NOW()) ON DUPLICATE KEY UPDATE idfa=VALUES(idfa), imei=VALUES(imei), mac=VALUES(mac), device_name=VALUES(device_name), system_version=VALUES(system_version), os_version=VALUES(os_version), language=VALUES(language), model=VALUES(model), timezone=VALUES(timezone), country=VALUES(country), lastping_at=VALUES(lastping_at)`
	var (
		deviceName    = "NULL"
		systemVersion = "NULL"
		osVersion     = "NULL"
		language      = "NULL"
		model         = "NULL"
		timezone      = "NULL"
		country       = "NULL"
		idfa          = "NULL"
		imei          = "NULL"
		mac           = "NULL"
		isEmulator    = 0
		isJailBrojen  = 0
		isTablet      = 0
	)
	if deviceRequest.DeviceName != "" {
		deviceName = fmt.Sprintf("'%s'", db.Escape(deviceRequest.DeviceName))
	}
	if deviceRequest.SystemVersion != "" {
		systemVersion = fmt.Sprintf("'%s'", db.Escape(deviceRequest.SystemVersion))
	}
	if deviceRequest.OsVersion != "" {
		osVersion = fmt.Sprintf("'%s'", db.Escape(deviceRequest.OsVersion))
	}
	if deviceRequest.Language != "" {
		language = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Language))
	}
	if deviceRequest.DeviceType != "" {
		model = fmt.Sprintf("'%s'", db.Escape(deviceRequest.DeviceType))
	}
	if deviceRequest.Timezone != "" {
		timezone = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Timezone))
	}
	if deviceRequest.Country != "" {
		country = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Country))
	}
	if deviceRequest.Idfa != "" {
		idfa = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Idfa))
	}
	if deviceRequest.Imei != "" {
		imei = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Imei))
	}
	if deviceRequest.Mac != "" {
		mac = fmt.Sprintf("'%s'", db.Escape(deviceRequest.Mac))
		mac = strings.Replace(mac, ":", "", -1)
	}
	if deviceRequest.IsEmulator {
		isEmulator = 1
	}
	if deviceRequest.IsJailBrojen {
		isJailBrojen = 1
	}
	if deviceRequest.IsTablet {
		isTablet = 1
	}
	_, _, err = db.Query(query, deviceId, deviceRequest.Platform, idfa, imei, mac, deviceName, systemVersion, osVersion, language, model, timezone, country, isEmulator, isJailBrojen, isTablet)
	if err != nil {
		raven.CaptureError(err, nil)
		return err
	}
	return nil
}

func saveApp(service *common.Service, deviceRequest common.DeviceRequest) error {
	deviceId := deviceRequest.DeviceId()
	appId := deviceRequest.AppId()
	if deviceId == "" || appId == "" {
		return nil
	}
	db := service.Db
	query := `INSERT INTO tmm.apps (id, platform, bundle_id, name, app_version, build_version, lastping_at) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', NOW()) ON DUPLICATE KEY UPDATE name=VALUES(name), app_version=VALUES(app_version), build_version=VALUES(build_version), lastping_at=VALUES(lastping_at)`
	_, _, err := db.Query(query, appId, deviceRequest.Platform, db.Escape(deviceRequest.AppBundleId), db.Escape(deviceRequest.AppName), db.Escape(deviceRequest.AppVersion), db.Escape(deviceRequest.AppBuildVersion))
	if err != nil {
		raven.CaptureError(err, nil)
		return err
	}
	_, _, err = db.Query(`INSERT INTO device_apps (device_id, app_id, lastping_at) VALUES ('%s', '%s', NOW()) ON DUPLICATE KEY UPDATE lastping_at=VALUES(lastping_at)`, deviceId, appId)
	if err != nil {
		raven.CaptureError(err, nil)
		return err
	}
	return nil
}
