package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"fmt"
	"strings"
)

type EditRequest struct {
	Id           int    `form:"id"`
	Title        string `form:"title"`
	OnlineStatus int    `form:"online_status"`
	EndDate 	 string `form:"end_date"`
}

func EditGroupHanler(c *gin.Context) {
	db := Service.Db
	var req EditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var setList []string
	if req.Title != "" {
		setList = append(setList, fmt.Sprintf(" title = '%s' ", db.Escape(req.Title)))
	}
	if req.OnlineStatus == 1 || req.OnlineStatus == -1 {
		setList = append(setList, fmt.Sprintf(" online_status = %d ", req.OnlineStatus))
	}
	_, _, err := db.Query(`UPDATE tmm.adgroups SET %s WHERE id = %d`, strings.Join(setList, " , "), req.Id)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    req,
	})
}
