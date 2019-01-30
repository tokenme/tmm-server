package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func RandHandler(c *gin.Context) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT r1.id
FROM tmm.articles AS r1
JOIN (SELECT CEIL(RAND() * (SELECT MAX(id) FROM tmm.articles)) AS id) AS r2
WHERE r1.id >= r2.id
ORDER BY r1.id ASC
LIMIT 1`)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	row := rows[0]
	articleId := row.Uint64(0)
	link := fmt.Sprintf("https://tmm.tianxi100.com/article/show/%d", articleId)
	c.Redirect(http.StatusFound, link)
}
