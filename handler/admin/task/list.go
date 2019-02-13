package task

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"strings"
	"strconv"
)

type SearchOptions struct {
	admin.Pages
	Cid                  int `form:"cid"`
	Nocid                int `form:"nocid"`
	Online               int `form:"online_status"`
	IsCrawled            int `form:"is_task"`
	SortByReadersNumbers int `form:"readersNumbersSort"`
}

func GetTaskListHandler(c *gin.Context) {
	var req SearchOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}

	var offset int
	if req.Page > 0 {
		offset = (req.Page - 1) * req.Limit
	}

	query := ` 
SELECT 
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
	COUNT(DISTINCT log.user_id) AS users,
	SUM(log.ts) AS ts,
	s.updated_at,
	GROUP_CONCAT(DISTINCT stc.cid) AS cid_str
FROM 
	tmm.share_tasks AS s
LEFT JOIN 
	tmm.reading_logs AS log ON (log.task_id = s.id)
LEFT JOIN 
	share_task_categories AS stc ON  (stc.task_id = s.id)
WHERE
	 %s 
GROUP BY 
	s.id 
ORDER BY 
	%s
LIMIT %d OFFSET %d `

	sumquery := `
SELECT 
	count(1) 
FROM 
	tmm.share_tasks as s 
LEFT JOIN 
	share_task_categories AS stc ON  stc.task_id = s.id
WHERE 
	%s `

	var where []string
	var sumwhere []string

	if req.Nocid == -1 {
		if req.Cid > 0 {
			where = append(where, fmt.Sprintf(` stc.cid = %d `, req.Cid))
			sumwhere = append(sumwhere, fmt.Sprintf(` stc.cid = %d `, req.Cid))
		}

	} else {
		isAuto := ` 
		stc.is_auto < 1
		`
		where = append(where, isAuto)
		sumwhere = append(sumwhere, isAuto)
	}

	var isOnline string
	if req.Online == 1 {
		isOnline = fmt.Sprintf(` s.online_status = %d AND s.points_left > s.bonus`, req.Online)
	} else {
		isOnline = fmt.Sprintf(` s.online_status = %d `, req.Online)
	}
	where = append(where, isOnline)
	sumwhere = append(sumwhere, isOnline)

	if req.IsCrawled != -1 {
		isCrawled := fmt.Sprint("  s.is_crawled = 0 ")
		where = append(where, isCrawled)
		sumwhere = append(sumwhere, isCrawled)
	}

	var order string
	if req.SortByReadersNumbers == 1 {
		order = ` users DESC  `
	} else {
		order = ` s.id DESC`
	}

	db := Service.Db
	rows, res, err := db.Query(query, strings.Join(where, ` AND `), order, req.Limit, offset)
	if CheckErr(err, c) {
		return
	}

	if Check(len(rows) == 0  ,admin.Not_Found,c){
		return
	}

	var sharelist []*common.ShareTask
	var taskIdList []string
	for _, row := range rows {
		var cidList []int
		taskIdList = append(taskIdList,row.Str(res.Map(`id`)))
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
		cidStrings:=strings.Split(row.Str(res.Map(`cid_str`)),",")

		for _,cidStr:=range cidStrings{
			cid,_:=strconv.Atoi(cidStr)
			cidList = append(cidList,cid)
		}

		share.Cid = cidList
		sharelist = append(sharelist, share)
	}

	rows,_,err=db.Query(`SELECT SUM(pv),task_id FROM tmm.share_task_stats WHERE task_id IN (%s) GROUP BY task_id `,strings.Join(taskIdList,`,`))

	if CheckErr(err,c){
		return
	}

	pvMap:= make(map[uint64]int)
	for _,row:=range rows{
		pvMap[row.Uint64(1)] = row.Int(0)
	}

	for _,share:=range sharelist{
		if value,ok:=pvMap[share.Id];ok{
			share.Pv = value
		}
	}
	rows, _, err = db.Query(sumquery, strings.Join(sumwhere, ` AND `))
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
			"curr_page": req.Page,
			"data":      sharelist,
			"amount":    count,
		},
	})
}
