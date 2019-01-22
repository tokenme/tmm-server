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
	if Check(req.Id == 0, `Invalid param`, c) {
		return
	}
	var setList []string
	if req.OnlineStatus != 0 {
		setList = append(setList, fmt.Sprintf(" online_status = %d ", req.OnlineStatus))
	}
	if req.EndDate != "" {
		setList = append(setList, fmt.Sprintf(" end_date = '%s' ", db.Escape(req.EndDate)))
	}

	if req.Title != "" {
		setList = append(setList, fmt.Sprintf(" title = '%s' ", db.Escape(req.Title)))
	}
	if req.AdMode == "cpc" || req.AdMode == "cps" {
		setList = append(setList, fmt.Sprintf(" ad_mode = '%s' ", db.Escape(req.AdMode)))
	}
	if req.AdIncome != 0 {
		setList = append(setList, fmt.Sprintf(" ad_income = ad_income+%f", req.AdIncome))
	}
	if len(setList) > 0 {
		_, _, err := db.Query("UPDATE tmm.creatives SET %s WHERE id = %d", strings.Join(setList, " , "), req.Id)
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
