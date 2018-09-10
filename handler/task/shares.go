package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"time"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

type SharesRequest struct {
	Page     uint            `json:"page" form:"page"`
	PageSize uint            `json:"page_size" form:"page_size"`
	Idfa     string          `json:"idfa" form:"idfa"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
}

func SharesHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}

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

	device := common.DeviceRequest{
		Idfa:     req.Idfa,
		Platform: req.Platform,
	}
	deviceId := device.DeviceId()

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
    st.inserted_at,
    st.updated_at
FROM tmm.share_tasks AS st
WHERE st.points_left>0 AND st.online_status = 1
ORDER BY st.bonus DESC, st.id DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, (req.Page-1)*req.PageSize, req.PageSize)
	if CheckErr(err, c) {
		return
	}
	var tasks []common.ShareTask
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(6))
		points, _ := decimal.NewFromString(row.Str(7))
		pointsLeft, _ := decimal.NewFromString(row.Str(8))
		task := common.ShareTask{
			Id:         row.Uint64(0),
			Title:      row.Str(1),
			Summary:    row.Str(2),
			Link:       row.Str(3),
			Image:      row.Str(4),
			MaxViewers: row.Uint(5),
			Bonus:      bonus,
			Points:     points,
			PointsLeft: pointsLeft,
			InsertedAt: row.ForceLocaltime(9).Format(time.RFC3339),
			UpdatedAt:  row.ForceLocaltime(10).Format(time.RFC3339),
		}
		task.ShareLink, _ = task.GetShareLink(deviceId, Config)
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
