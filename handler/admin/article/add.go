package article

import (
	"github.com/gin-gonic/gin"
	"fmt"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func AddArticleHandler(c *gin.Context) {
	db := Service.Db
	var article Article
	if CheckErr(c.Bind(&article), c) {
		return
	}
	query := `INSERT INTO tmm.articles(fileid,author,title,link,source_url,
cover,published_at,digest,content,sortid,Published)
VALUES (%d,'%s','%s','%s','%s','%s','%s','%s','%s',%d,%d)`
	_, result, err := db.Query(query, article.Fileid, db.Escape(article.Author), db.Escape(article.Title),
		db.Escape(article.Link), db.Escape(article.SourceUrl), db.Escape(article.Cover),
		db.Escape(article.PublishedOn), db.Escape(article.Digest), db.Escape(article.Content),
		article.Sortid, 1)
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf("https://tmm.tokenmama.io/article/show/%d", result.InsertId())

	query = `INSERT INTO tmm.share_tasks (creator, title, summary, link,
	image, points, points_left, bonus, max_viewers,online_status) VALUES(0, '%s', '%s', '%s', '%s', 5000, 5000, 10, 10,-1)`

	_, _, err = db.Query(query, db.Escape(article.Title), db.Escape(article.Digest),
		db.Escape(link), db.Escape(article.Cover))
	if CheckErr(err, c) {
		return
	}
	rows, _, err := db.Query(`select id from tmm.share_tasks where link = '%s'`, db.Escape(link))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `select Error `, c) {
		return
	}
	_, _, err = db.Query(`INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUES (%d,%d,%d)`, rows[0].Int(0), article.Sortid, 0)
	if CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "ok",
		"data": article},
	)
}
