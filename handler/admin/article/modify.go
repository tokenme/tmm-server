package article

import (
	"github.com/gin-gonic/gin"
	"fmt"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
)

func ModifyArticleHandler(c *gin.Context) {
	db := Service.Db
	var up = Article{}
	if CheckErr(c.Bind(&up), c) {
		return
	}
	var where []string
	if up.Fileid != 0 {
		where = append(where, fmt.Sprintf(`art.fileid = %d`, up.Fileid))
	}
	if up.Author != "" {
		where = append(where, fmt.Sprintf(`art.author='%s'`, db.Escape(up.Author)))
	}
	if up.Title != "" {
		where = append(where, fmt.Sprintf(`art.title='%s'`, db.Escape(up.Title)))
		where = append(where, fmt.Sprintf(`task.title='%s'`, db.Escape(up.Title)))
	}
	if up.Link != "" {
		where = append(where, fmt.Sprintf(`art.link='%s'`, db.Escape(up.Link)))
	}
	if up.SourceUrl != "" {
		where = append(where, fmt.Sprintf(`art.source_url='%s'`, db.Escape(up.SourceUrl)))
	}
	if up.Content != "" {
		where = append(where, fmt.Sprintf(`art.content='%s'`, db.Escape(up.Content)))

	}
	if up.Digest != "" {
		where = append(where, fmt.Sprintf(`art.digest='%s'`, db.Escape(up.Digest)))
		where = append(where, fmt.Sprintf(`task.summary='%s'`, db.Escape(up.Digest)))
	}
	if up.Cover != "" {
		where = append(where, fmt.Sprintf(`art.cover='%s'`, db.Escape(up.Cover)))
		where = append(where, fmt.Sprintf(`task.image='%s'`, db.Escape(up.Cover)))
	}
	if up.PublishedOn != "" {
		where = append(where, fmt.Sprintf(`art.published_at='%s'`, db.Escape(up.PublishedOn)))
	}
	if up.Published != 0 {
		where = append(where, fmt.Sprintf(`art.published=%d`, up.PublishedOn))
	}

	query := `
	UPDATE tmm.articles AS art ,
	tmm.share_tasks AS task
	%s
	WHERE art.id = %d AND task.link = '%s' `
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, up.Id)

	if _, _, err := db.Query(query, strings.Join(where, `,`), up.Id, link); CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: `ok`})
}
