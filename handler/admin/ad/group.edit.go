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
	Id           int             `form:"id",json:"id,omitempty"`
	Title        string          `form:"title",json:"title,omitempty"`
	OnlineStatus int             `form:"online_status",json:"online_status,omitempty"`
	EndDate      string          `form:"end_date",json:"end_date,omitempty"`
	AdzoneId     int             `form:"adzone_id",json:"adzone_id,omitempty"`
	AdMode       string          `form:"ad_mode",json:"ad_mode,omitempty"`
	AdIncome     float64 `form:"ad_income",json:"ad_income,omitempty"`
}

func EditGroupHanler(c *gin.Context) {
	db := Service.Db
	var req EditRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.Id == 0, `Invalid param`, c) {
		return
	}
	var setList []string
	if req.Title != "" {
		setList = append(setList, fmt.Sprintf(" title = '%s' ", db.Escape(req.Title)))
	}
	if req.OnlineStatus == 1 || req.OnlineStatus == -1 {
		setList = append(setList, fmt.Sprintf(" online_status = %d ", req.OnlineStatus))
	}
	if req.AdzoneId != 0 {
		setList = append(setList, fmt.Sprintf(" adzone_id = %d ", req.AdzoneId))
	}
	if len(setList) > 0 {
		_, _, err := db.Query(`UPDATE tmm.adgroups SET %s WHERE id = %d`, strings.Join(setList, " , "), req.Id)
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    req,
	})
}
