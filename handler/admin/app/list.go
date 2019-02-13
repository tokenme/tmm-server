package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
)

type SearchOptions struct {
	Online   int    `form:"online"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Title    string `form:"title"`
}

func GetShareAppHandler(c *gin.Context) {
	db := Service.Db
	var req SearchOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var where []string
	var sumwhere []string
	if req.Online == 1 || req.Online == -1 {
		where = append(where, fmt.Sprintf(" AND online_status = %d ", req.Online))
		sumwhere = append(sumwhere, fmt.Sprintf(" AND online_status = %d ", req.Online))
	}
	var offset, limit int
	if req.PageSize > 0 {
		limit = req.PageSize
	} else {
		limit = 25
	}
	if req.Page > 1 {
		offset = limit * (req.Page - 1)
	}

	query := `SELECT id, bundle_id, name, size, bonus, download_url, icon, 
				(SELECT COUNT(1) FROM tmm.device_app_tasks WHERE task_id = id AND (status = 1  OR status = 2))  AS downloads,
				points, points_left, online_status, inserted_at,details
				FROM tmm.app_tasks 
				WHERE 1 = 1 %s 
				ORDER BY id DESC 
				LIMIT  %d 
				OFFSET %d `
	sumquery := `SELECT COUNT(1) FROM tmm.app_tasks WHERE 1 = 1 %s `
	rows, res, err := db.Query(query, strings.Join(where, " "), limit, offset)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0, admin.Not_Found, c) {
		return
	}

	var lists []*common.AppTask
	for _, row := range rows {
		points, err := decimal.NewFromString(row.Str(res.Map(`points`)))
		if CheckErr(err, c) {
			return
		}
		bonus, err := decimal.NewFromString(row.Str(res.Map(`bonus`)))
		if CheckErr(err, c) {
			return
		}
		pointsLeft, err := decimal.NewFromString(row.Str(res.Map(`points_left`)))
		if CheckErr(err, c) {
			return
		}
		apps := &common.AppTask{
			Id:           row.Uint64(res.Map(`id`)),
			BundleId:     row.Str(res.Map(`bundle_id`)),
			Name:         row.Str(res.Map(`name`)),
			Size:         row.Uint(res.Map(`size`)),
			Bonus:        bonus,
			DownloadUrl:  row.Str(res.Map(`download_url`)),
			Icon:         row.Str(res.Map(`icon`)),
			Downloads:    row.Uint(res.Map(`downloads`)),
			Points:       points,
			PointsLeft:   pointsLeft,
			OnlineStatus: int8(row.Int(res.Map(`online_status`))),
			InsertedAt:   row.Str(res.Map(`inserted_at`)),
			Details:      row.Str(res.Map(`details`)),
		}
		lists = append(lists, apps)
	}

	rows, _, err = db.Query(sumquery, strings.Join(sumwhere, ` `))
	if CheckErr(err, c) {
		return
	}
	var count int

	if len(rows) > 0 {
		count = rows[0].Int(0)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"data":      lists,
			"amount":    count,
		},
	})
}
