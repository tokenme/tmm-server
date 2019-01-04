package device

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type UnbindRequest struct {
	Id string `json:"id" form:"id" binding:"required"`
}

func UnbindHandler(c *gin.Context) {
	var req UnbindRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db

	if Check(req.Id == "", "invalid request", c) {
		return
	}
	_, _, err := db.Query(`INSERT IGNORE INTO tmm.user_devices (user_id, device_id) VALUES (%d, '%s')`, user.Id, db.Escape(req.Id))
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`UPDATE tmm.devices SET user_id=0 WHERE user_id=%d AND id='%s'`, user.Id, req.Id)
	if CheckErr(err, c) {
		log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
