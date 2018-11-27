package admin

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"time"
	"net/http"
	"github.com/shopspring/decimal"
	"strconv"
	"fmt"
)

type ShareAddRequest struct {
	Title         string `json:"title" form:"title" binding:"required"`
	Summary       string `json:"summary" form:"summary" binding:"required"`
	Link          string `json:"link" form:"link" binding:"required"`
	Image         string `json:"image" form:"image"`
	FileExtension string `json:"image_extension" from:"image_extension"`
	Points        string `json:"points" form:"points" binding:"required"`
	Bonus         string `json:"bonus" form:"bonus" binding:"required"`
	MaxViewers    uint   `json:"max_viewers" form:"max_viewers" binding:"required"`
}

func AddShareHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	var req ShareAddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}

	db := Service.Db
	query := `SELECT
d.id
FROM tmm.devices AS d
WHERE d.user_id = %d ORDER BY d.points DESC LIMIT 1`
	rows, _, err := db.Query(query, user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	deviceId := rows[0].Str(0)

	_, ret, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, points, points_left, bonus, max_viewers) VALUES (%d, '%s', '%s', '%s', '%s', %s, %s, %s, %d)`, user.Id, db.Escape(req.Title), db.Escape(req.Summary), db.Escape(req.Link), db.Escape(req.Image), db.Escape(req.Points), db.Escape(req.Points), db.Escape(req.Bonus), req.MaxViewers)
	if CheckErr(err, c) {
		return
	}
	now := time.Now().Format(time.RFC3339)
	Points, err := decimal.NewFromString(req.Points)
	if CheckErr(err, c) {
		return
	}
	Bonus, err := decimal.NewFromString(req.Bonus)
	if CheckErr(err, c) {
		return
	}
	task_ := common.ShareTask{
		Id:         ret.InsertId(),
		Title:      req.Title,
		Summary:    req.Summary,
		Link:       req.Link,
		Image:      req.Image,
		Points:     Points,
		PointsLeft: Points,
		Bonus:      Bonus,
		MaxViewers: req.MaxViewers,
		InsertedAt: now,
		UpdatedAt:  now,
		Creator:    user.Id,
	}
	task_.ShareLink, _ = task_.GetShareLink(deviceId, Config)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": task_})
}

func GetShareListHandler(c *gin.Context) {
	db := Service.Db
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "5"))
	types, _ := strconv.Atoi(c.DefaultQuery(`type`, "0"))
	var (
		offset, count int
		Query         string
	)
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}
	if types == 0 {
		Query = fmt.Sprintf(`select id,creator,title,summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,inserted_at,updated_at from tmm.share_tasks order by inserted_at,updated_at DESC limit %d offset %d`, limit, offset)
	} else {
		Query = fmt.Sprintf(`select id,creator,title,summary,link,image,points,points_left,
	bonus,max_viewers,viewers,online_status,inserted_at,updated_at from tmm.share_tasks LEFT JOIN 
	share_task_categories ON(id = task_id) where cid = %d order by inserted_at,updated_at DESC 
	limit %d offset %d`, types, limit, offset)
	}
	Rows, result, err := db.Query(Query)
	if CheckErr(err, c) {
		return
	}
	var sharelist []common.ShareTask
	for _, row := range Rows {
		Points := decimal.NewFromFloat(row.Float(result.Map(`points`)))
		Bonus := decimal.NewFromFloat(row.Float(result.Map(`bonus`)))
		points_left := decimal.NewFromFloat(row.Float(result.Map(`points_left`)))
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
			Points:       Points,
			PointsLeft:   points_left,
			Bonus:        Bonus,
			MaxViewers:   row.Uint(result.Map(`max_viewers`)),
			Viewers:      row.Uint(result.Map(`viewers`)),
			OnlineStatus: int8(row.Int(result.Map(`online_status`))),
			InsertedAt:   row.Str(result.Map(`inserted_at`)),
			UpdatedAt:    row.Str(result.Map(`updated_at`)),
		}
		sharelist = append(sharelist, share)
	}
	if types == 0 {
		Rows, _, err := db.Query(`select count(*) from tmm.share_tasks`)
		if err != nil {
			return
		}
		count = Rows[0].Int(0)

	} else {
		Rows, _, err := db.Query(`select count(*) from tmm.share_tasks LEFT JOIN 
	    share_task_categories ON(id = task_id) where cid = %d`, types)
		if err != nil {
			return
		}
		count = Rows[0].Int(0)
	}
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
