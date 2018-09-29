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

type ShareUpdateRequest struct {
	Id           uint64          `json:"id" from:"id" binding:"required"`
	Title        string          `json:"title" form:"title"`
	Summary      string          `json:"summary" form:"summary"`
	Link         string          `json:"link" form:"link"`
	Image        string          `json:"image" form:"image"`
	Points       decimal.Decimal `json:"points" form:"points"`
	Bonus        decimal.Decimal `json:"bonus" form:"bonus"`
	MaxViewers   uint            `json:"max_viewers" form:"max_viewers"`
	OnlineStatus int8            `json:"online_status" from:"online_status"`
}

func ShareUpdateHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req ShareUpdateRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT points, points_left FROM tmm.share_tasks WHERE id=%d AND creator=%d LIMIT 1`, req.Id, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	oriPoints, _ := decimal.NewFromString(rows[0].Str(0))
	oriPointsLeft, _ := decimal.NewFromString(rows[0].Str(1))
	pointsLeft := oriPointsLeft
	if req.Points.GreaterThan(decimal.Zero) {
		query := `SELECT
d.id
FROM tmm.devices AS d
WHERE d.user_id = %d AND d.points >= %s
ORDER BY d.points DESC LIMIT 1`
		rows, _, err := db.Query(query, user.Id, req.Points)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		deviceId := rows[0].Str(0)
		pointsPerTs, err := common.GetPointsPerTs(Service)
		if CheckErr(err, c) {
			return
		}
		if req.Points.GreaterThan(oriPoints) {
			depositPoints := req.Points.Sub(oriPoints)
			ts := depositPoints.Div(pointsPerTs)
			_, ret, err := db.Query(`UPDATE tmm.devices AS d SET d.points=d.points-%s, d.consumed_ts = d.consumed_ts + %d WHERE id='%s' AND d.points>=%s AND d.user_id=%d`, depositPoints.String(), ts.IntPart(), db.Escape(deviceId), depositPoints.String(), user.Id)
			if CheckErr(err, c) {
				return
			}
			if CheckWithCode(ret.AffectedRows() == 0, NOT_ENOUGH_POINTS_ERROR, "not enough points in device", c) {
				return
			}
			_, _, err = db.Query(`UPDATE tmm.share_tasks AS st SET st.points=%s, st.points_left=st.points_left + %s WHERE st.id=%d AND st.creator=%d`, req.Points.String(), depositPoints.String(), req.Id, user.Id)
			if CheckErr(err, c) {
				return
			}
			pointsLeft = oriPointsLeft.Add(depositPoints)
		} else {
			_, ret, err := db.Query(`UPDATE tmm.share_tasks AS st SET st.points=%s, st.points_left=IF(st.points_left > %s, %s, st.points_left) WHERE st.id=%d AND st.creator=%d`, req.Points.String(), req.Points.String(), req.Points.String(), req.Id, user.Id)
			if CheckErr(err, c) {
				return
			}
			if ret.AffectedRows() > 0 {
				depositPoints := oriPoints.Sub(req.Points)
				ts := depositPoints.Div(pointsPerTs)
				_, _, err := db.Query(`UPDATE tmm.devices AS d SET d.points=d.points+%s, d.total_ts = d.total_ts + %d WHERE id='%s' AND d.user_id=%d`, depositPoints.String(), ts.IntPart(), db.Escape(deviceId), user.Id)
				if CheckErr(err, c) {
					return
				}
				if req.Points.GreaterThan(oriPointsLeft) {
					pointsLeft = oriPointsLeft.Add(depositPoints)
				} else {
					pointsLeft = oriPointsLeft.Sub(depositPoints)
				}
			}
		}
	}
	var updates []string
	if req.Title != "" {
		updates = append(updates, fmt.Sprintf("title='%s'", db.Escape(req.Title)))
	}
	if req.Summary != "" {
		updates = append(updates, fmt.Sprintf("summary='%s'", db.Escape(req.Summary)))
	}
	if req.Link != "" {
		updates = append(updates, fmt.Sprintf("link='%s'", db.Escape(req.Link)))
	}
	if req.Image != "" {
		updates = append(updates, fmt.Sprintf("image='%s'", db.Escape(req.Image)))
	}
	if req.Bonus.GreaterThan(decimal.Zero) {
		updates = append(updates, fmt.Sprintf("bonus=%s", req.Bonus.String()))
	}
	if req.MaxViewers > 0 {
		updates = append(updates, fmt.Sprintf("max_viewers=%d", req.MaxViewers))
	}
	if req.OnlineStatus == 1 || req.OnlineStatus == -1 {
		updates = append(updates, fmt.Sprintf("online_status=%d", req.OnlineStatus))
	}

	if req.OnlineStatus == -2 {
		query := `SELECT
d.id
FROM tmm.devices AS d
WHERE d.user_id = %d AND d.points >= %s
ORDER BY d.points DESC LIMIT 1`
		rows, _, err := db.Query(query, user.Id, req.Points)
		if CheckErr(err, c) {
			return
		}
		if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
			return
		}
		deviceId := rows[0].Str(0)
		pointsPerTs, err := common.GetPointsPerTs(Service)
		if CheckErr(err, c) {
			return
		}
		ts := oriPointsLeft.Div(pointsPerTs)
		_, _, err = db.Query(`UPDATE tmm.devices AS d, tmm.share_tasks AS st SET d.points=d.points+%s, d.total_ts = d.total_ts + %d, st.points_left = 0, st.online_status=-2 WHERE d.id='%s' AND d.user_id=%d AND st.id=%d AND st.creator=%d`, oriPointsLeft.String(), ts.IntPart(), db.Escape(deviceId), user.Id, req.Id, user.Id)
		if CheckErr(err, c) {
			return
		}
	}

	if len(updates) > 0 {
		_, _, err := db.Query(`UPDATE tmm.share_tasks SET %s WHERE id=%d AND creator=%d`, strings.Join(updates, ","), req.Id, user.Id)
		if CheckErr(err, c) {
			return
		}
	}
	now := time.Now().Format(time.RFC3339)
	task := common.ShareTask{
		Id:         req.Id,
		Title:      req.Title,
		Summary:    req.Summary,
		Link:       req.Link,
		Image:      req.Image,
		Points:     req.Points,
		PointsLeft: pointsLeft,
		Bonus:      req.Bonus,
		MaxViewers: req.MaxViewers,
		UpdatedAt:  now,
		Creator:    user.Id,
	}
	c.JSON(http.StatusOK, task)
}
