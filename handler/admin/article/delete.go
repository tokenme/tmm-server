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
	_, _, err = db.Query(`delete tmm.share_tasks  ,tmm.share_task_categories,tmm.articles  as art
	from tmm.share_tasks,tmm.share_task_categories,tmm.articles
	 where tmm.share_tasks.link = '%s'
	 and  tmm.share_task_categories.task_id = tmm.share_tasks.id and art.id =%d
		`, db.Escape(link),articleId)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": ""})

}