package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"
)

func CategoryListHandler(c *gin.Context) {
	db := Service.Db
	query := `select id,name from tmm.article_categories`
	maps := make(map[int]string)
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	for _, row := range rows {
		maps[row.Int(0)] = row.Str(1)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": maps},
	)

}
