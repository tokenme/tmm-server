package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"github.com/tokenme/tmm/common"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type AddGroupRequest struct {
	GroupTitle  string `json:"group_title"`
	Location    int    `json:"location"`
	Page        int    `json:"page"`
	ChannelId   int    `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

type Creatives struct {
	Title        string  `json:"title,omitempty"`
	InsertAt     string  `json:"insert_at,omitempty"`
	UpdateAt     string  `json:"update_at,omitempty"`
	StartDate    string  `json:"start_date,omitempty"`
	EndDate      string  `json:"end_date,omitempty"`
	Platform     int     `json:"platform,omitempty"`
	OnlineStatus int     `json:"online_status,omitempty"`
	ShareImage   string  `json:"share_image,omitempty"`
	AdGroup      AdGroup `json:"ad_group,omitempty"`
	Adzone       Adzone  `json:"adzone,omitempty"`
	AdMode       string  `json:"ad_mode"`
	AdIncome     string  `json:"ad_income"`
	common.Creative
	Data
}

type AdGroup struct {
	common.Adgroup
	Title string `json:"title,omitempty"`
	Data
}
type Adzone struct {
	common.Adzone
	Data
}

type Data struct {
	Imp   int `json:"impress,omitempty"`
	Click int `json:"click,omitempty"`
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
	if Check(req.ChannelId < 0 || req.Page < 0 || req.Location < 0, `Invalid param`, c) {
		return
	}
	var adzId uint64
	rows, _, err := db.Query(`SELECT id  FROM tmm.adzones WHERE cid = %d AND page = %d AND idx = %d`, req.ChannelId, req.Page, req.Location)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		summary := fmt.Sprintf("文章%s-%d刷-%d", req.ChannelName, req.Page, req.Location)
		_, res, err := db.Query(`INSERT INTO tmm.adzones(summary,cid,page,idx) VALUES('%s',%d,%d,%d)`, db.Escape(summary), req.ChannelId, req.Page, req.Location)
		if CheckErr(err, c) {
			return
		}
		adzId = res.InsertId()
	} else {
		adzId = rows[0].Uint64(0)
	}
	query := `INSERT tmm.adgroups(adzone_id,user_id,title)VALUES(%d,%d,'%s')  `
	_, _, err = db.Query(query, adzId, user.Id, db.Escape(req.GroupTitle))
	if CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    req,
	})

}
