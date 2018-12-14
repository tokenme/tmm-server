package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"github.com/tokenme/tmm/common"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"strings"
)

type AddGroupRequest struct {
	GroupTitle  string `json:"group_title"`
	Location    int    `json:"location"`
	Page        int    `json:"page"`
	ChannelId   int    `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Creatives
}

type Creatives struct {
	Title        string `json:"title"`
	InsertAt     string `json:"insert_at,omitempty"`
	UpdateAt     string `json:"update_at,omitempty"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
	Platform     int    `json:"platform,omitempty"`
	OnlineStatus int    `json:"online_status,omitempty"`
	ShareImage   string `json:"share_image,omitempty"`
	common.Creative
	Data
}

type AdGroup struct {
	common.Adgroup
	Title     string    `json:"title"`
	Creatives Creatives `json:"creatives"`
	Data
}
type Adzone struct {
	common.Adzone
	Group AdGroup `json:"group"`
	Data
}

type Data struct {
	Img   int `json:"img"`
	Click int `json:"click"`
}

func AddGroupsHandler(c *gin.Context) {
	db := Service.Db
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req AddGroupRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if Check(req.Title == "" || req.ChannelId <  0 || req.Page < 0 || req.Location < 0 || req.Image == "" || req.Link == "" || req.Platform != 1 && req.Platform != -1 , `Invalid param`, c) {
		return
	}
	var adzId uint64
	rows, _, err := db.Query(`SELECT id  FROM tmm.adzones WHERE cid = %d AND page = %d AND idx = %d`, req.ChannelId, req.Page, req.Location)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		summary := fmt.Sprintf("文章%s%d刷%d", req.ChannelName, req.Page, req.Location)
		_, res, err := db.Query(`INSERT INTO tmm.adzones(summary,cid,page,idx) VALUES('%s',%d,%d,%d)`, db.Escape(summary), req.ChannelId, req.Page, req.Location)
		if CheckErr(err, c) {
			return
		}
		adzId = res.InsertId()
	} else {
		adzId = rows[0].Uint64(0)
	}
	query := `INSERT tmm.adgroups(adzone_id,user_id,title)VALUES(%d,%d,'%s')  `
	_, res, err := db.Query(query, adzId, user.Id, db.Escape(req.GroupTitle))
	if CheckErr(err, c) {
		return
	}
	var valueList, fieldList []string
	id := res.InsertId()
	if id != 0 {
		valueList = append(valueList, fmt.Sprintf("%d", id))
		fieldList = append(fieldList, fmt.Sprintf("adgroup_id"))
	}

	if req.Title != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(req.Title)))
		fieldList = append(fieldList, fmt.Sprintf("title"))
	}

	if req.Image != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(req.Image)))
		fieldList = append(fieldList, fmt.Sprintf("image"))
	}

	if req.Link != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(req.Link)))
		fieldList = append(fieldList, fmt.Sprintf("link"))
	}

	if req.ShareImage != "" {
		valueList = append(valueList, fmt.Sprintf(" '%s' ", db.Escape(req.ShareImage)))
		fieldList = append(fieldList, fmt.Sprintf("share_image"))
	}

	if req.Platform == 1 || req.Platform == 2 {
		valueList = append(valueList, fmt.Sprintf(" %d ", req.Platform))
		fieldList = append(fieldList, fmt.Sprintf("platform"))
	}

	if req.Width != 0 {
		valueList = append(valueList, fmt.Sprintf(" %d ", req.Width))
		fieldList = append(fieldList, fmt.Sprintf("width"))
	}
	if req.Height != 0 {
		valueList = append(valueList, fmt.Sprintf(" %d ", req.Height))
		fieldList = append(fieldList, fmt.Sprintf("height"))
	}

	if len(valueList) > 0 {
		_, _, err = db.Query(`INSERT tmm.creatives(%s )
		VALUES (%s)`, strings.Join(fieldList, ","), strings.Join(valueList, ","))
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
