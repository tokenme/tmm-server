package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"

	"strconv"
	"fmt"
	"time"
	"github.com/tokenme/tmm/handler/admin"
)

func GetArticleHandler(c *gin.Context) {
	db := Service.Db
	var (
		up          = Article{}
		query       string
		onlineQuery = `SELECT online_status FROM tmm.share_tasks WHERE link = '%s'`
		link        string
	)
	artId, err := strconv.Atoi(c.Query(`artId`))
	if CheckErr(err, c) {
		return
	}
	link = fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, artId)

	query = `SELECT 
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
	published FROM tmm.articles WHERE id = %d`
	rows, result, err := db.Query(query, artId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `Not found`, c) {
		return
	}
	row := rows[0]
	up = Article{
		Id:          artId,
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
	rows, _, err = db.Query(onlineQuery, link)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0,`Not found Online field`,c){
		return
	}
	up.Online = rows[0].Int(0)
		c.JSON(http.StatusOK, admin.Response{
			Message:  admin.API_OK,
			Data: up,
		})
}
