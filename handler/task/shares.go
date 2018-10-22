package task

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
	"time"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

type SharesRequest struct {
	Page       uint            `json:"page" form:"page"`
	PageSize   uint            `json:"page_size" form:"page_size"`
	Idfa       string          `json:"idfa" form:"idfa"`
	Imei       string          `json:"imei" form:"imei"`
	Mac        string          `json:"mac" form:"mac"`
	Platform   common.Platform `json:"platform" form:"platform" binding:"required"`
	MineOnly   bool            `json:"mine_only" form:"mine_only"`
	AppVersion string          `json:"app_version" form:"app_version"`
}

func SharesHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req SharesRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}

	onlineStatusConstrain := "AND st.online_status = 1"
	orderBy := "st.bonus DESC, st.id DESC"
	if req.MineOnly {
		onlineStatusConstrain = fmt.Sprintf("AND st.creator = %d", user.Id)
		orderBy = "st.id DESC"
	}
	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	showBonusHint := true
	if req.Idfa != "" && strings.Compare(req.AppVersion, Config.AppReleaseVersion.IOS) == 1 {
		showBonusHint = false
	}

	db := Service.Db
	query := `SELECT
    st.id,
    st.title,
    st.summary,
    st.link,
    st.image,
    st.max_viewers,
    st.bonus,
    st.points,
    st.points_left,
    st.viewers,
    st.inserted_at,
    st.updated_at,
    st.creator,
    st.online_status
FROM tmm.share_tasks AS st
WHERE st.points_left>0 %s
ORDER BY %s LIMIT %d, %d`
	rows, _, err := db.Query(query, onlineStatusConstrain, orderBy, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var tasks []common.ShareTask
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(6))
		points, _ := decimal.NewFromString(row.Str(7))
		pointsLeft, _ := decimal.NewFromString(row.Str(8))
		creator := row.Uint64(12)
		task := common.ShareTask{
			Id:            row.Uint64(0),
			Title:         row.Str(1),
			Summary:       row.Str(2),
			Link:          row.Str(3),
			Image:         row.Str(4),
			MaxViewers:    row.Uint(5),
			Bonus:         bonus,
			Points:        points,
			PointsLeft:    pointsLeft,
			InsertedAt:    row.ForceLocaltime(10).Format(time.RFC3339),
			UpdatedAt:     row.ForceLocaltime(11).Format(time.RFC3339),
			ShowBonusHint: showBonusHint,
		}
		if creator == user.Id {
			task.Viewers = row.Uint(9)
			task.Creator = creator
			task.OnlineStatus = int8(row.Int(13))
		}
		task.ShareLink, _ = task.GetShareLink(deviceId, Config)
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
