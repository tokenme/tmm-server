package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func AddCreativeHandler(c *gin.Context) {
	db := Service.Db
	var creat Creatives
	if CheckErr(c.Bind(&creat), c) {
		return
	}
	if Check(creat.AdgroupId == 0 || creat.Platform != 1 && creat.Platform != -1 || creat.Title == "" || creat.Link == "" || creat.Image == "", `Invalid param`, c) {
		return
	}

	var valueList, fieldList []string
	if creat.AdgroupId != 0 {
		valueList = append(valueList, fmt.Sprintf("%d", creat.AdgroupId))
		fieldList = append(fieldList, fmt.Sprintf("adgroup_id"))
	}

	if creat.Title != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(creat.Title)))
		fieldList = append(fieldList, fmt.Sprintf("title"))
	}

	if creat.Image != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(creat.Image)))
		fieldList = append(fieldList, fmt.Sprintf("image"))
	}

	if creat.ShareImage != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(creat.ShareImage)))
		fieldList = append(fieldList, fmt.Sprintf("share_image"))
	}

	if creat.Link != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(creat.Link)))
		fieldList = append(fieldList, fmt.Sprintf("link"))
	}
	if creat.Width != 0 {
		valueList = append(valueList, fmt.Sprintf(" %d ", creat.Width))
		fieldList = append(fieldList, fmt.Sprintf("width"))
	}
	if creat.Height != 0 {
		valueList = append(valueList, fmt.Sprintf(" %d ", creat.Height))
		fieldList = append(fieldList, fmt.Sprintf("height"))
	}
	if creat.Platform == 1 || creat.Platform == 2 {
		valueList = append(valueList, fmt.Sprintf(" %d ", creat.Platform))
		fieldList = append(fieldList, fmt.Sprintf("platform"))
	}

	if len(valueList) > 0 {
		_, _, err := db.Query(`INSERT tmm.creatives(%s )
		VALUES (%s)`, strings.Join(fieldList, ","), strings.Join(valueList, ","))
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    creat,
	})
}
