package article

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"html/template"
	"net/http"
	"strconv"
)

type Article struct {
	Author      string
	Title       string
	SourceUrl   string
	Content     template.HTML
	Digest      string
	Cover       string
	PublishedOn string
}

func ShowHandler(c *gin.Context) {
	articleId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT author, title, source_url, digest, content, cover, published_at FROM tmm.articles WHERE id=%d`, articleId)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	row := rows[0]
	article := Article{
		Author:      row.Str(0),
		Title:       row.Str(1),
		SourceUrl:   row.Str(2),
		Digest:      row.Str(3),
		Content:     template.HTML(row.Str(4)),
		Cover:       row.Str(5),
		PublishedOn: row.ForceLocaltime(6).Format("2006-01-02"),
	}
	c.HTML(http.StatusOK, "article.tmpl", article)
}
