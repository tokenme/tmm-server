package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func AdzoneListHandler(c *gin.Context) {
	db := Service.Db

	query := `SELECT id,summary FROM tmm.adzones ORDER BY id DESC `
	var list []*common.Adzone
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}

	for _, row := range rows {
		list = append(list, &common.Adzone{
			Id:      row.Uint64(0),
			Summery: row.Str(1),
		}, )
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    list,
	})
}
