package task

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"

)

func TaskUpdateHandler(c *gin.Context) {
	db := Service.Db
	var res OnlineStatusRequest
	if CheckErr(c.Bind(&res), c) {
		return
	}
	_, _, err := db.Query(`update tmm.share_tasks set online_status = %d where id = %d`, res.Status, res.TaskId)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: ""})
}