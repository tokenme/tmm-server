package article

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func DeleteArticleHandler(c *gin.Context) {
	db := Service.Db
	articleId, err := strconv.Atoi(c.Param("id"))
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, articleId)
	_, _, err = db.Query(`DELETE tmm.share_tasks  ,tmm.share_task_categories,tmm.articles  AS art
	from tmm.share_tasks,tmm.share_task_categories,tmm.articles
	 WHERE tmm.share_tasks.link = '%s'
	 AND  tmm.share_task_categories.task_id = tmm.share_tasks.id AND art.id =%d
		`, db.Escape(link), articleId)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: `ok`})

}
