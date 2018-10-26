package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type PushTokenRequest struct {
	Idfa  string `form:"idfa" json:"idfa"`
	Imei  string `form:"imei" json:"imei"`
	Mac   string `form:"mac" json:"mac"`
	Token string `form:"token" json:"token" binding:"required"`
}

func PushTokenHandler(c *gin.Context) {
	var req PushTokenRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
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
	_, _, err := db.Query(`UPDATE tmm.devices SET push_token='%s' WHERE id='%s' AND user_id=%d`, req.Token, deviceId, user.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
