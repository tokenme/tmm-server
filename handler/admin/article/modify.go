package article

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	. "github.com/tokenme/tmm/handler"
)

func ModifyArticleHandler(c *gin.Context) {
	db := Service.Db
	var up = Article{}
	if CheckErr(c.Bind(&up), c) {
		return
	}
	query := `update tmm.articles as art ,tmm.share_tasks as task
	set art.fileid = %d ,art.author='%s',art.title='%s',art.link='%s',art.source_url='%s',
	art.content='%s' ,art.digest = '%s',art.cover='%s',
	art.published_at='%s',art.published = %d ,task.summary='%s',task.title='%s',task.image='%s'
	where art.id = %d and task.link = '%s' `
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, up.Id)

	_, _, err := db.Query(query, up.Fileid, db.Escape(up.Author), db.Escape(up.Title),
		db.Escape(up.Link), db.Escape(up.SourceUrl),
		db.Escape(up.Content), db.Escape(up.Digest), db.Escape(up.Cover),
		db.Escape(up.PublishedOn), up.Published, db.Escape(up.Digest), db.Escape(up.Title),
		db.Escape(up.Cover), up.Id, link)

	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: `ok`})
}
