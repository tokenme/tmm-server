package ad

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "github.com/tokenme/tmm/handler"
	"time"
	"fmt"
	"strings"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func CreativeListHandler(c *gin.Context) {
	db := Service.Db
	online, err := strconv.Atoi(c.DefaultQuery(`online`, "0"))
	if CheckErr(err, c) {
		return
	}
	adzoneId, err := strconv.Atoi(c.DefaultQuery(`adzoneId`, "-1"))
	if CheckErr(err, c) {
		return
	}
	adGroupId, err := strconv.Atoi(c.DefaultQuery(`groupId`, "-1"))
	if CheckErr(err, c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, "0"))
	if CheckErr(err, c) {
		return
	}
	channelId, err := strconv.Atoi(c.DefaultQuery(`channelid`, "-1"))
	if CheckErr(err, c) {
		return
	}
	when := c.DefaultQuery(`when`, "")
	if when == "" {
		when = time.Now().AddDate(0, 0, -7).String()
	}
	var where []string
	if online == 1 || online == -1 {
		where = append(where, fmt.Sprintf(" AND creat.online_status = %d ", online))
	}
	if adzoneId != -1 {
		where = append(where, fmt.Sprintf(" AND adz.id = %d ", adzoneId))
	}
	if adGroupId != -1 {
		where = append(where, fmt.Sprintf(" AND creat.adgroup_id = %d ", adGroupId))
	}
	if channelId != -1 {
		where = append(where, fmt.Sprintf(" AND adz.cid = %d ", channelId))
	}
	var limit, offset int
	if page == 0 {
		limit = 10
		offset = 0
	} else {
		limit = 10
		offset = 10 * page
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
	stats.cik AS cik,
	stats.imp AS imp,
	adz.id AS adz_id,
	adz.summary AS summary,
	group_.id AS group_id,
	group_.title AS group_title
FROM tmm.creatives  AS creat 
LEFT JOIN (SELECT 
		   SUM(imp) AS imp ,
           SUM(cik) AS cik,
		   creative_id AS id FROM tmm.creative_stats
		   WHERE record_on > '%s'
		   GROUP BY id) AS stats 
		   ON (stats.id = creat.id )
INNER JOIN tmm.adgroups AS group_ ON (group_.id = creat.adgroup_id)
INNER JOIN tmm.adzones AS  adz ON (adz.id = group_.adzone_id)
WHERE 1 = 1 %s
ORDER BY adz.id,stats.imp DESC
LIMIT  %d 
OFFSET %d `

	rows, res, err := db.Query(query, db.Escape(when),
		strings.Join(where, " "), limit, offset)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	var list []*Adzone
	for _, row := range rows {
		adzone := &Adzone{}
		adzone.Id = row.Uint64(res.Map(`adz_id`))
		adzone.Summery = row.Str(res.Map(`summary`))
		adzone.Group.Title = row.Str(res.Map(`group_title`))
		adzone.Group.Id = row.Uint64(res.Map(`group_id`))
		adzone.Group.Creatives = Creatives{
			InsertAt:     row.Str(res.Map(`inserted_at`)),
			StartDate:    row.Str(res.Map(`start_date`)),
			EndDate:      row.Str(res.Map(`end_date`)),
			UpdateAt:     row.Str(res.Map(`updated`)),
			OnlineStatus: row.Int(res.Map(`online_status`)),
			Platform:     row.Int(res.Map(`platform`)),
		}
		adzone.Group.Creatives.Img = row.Int(res.Map(`imp`))
		adzone.Group.Creatives.Click = row.Int(res.Map(`cik`))
		adzone.Group.Creatives.Id = row.Uint64(res.Map(`id`))
		adzone.Group.Creatives.Title = row.Str(res.Map(`title`))
		adzone.Group.Creatives.Image = row.Str(res.Map(`image`))
		adzone.Group.Creatives.Link = row.Str(res.Map(`link`))
		list = append(list, adzone)
	}

	groupImpClick := make(map[uint64]*Data)
	adzoneImpClick := make(map[uint64]*Data)
	for _, adzone := range list {
		creat := adzone.Group.Creatives
		if data, ok := groupImpClick[adzone.Group.Id]; ok {
			data.Click += creat.Click
			data.Img += creat.Img
			groupImpClick[adzone.Group.Id] = data
		} else {
			groupImpClick[adzone.Group.Id] = &Data{
				Click: creat.Click,
				Img:   creat.Img,
			}
		}
		if data, ok := adzoneImpClick[adzone.Id]; ok {
			data.Click += creat.Click
			data.Img += creat.Img
			adzoneImpClick[adzone.Id] = data
		} else {
			adzoneImpClick[adzone.Id] = &Data{
				Click: creat.Click,
				Img:   creat.Img,
			}
		}
	}

	for _, adzone := range list {
		adzone.Img = adzoneImpClick[adzone.Id].Img
		adzone.Click = adzoneImpClick[adzone.Id].Click
		adzone.Group.Click = groupImpClick[adzone.Group.Id].Click
		adzone.Group.Img = groupImpClick[adzone.Group.Id].Img
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
		},
	})

}
