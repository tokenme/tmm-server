package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func EditCreativeHanlder(c *gin.Context) {
	var req EditRequest
	db := Service.Db
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.OnlineStatus == 0 && req.Id == 0, `Invalid param`, c) {
		return
	}
	_, _, err := db.Query("UPDATE tmm.creatives SET online_status = %d WHERE id = %d", req.OnlineStatus, req.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    req,
	})
}
