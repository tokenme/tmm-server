package article

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"net/http"
	. "github.com/tokenme/tmm/handler"
	"time"
)

func GetArticleHandler(c *gin.Context) {
	db := Service.Db
	var (
		offset, count   int
		query, sumquery string
		articleList     []Article
	)
	sortid, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "3"))
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}

	if sortid == 0 {
		query = fmt.Sprintf(`
	select 	id,fileid,author,
	title,link,source_url,cover,published_at,
	digest,content,sortid,published from tmm.articles order by id DESC
	limit %d offset %d
	`, limit, offset)
		sumquery = `select count(*) from tmm.articles`
	} else {
		query = fmt.Sprintf(`
		select 	id,fileid,author,
		title,link,source_url,cover,published_at,
		digest,content,sortid,published from tmm.articles
		order by id DESC where sortid = %d  limit %d offset %d`, sortid, limit, offset)
		sumquery = fmt.Sprintf(`select count(*) from tmm.articles  where sortid = %d `, sortid)
	}

	rows, result, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "没有到数据",
			"data": gin.H{
				"curr_page": page,
				"data":      "",
			},
		})
		return
	}

	for _, row := range rows {
		Link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, row.Int(result.Map(`id`)))
		query = `select online_status from tmm.share_tasks where link = '%s'`
		a, result, err := db.Query(query, Link)
		if CheckErr(err, c) {
			return
		}
		article := Article{
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
	count = rows[0].Int(0)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": gin.H{
			"curr_page": page,
			"data":      articleList,
			"amount":    count,
		},
	}, )

}
