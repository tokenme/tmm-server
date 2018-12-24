package task

import (
	"net/http"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/shopspring/decimal"
	"fmt"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler/admin"
	"strings"
)

type SearchOptions struct {
	Page                 int `form:"page"`
	Limit                int `form:"limit"`
	Cid                  int `form:"cid"`
	Nocid                int `form:"nocid"`
	Online               int `form:"online_status"`
	IsTask               int `form:"is_task"`
	SortByReadersNumbers int `form:"readersNumbersSort"`
}

func GetTaskListHandler(c *gin.Context) {
	db := Service.Db
	var req SearchOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	var offset int

	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
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
	log.users AS users,
	log.ts AS ts,	
	s.updated_at FROM tmm.share_tasks AS s
	LEFT JOIN (SELECT COUNT(1) AS users ,task_id,SUM(ts) AS ts FROM tmm.reading_logs GROUP BY task_id) AS log ON (log.task_id = s.id)
	%s 
	ORDER BY %s
	LIMIT %d OFFSET %d `
	sumquery = `SELECT count(*) FROM tmm.share_tasks as s %s `
	var where []string
	var sumwhere []string
	if req.Nocid != 1 {
		if req.Cid != 0 {
			taskType := fmt.Sprintf(` INNER JOIN tmm.share_task_categories AS stc ON ( stc.task_id = s.id )  
			where stc.cid = %d `, req.Cid)
			where = append(where, taskType)
			sumwhere = append(sumwhere, taskType)
		} else {
			param := ` WHERE 1 = 1 `
			where = append(where, param)
			sumwhere = append(sumwhere, param)
		}
	} else {
		isAuto := ` WHERE NOT EXISTS (
		SELECT 1 FROM tmm.share_task_categories AS stc
		WHERE stc.is_auto = 1 AND stc.task_id = s.id
		LIMIT 1 ) `
		where = append(where, isAuto)
		sumwhere = append(sumwhere, isAuto)
	}
	if req.Online != 0 {
		isOnline := fmt.Sprintf(` And s.online_status = %d `, req.Online)
		where = append(where, isOnline)
		sumwhere = append(sumwhere, isOnline)
	}
	if req.IsTask != -1 {
		istask := fmt.Sprint(" AND s.link not like  'https://tmm.tokenmama.io/article/show%' ")
		where = append(where, istask)
		sumwhere = append(sumwhere, istask)
	}
	var order string
	if req.SortByReadersNumbers != -1 {
		order = ` log.users DESC ,s.id DESC `
	} else {
		order = ` s.id DESC`
	}
	rows, res, err := db.Query(query, strings.Join(where, ` `), order, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: "Not Find Data",
			Data: gin.H{
				"curr_page": req.Page,
				"data":      nil,
				"amount":    0,
			},
		})
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
			Id:            row.Uint64(res.Map(`id`)),
			Creator:       row.Uint64(res.Map(`creator`)),
			Title:         row.Str(res.Map(`title`)),
			Summary:       row.Str(res.Map(`summary`)),
			Link:          row.Str(res.Map(`link`)),
			Image:         row.Str(res.Map(`image`)),
			Points:        points,
			PointsLeft:    pointsLeft,
			Bonus:         bonus,
			MaxViewers:    row.Uint(res.Map(`max_viewers`)),
			Viewers:       row.Uint(res.Map(`viewers`)),
			OnlineStatus:  int8(row.Int(res.Map(`online_status`))),
			InsertedAt:    row.Str(res.Map(`inserted_at`)),
			UpdatedAt:     row.Str(res.Map(`updated_at`)),
			TotalReadUser: row.Int(res.Map(`users`)),
			ReadDuration:  row.Int(res.Map(`ts`)),
		}
		cidquery := `SELECT cid FROM tmm.share_task_categories WHERE task_id = %d`
		rows, _, err = db.Query(cidquery, share.Id)
		for _, row := range rows {
			cidList = append(cidList, row.Int(0))
		}
		share.Cid = cidList
		sharelist = append(sharelist, share)
	}
	rows, _, err = db.Query(sumquery, strings.Join(sumwhere, ` `))
	if CheckErr(err, c) {
		return
	}
	count := rows[0].Int(0)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"curr_page": req.Page,
			"data":      sharelist,
			"amount":    count,
		},
	})
}
