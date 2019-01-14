package account_controller

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
)

type EditRequest struct {
	Id       int    `json:"id"`
	Ban      bool   `json:"ban"`
	Comments string `json:"blocked_message"`
}

func EditAccountHandler(c *gin.Context) {
	var req EditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	query := `
	INSERT INTO tmm.user_settings (user_id,blocked,block_whitelist,comments) 
	VALUES(%d,%d,%d,'%s')  
    ON DUPLICATE KEY UPDATE blocked=VALUES(blocked),block_whitelist=VALUES(block_whitelist),comments=VALUES(comments)`
	var blockWhitelist int
	if !req.Ban {
		blockWhitelist = 1
	}

	db := Service.Db
	if _, _, err := db.Query(query, req.Id, 1, blockWhitelist, req.Comments); CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	})
}
