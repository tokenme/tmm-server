package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"time"
	"fmt"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

type ReturnData struct {
	AdGroup   AdGroup   `json:"ad_group"`
	Creatives Creatives `json:"creatives"`
	Adzone    Adzone    `json:"adzone"`
}

type SearchOptions struct {
	Online   int    `form:"online"`
	AdzoneId int    `form:"adzoneId"`
	GroupId  int    `form:"groupId"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Title    string `form:"title"`
}

func CreativeListHandler(c *gin.Context) {
	db := Service.Db
	var req SearchOptions
	if CheckErr(c.Bind(&req), c) {
		return
	}
	ChannelId := c.DefaultQuery(`channelid`, `-1`)
	when := c.DefaultQuery(`when`, "")
	if when == "" {
		when = time.Now().AddDate(0, 0, -7).String()
	}
	var where []string
	if req.Online == 1 || req.Online == -1 {
		where = append(where, fmt.Sprintf(" AND creat.online_status = %d ", req.Online))
	}
	if req.AdzoneId != 0 {
		where = append(where, fmt.Sprintf(" AND adz.id = %d ", req.AdzoneId))
	}
	if req.GroupId != 0 {
		where = append(where, fmt.Sprintf(" AND creat.adgroup_id = %d ", req.GroupId))
	}
	if ChannelId != "-1" {
		where = append(where, fmt.Sprintf(" AND adz.cid = %s ", db.Escape(ChannelId)))
	}
	if req.Title != "" {
		where = append(where, fmt.Sprintf(" AND title like '%s%s'", db.Escape(req.Title), "%"))
	}
	var offset, limit int
	if req.PageSize > 0 {
		limit = req.PageSize
	} else {
		limit = 25
	}
	if req.Page <= 0 {
		offset = 0
	} else {
		offset = limit * (req.Page - 1)
	}

	query := `
SELECT 
	creat.title AS title,
	creat.image AS image,
	creat.inserted_at AS inserted_at,
	creat.start_date AS start_date,
	creat.end_date   AS end_date,
	creat.id AS id ,
	creat.online_status AS online_status,
	creat.updated_at AS updated,
	creat.link AS link,
	creat.platform AS platform,
	creat.ad_mode AS mode,
	creat.ad_income AS income,
	stats.clk AS clk,
	stats.imp AS imp,
	adz.id AS adz_id,
	adz.summary AS summary,
	group_.id AS group_id,
	group_.title AS group_title
FROM tmm.creatives  AS creat 
LEFT JOIN (SELECT 
		   SUM(imp) AS imp,
           SUM(clk) AS clk,
		   creative_id AS id FROM tmm.creative_stats
		   WHERE record_on > '%s'
		   GROUP BY id) AS stats 
		   ON (stats.id = creat.id )
INNER JOIN tmm.adgroups AS group_ ON (group_.id = creat.adgroup_id)
INNER JOIN tmm.adzones AS  adz ON (adz.id = group_.adzone_id)
WHERE 1 = 1 %s 
ORDER BY adz.id DESC,stats.imp DESC
LIMIT  %d 
OFFSET %d `

	rows, res, err := db.Query(query, db.Escape(when),
		strings.Join(where, " "), limit, offset)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, admin.Response{
			Code:    0,
			Message: "没有找到数据",
			Data: gin.H{
				"total": 0,
				"data":  nil,
			},
		})
		return
	}
	var list []*Creatives
	for _, row := range rows {
		creatives := &Creatives{
			InsertAt:     row.Str(res.Map(`inserted_at`)),
			StartDate:    row.Str(res.Map(`start_date`)),
			EndDate:      row.Str(res.Map(`end_date`)),
			UpdateAt:     row.Str(res.Map(`updated`)),
			OnlineStatus: row.Int(res.Map(`online_status`)),
			Platform:     row.Int(res.Map(`platform`)),
			AdMode:       row.Str(res.Map(`mode`)),
			AdIncome:     fmt.Sprintf("%.2f", row.Float(res.Map(`income`))),
		}
		creatives.Imp = row.Int(res.Map(`imp`))
		creatives.Click = row.Int(res.Map(`clk`))
		creatives.Id = row.Uint64(res.Map(`id`))
		creatives.Title = row.Str(res.Map(`title`))
		creatives.Image = row.Str(res.Map(`image`))
		creatives.Link = row.Str(res.Map(`link`))
		creatives.Adzone.Id = row.Uint64(res.Map(`adz_id`))
		creatives.Adzone.Summery = row.Str(res.Map(`summary`))
		creatives.AdGroup.Title = row.Str(res.Map(`group_title`))
		creatives.AdGroup.Id = row.Uint64(res.Map(`group_id`))
		list = append(list, creatives)
	}

	groupImpClick := make(map[uint64]*Data)
	adzoneImpClick := make(map[uint64]*Data)
	for _, creatives := range list {
		creat := creatives
		if data, ok := groupImpClick[creat.AdGroup.Id]; ok {
			data.Click += creat.Click
			data.Imp += creat.Imp
			groupImpClick[creat.AdGroup.Id] = data
		} else {
			groupImpClick[creat.AdGroup.Id] = &Data{
				Click: creat.Click,
				Imp:   creat.Imp,
			}
		}
		if data, ok := adzoneImpClick[creat.Adzone.Id]; ok {
			data.Click += creat.Click
			data.Imp += creat.Imp
			adzoneImpClick[creat.Adzone.Id] = data
		} else {
			adzoneImpClick[creat.Adzone.Id] = &Data{
				Click: creat.Click,
				Imp:   creat.Imp,
			}
		}
	}

	for _, creat := range list {
		creat.Adzone.Imp = adzoneImpClick[creat.Adzone.Id].Imp
		creat.Adzone.Click = adzoneImpClick[creat.Adzone.Id].Click
		creat.AdGroup.Click = groupImpClick[creat.AdGroup.Id].Click
		creat.AdGroup.Imp = groupImpClick[creat.AdGroup.Id].Imp
	}
	rows, _, err = db.Query(`SELECT COUNT(1) FROM tmm.creatives AS creat  
	INNER JOIN tmm.adgroups AS adg ON (adg.id = creat.adgroup_id)
	INNER JOIN tmm.adzones AS  adz ON (adz.id = adg.adzone_id)
	WHERE 1 = 1 %s `, strings.Join(where, " "))
	if CheckErr(err, c) {
		return
	}
	var total int
	if len(rows) != 0 {
		total = rows[0].Int(0)
	}
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"total": total,
			"data":  list,
			"page":  req.Page,
		},
	})

}
