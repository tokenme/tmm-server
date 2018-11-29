package article

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"net/http"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
	"time"
)

func GetArticleListHandler(c *gin.Context) {
	db := Service.Db
	sortid, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "3"))
	var offset int
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}

	var query, sumquery string
	query = `
	SELECT 	
	id,
	fileid,
	author,
	title,
	link,
	source_url,
	cover,
	published_at,			
	digest,
	content,
	sortid,
	published 
	FROM tmm.articles 
	ORDER BY id DESC
	%s	LIMIT %d OFFSET %d`
	if sortid == 0 {
		query = fmt.Sprintf(query, " ", limit, offset)
		sumquery = `select count(*) FROM tmm.articles`
	} else {
		query = fmt.Sprintf(query, fmt.Sprintf(`WHERE sortid = %d`, sortid), limit, offset)
		sumquery = fmt.Sprintf(`SELECT count(*) FROM tmm.articles  WHERE sortid = %d `, sortid)
	}

	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK,
			admin.Response{
				Code:    1,
				Message: "没有到数据",
				Data: gin.H{
					"curr_page": page,
					"data":      "",
				},
			})
		return
	}

	var articleList []*Article
	for _, row := range rows {
		Link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, row.Int(res.Map(`id`)))
		query = `SELECT online_status FROM tmm.share_tasks WHERE link = '%s'`
		a, result, err := db.Query(query, Link)
		if CheckErr(err, c) {
			return
		}
		article := &Article{
			Id:          row.Int(result.Map(`id`)),
			Fileid:      row.Int(result.Map(`fileid`)),
			Author:      row.Str(result.Map(`author`)),
			Link:        row.Str(result.Map(`link`)),
			SourceUrl:   row.Str(result.Map(`source_url`)),
			Cover:       row.Str(result.Map(`cover`)),
			PublishedOn: row.ForceLocaltime(result.Map(`published_at`)).Format(time.RFC3339),
			Digest:      row.Str(result.Map(`digest`)),
			Content:     row.Str(result.Map(`content`)),
			Sortid:      row.Int(result.Map(`sortid`)),
			Published:   row.Int(result.Map(`published`)),
		}
		if len(a) == 0 {
			continue
		} else {
			article.Online = a[0].Int(result.Map(`online_status`))
		}
		articleList = append(articleList, article)
	}
	rows, _, err = db.Query(sumquery)
	if CheckErr(err, c) {
		return
	}
	count := rows[0].Int(0)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data: gin.H{
			"curr_page": page,
			"data":      articleList,
			"amount":    count,
		},
	}, )
}
