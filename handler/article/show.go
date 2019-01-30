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

const (
	ARTICLE_CACKE_KEY = "article-%s"
)

type Article struct {
	Author      string        `json:"author,omitempty"`
	Title       string        `json:"title,omitempty"`
	SourceUrl   string        `json:"source_url,omitempty"`
	Content     template.HTML `json:"-"`
	RawContent  string        `json:"raw,omitempty"`
	Digest      string        `json:"digest,omitempty"`
	Cover       string        `json:"cover,omitempty"`
	PublishedOn string        `json:"published_on,omitempty"`
}

func ShowHandler(c *gin.Context) {
	articleId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT author, title, source_url, digest, content, cover, published_at, inserted_at FROM tmm.articles WHERE id=%d`, articleId)
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
		RawContent:  row.Str(4),
		Cover:       row.Str(5),
		PublishedOn: row.ForceLocaltime(6).Format("2006-01-02"),
	}
	if article.PublishedOn == "1970-01-01" {
		article.PublishedOn = row.ForceLocaltime(7).Format("2006-01-02")
	}
	article.SourceUrl = "https://tmm.tianxi100.com/article/rand"
	article.Content = template.HTML(article.RawContent)
	c.HTML(http.StatusOK, "article.tmpl", article)
}
