package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type BindRequest struct {
    Idfa     string `form:"idfa" json:"idfa"`
    Platform string `form:"platform" json:"platform"`
    Imei     string `form:"imei" json:"imei"`
    Mac      string `form:"mac" json:"mac"`
}

func BindHandler(c *gin.Context) {
	var req common.DeviceRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	err := saveDevice(Service, req, c)
	if CheckErr(err, c) {
		return
	}
	err = saveApp(Service, req)
	if CheckErr(err, c) {
		return
	}

	db := Service.Db

    if Check(req.Idfa == "" && req.Imei == "" && req.Mac == "", "invalid request", c) {
        return
    }
	rows, _, err := db.Query(`SELECT COUNT(*) FROM tmm.devices WHERE user_id=%d`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var deviceCount int
	if len(rows) > 0 {
		deviceCount = rows[0].Int(0)
	}
	if CheckWithCode(deviceCount >= Config.MaxBindDevice, MAX_BIND_DEVICE_ERROR, "exceeded maximum binding devices", c) {
		return
	}
	device := common.DeviceRequest{
		Idfa:     req.Idfa,
        Imei:     req.Imei,
        Mac:      req.Mac,
	}
	deviceId := device.DeviceId()
    if Check(len(deviceId) == 0, "not found", c) {
		return
	}
	_, ret, err := db.Query(`UPDATE tmm.devices SET user_id=%d WHERE id='%s' AND user_id=0`, user.Id, deviceId)
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() == 0 {
		rows, _, err := db.Query(`SELECT 1 FROM tmm.devices WHERE user_id=%d AND id='%s' LIMIT 1`, user.Id, deviceId)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, OTHER_BIND_DEVICE_ERROR, "the device has been bind by others", c) {
			return
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
