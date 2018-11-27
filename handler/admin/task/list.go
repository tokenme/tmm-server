package task

import (
	"net/http"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"strconv"
	"fmt"
	"github.com/tokenme/tmm/common"
)

func GetShareListHandler(c *gin.Context) {
	db := Service.Db
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "5"))
	cid, _ := strconv.Atoi(c.DefaultQuery(`cid`, "0"))
	var (
		offset, count   int
		query, sumquery string
		cidQuery        = `select cid from tmm.share_task_categories where task_id = %d`
	)
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}
	if cid == 0 {
		query = fmt.Sprintf(`select id,creator,title,summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,inserted_at,updated_at from tmm.share_tasks order by id DESC limit %d offset %d`, limit, offset)
		sumquery = `select count(*) from tmm.share_tasks`
	} else {
		query = fmt.Sprintf(`select id,creator,title,summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,inserted_at,updated_at from tmm.share_tasks LEFT JOIN 
	share_task_categories ON(id = task_id) where cid = %d order by id DESC 
	limit %d offset %d`, cid, limit, offset)
		sumquery = fmt.Sprintf(`select count(*) from tmm.share_tasks INNER JOIN
	    share_task_categories ON(id = task_id) where cid = %d`, cid)
	}
	rows, result, err := db.Query(query)
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
		rows, _, err = db.Query(cidQuery, share.Id)
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
		"code": http.StatusOK,
		"msg":  "",
		"data": gin.H{
			"curr_page": page,
			"data":      sharelist,
			"amount":    count,
		},
	})
}
