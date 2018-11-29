package task

import (
	"net/http"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"strconv"
	"fmt"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler/admin"
)

func GetTaskListHandler(c *gin.Context) {
	db := Service.Db
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "5"))
	cid, _ := strconv.Atoi(c.DefaultQuery(`cid`, "0"))
	nocid, _ := strconv.Atoi(c.DefaultQuery(`nocid`, "0"))
	var offset, count int
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}
	var sumquery string
	query := ` SELECT 
	s.id,
	s.creator,
	s.title,
    s.summary,
	s.link,
	s.image,
	s.points,
	s.points_left,
	s.bonus,
	s.max_viewers,
	s.viewers,
	s.online_status,
    s.inserted_at,
	s.updated_at FROM tmm.share_tasks AS s
	%s ORDER BY s.id LIMIT %d OFFSET %d `
	sumquery = `SELECT count(*) FROM tmm.share_tasks as s %s `
	if nocid != 1 {
		if cid != 0 {
			param := fmt.Sprintf(` INNER JOIN tmm.share_task_categories AS stc ON ( stc.task_id = s.id )  where stc.cid = %d `, cid)
			query = fmt.Sprintf(query, param, limit, offset)
			sumquery = fmt.Sprintf(sumquery, param)
		} else {
			query = fmt.Sprintf(query, "", limit, offset)
			sumquery = fmt.Sprintf(sumquery, "")
		}
	} else {
		isAuto := `WHERE NOT EXISTS (
		SELECT 1 FROM tmm.share_task_categories AS stc
		WHERE stc.is_auto = 1 AND stc.task_id = s.id
		LIMIT 1
		)`
		query = fmt.Sprintf(query, isAuto, limit, offset)
		sumquery = fmt.Sprintf(sumquery, isAuto)
	}
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `Not Find`, c) {
		return
	}
	var sharelist []*common.ShareTask
	for _, row := range rows {
		var cidList []int
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
		share := &common.ShareTask{
			Id:           row.Uint64(res.Map(`id`)),
			Creator:      row.Uint64(res.Map(`creator`)),
			Title:        row.Str(res.Map(`title`)),
			Summary:      row.Str(res.Map(`summary`)),
			Link:         row.Str(res.Map(`link`)),
			Image:        row.Str(res.Map(`image`)),
			Points:       points,
			PointsLeft:   pointsLeft,
			Bonus:        bonus,
			MaxViewers:   row.Uint(res.Map(`max_viewers`)),
			Viewers:      row.Uint(res.Map(`viewers`)),
			OnlineStatus: int8(row.Int(res.Map(`online_status`))),
			InsertedAt:   row.Str(res.Map(`inserted_at`)),
			UpdatedAt:    row.Str(res.Map(`updated_at`)),
		}
		cidquery := `SELECT cid FROM tmm.share_task_categories WHERE task_id = %d`
		rows, _, err = db.Query(cidquery, share.Id)
		for _, row := range rows {
			cidList = append(cidList, row.Int(0))
		}
		share.Cid = cidList
		sharelist = append(sharelist, share)
	}
	rows, _, err = db.Query(sumquery)
	if CheckErr(err, c) {
		return
	}
	count = rows[0].Int(0)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"curr_page": page,
			"data":      sharelist,
			"amount":    count,
		},
	})
}
