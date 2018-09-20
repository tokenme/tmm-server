package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type BindRequest struct {
	Idfa string `form:"idfa" json:"idfa"`
}

func BindHandler(c *gin.Context) {
	var req BindRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db

	if Check(req.Idfa == "", "invalid request", c) {
		return
	}
	deviceRequest := common.DeviceRequest{
		Platform: common.IOS,
		Idfa:     req.Idfa,
	}
	_, _, err := db.Query(`UPDATE tmm.devices SET user_id=%d WHERE id='%s' AND user_id=0`, user.Id, deviceRequest.DeviceId())
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
