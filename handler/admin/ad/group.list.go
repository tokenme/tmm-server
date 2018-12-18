package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/tokenme/tmm/common"
)

func GroupListHanlder(c *gin.Context) {
	db := Service.Db
	var List []common.Adgroup

	query := `SELECT id,title,online_status,adzone_id FROM tmm.adgroups ORDER BY id DESC`

	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	for _, row := range rows {
		adGroup := common.Adgroup{
			Title:        row.Str(1),
			Id:           row.Uint64(0),
			OnlineStatus: row.Int(2),
		}
		adGroup.Adzone = &common.Adzone{Id:row.Uint64(3)}
		List = append(List, adGroup)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    List,
	})

}
