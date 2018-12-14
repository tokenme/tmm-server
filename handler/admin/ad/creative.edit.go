package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
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
	var SetList []string
	if req.OnlineStatus != 0 {
		SetList = append(SetList, fmt.Sprintf(" online_status = %d ", req.OnlineStatus))
	}
	if req.EndDate != "" {
		SetList = append(SetList, fmt.Sprintf(" end_date = '%s' ", db.Escape(req.EndDate)))
	}
	_, _, err := db.Query("UPDATE tmm.creatives SET %s WHERE id = %d", strings.Join(SetList, " , "), req.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    req,
	})
}
