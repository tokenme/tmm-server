package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type EditRequest struct {
	Id  int  `json:"id"`
	Ban bool `json:"ban"`
}

func EditAccountHandler(c *gin.Context) {
	db := Service.Db
	var req EditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	query := `
INSERT INTO tmm.user_settings (user_id,blocked,block_whitelist) VALUES(%d,%d,%d)  ON DUPLICATE KEY UPDATE blocked=VALUES(blocked),block_whitelist=VALUES(block_whitelist)`
	var err error
	if req.Ban {
		_, _, err = db.Query(query, req.Id, 1, 0)

	} else {
		_, _, err = db.Query(query, req.Id, 0, 0)
	}
	if CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	})
}
