package task

import (
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/videospider"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

type SharesRequest struct {
	Page     uint            `json:"page" form:"page"`
	PageSize uint            `json:"page_size" form:"page_size"`
	Idfa     string          `json:"idfa" form:"idfa"`
	Imei     string          `json:"imei" form:"imei"`
	Mac      string          `json:"mac" form:"mac"`
	Platform common.Platform `json:"platform" form:"platform" binding:"required"`
	MineOnly bool            `json:"mine_only" form:"mine_only"`
	IsVideo  bool            `json:"is_video" form:"is_video"`
	Build    uint            `json:"build" form:"build"`
	Cid      uint            `json:"cid" form:"cid"`
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

	device := common.DeviceRequest{
		Idfa: req.Idfa,
		Imei: req.Imei,
		Mac:  req.Mac,
	}
	deviceId := device.DeviceId()
	if Check(len(deviceId) == 0, "not found", c) {
		return
	}

	var taskIds []uint64
	limitState := fmt.Sprintf("LIMIT %d, %d", (req.Page-1)*req.PageSize, req.PageSize)
	onlineStatusConstrain := "st.points_left>0 AND st.online_status=1"
	var inCidConstrain string
	orderBy := "st.bonus DESC, st.id DESC"
	if req.MineOnly {
		onlineStatusConstrain = fmt.Sprintf("st.creator = %d", user.Id)
		orderBy = "st.id DESC"
	}
	if req.IsVideo {
		onlineStatusConstrain = "st.is_video=1 AND st.points_left>0 AND st.online_status=1"
	}
	if req.Cid > 0 {
		inCidConstrain = fmt.Sprintf("INNER JOIN tmm.share_task_categories AS stc ON (stc.task_id=st.id AND stc.cid=%d)", req.Cid)
	} else if !req.MineOnly && !req.IsVideo {
		taskIds = SuggestEngine.Match(user.Id, req.Page, req.PageSize)
	}
	if len(taskIds) > 0 {
		var tids []string
		for _, tid := range taskIds {
			tids = append(tids, fmt.Sprintf("%d", tid))
		}
		onlineStatusConstrain = fmt.Sprintf("st.id IN (%s)", strings.Join(tids, ","))
		orderBy = "st.bonus DESC"
		limitState = ""
	}

	showBonusHint := true
	if req.Idfa != "" && req.Build == Config.App.SubmitBuild {
		showBonusHint = true
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
    st.video_link,
    st.is_video,
    st.online_status
FROM tmm.share_tasks AS st
%s
WHERE %s
ORDER BY %s %s`
	rows, _, err := db.Query(query, inCidConstrain, onlineStatusConstrain, orderBy, limitState)
	if CheckErr(err, c) {
		return
	}
	var tasks []*common.ShareTask
	var wg sync.WaitGroup
	videopider := videospider.NewClient(Service, Config)
	videoFetchPool, _ := ants.NewPoolWithFunc(10, func(req interface{}) error {
		defer wg.Done()
		task := req.(*common.ShareTask)
		video, err := videopider.Get(task.Link)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		if len(video.Files) == 0 {
			return errors.New("invalid video")
		}
		sorter := videospider.NewVideoSorter(video.Files)
		sort.Sort(sort.Reverse(sorter))
		task.VideoLink = sorter[0].Link
		return nil
	})
	for _, row := range rows {
		bonus, _ := decimal.NewFromString(row.Str(6))
		points, _ := decimal.NewFromString(row.Str(7))
		pointsLeft, _ := decimal.NewFromString(row.Str(8))
		creator := row.Uint64(12)
		task := &common.ShareTask{
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
			VideoLink:     row.Str(13),
			IsVideo:       uint8(row.Uint(14)),
			ShowBonusHint: showBonusHint,
		}
		if strings.HasPrefix(task.Link, "https://tmm.tokenmama.io/article/show") {
			task.Link = strings.Replace(task.Link, "https://tmm.tokenmama.io/article/show", "https://static.tianxi100.com/article/show", -1)
		}
		if creator == user.Id {
			task.Viewers = row.Uint(9)
			task.Creator = creator
			task.OnlineStatus = int8(row.Int(15))
		}
		task.ShareLink, _ = task.GetShareLink(deviceId, Config)
		if task.IsVideo == 1 && (strings.Contains(task.VideoLink, "krcom.cn") || strings.Contains(task.VideoLink, "sinaimg.cn")) {
			wg.Add(1)
			videoFetchPool.Serve(task)
		}
		tasks = append(tasks, task)
	}
	wg.Wait()
	c.JSON(http.StatusOK, tasks)
}
