package task

import (
	"net/http"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"strconv"
	"fmt"
	"github.com/tokenme/tmm/common"
	"strings"
)

func GetTaskListHandler(c *gin.Context) {
	db := Service.Db
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "5"))
	cid, _ := strconv.Atoi(c.DefaultQuery(`cid`, "0"))
	nocid, _ := strconv.Atoi(c.DefaultQuery(`nocid`, "0"))
	var (
		offset, count           int
		sumquery, param, isAuto string
		query                   []string
		cidquery                = `SELECT cid FROM tmm.share_task_categories WHERE task_id = %d`
	)
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}
	query = append(query, `SELECT id,creator,title,
    summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,
    inserted_at,updated_at FROM tmm.share_tasks`)

	sumquery = `SELECT count(*) FROM tmm.share_tasks %s`
	if nocid != 1 {
		if cid != 0 {
			param = fmt.Sprintf(`INNER JOIN share_task_categories ON(id = task_id) where cid = %d`, cid)
			query = append(query, param)
			sumquery = fmt.Sprintf(sumquery, param)
		}
	} else {
		isAuto = `WHERE id NOT IN
(SELECT id  FROM tmm.share_tasks INNER JOIN share_task_categories ON(id = task_id) WHERE is_auto = 1)`
		query = append(query, isAuto)
		sumquery = fmt.Sprintf(sumquery, isAuto)
	}
	rows, result, err := db.Query(`%s ORDER BY id DESC LIMIT %d OFFSET %d`, strings.Join(query, ` `), limit, offset)
	if CheckErr(err, c) {
		return
	}
	var sharelist []common.ShareTask
	for _, row := range rows {
		var cidList []int
		points, err := decimal.NewFromString(row.Str(result.Map(`points`)))
		if CheckErr(err, c) {
			return
		}
		bonus, err := decimal.NewFromString(row.Str(result.Map(`bonus`)))
		if CheckErr(err, c) {
			return
		}
		pointsLeft, err := decimal.NewFromString(row.Str(result.Map(`points_left`)))
		if CheckErr(err, c) {
			return
		}
		share := common.ShareTask{
			Id:           row.Uint64(result.Map(`id`)),
			Creator:      row.Uint64(result.Map(`creator`)),
			Title:        row.Str(result.Map(`title`)),
			Summary:      row.Str(result.Map(`summary`)),
			Link:         row.Str(result.Map(`link`)),
			Image:        row.Str(result.Map(`image`)),
			Points:       points,
			PointsLeft:   pointsLeft,
			Bonus:        bonus,
			MaxViewers:   row.Uint(result.Map(`max_viewers`)),
			Viewers:      row.Uint(result.Map(`viewers`)),
			OnlineStatus: int8(row.Int(result.Map(`online_status`))),
			InsertedAt:   row.Str(result.Map(`inserted_at`)),
			UpdatedAt:    row.Str(result.Map(`updated_at`)),
		}
		rows, _, err = db.Query(cidquery, share.Id)
		for _, row := range rows {
			cidList = append(cidList, row.Int(0))
		}
		share.Cid = cidList
		sharelist = append(sharelist, share)
	}
	rows, _, err = db.Query(sumquery)
	if err != nil {
		return
	}
	count = rows[0].Int(0)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": gin.H{
			"curr_page": page,
			"data":      sharelist,
			"amount":    count,
		},
	})
}
