package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/handler/admin"
)

func CategoryListHandler(c *gin.Context) {
	db := Service.Db
	query := `SELECT id,name FROM tmm.article_categories ORDER BY id`
	maps := make(map[int]string)
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	for _, row := range rows {
		maps[row.Int(0)] = row.Str(1)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    maps,
	},
	)

}
