package withdraw

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type EditRequest struct {
	Id     string `json:"id"`
	Types  int    `json:"types"`
	Status int    `json:"status"`
	UserId int    `json:"user_id"`
}

func EditWithDrawHandler(c *gin.Context) {
	db := Service.Db

	var req EditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	query := ``
	if req.Types == UC {
		query = `UPDATE tmm.withdraw_txs SET verified = %d WHERE tx = '%s' AND user_id = %d AND verified = 0`
	} else if req.Types == Point {
		query = `UPDATE tmm.point_withdraws SET verified = %d WHERE id = '%s' AND user_id = %d AND verified = 0`
	}
	if _, _, err := db.Query(query, req.Status, req.Id, req.UserId); CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
	})
}
