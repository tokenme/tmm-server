package operation

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func HeadlinesHandler(c *gin.Context) {
    db := Service.Db
    query := `
        SELECT h.id,
            h.online_status,
            h.title,
            h.link
        FROM tmm.headlines AS h
        WHERE h.online_status = 1
    `
    rows, _, err := db.Query(query)
    if CheckErr(err, c) {
		return
	}
    var headlines []common.Headline
    for _, row := range rows {
        headline := common.Headline {
            Id:             row.Uint64(0),
            OnlineStatus:   int8(row.Int(1)),
            Title:          row.Str(2),
            Link:           row.Str(3),
        }
        headlines = append(headlines, headline)
    }
    c.JSON(http.StatusOK, headlines)
}
