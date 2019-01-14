package operation

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func ActivitiesHandler(c *gin.Context) {
    db := Service.Db
    query := `
        SELECT a.id,
            a.row_id,
            a.image,
            a.share_image,
            a.online_status,
            a.title,
            a.link,
            a.width,
            a.height,
            a.action
        FROM tmm.activities AS a
        WHERE a.online_status = 1
          AND a.row_id > 0
        ORDER BY a.row_id
    `
    rows, _, err := db.Query(query)
    if CheckErr(err, c) {
		return
	}
    var activities []common.Activity
    for _, row := range rows {
        activity := common.Activity {
            Id:             row.Uint64(0),
            RowId:          row.Uint(1),
            Image:          row.Str(2),
            ShareImage:     row.Str(3),
            OnlineStatus:   int8(row.Int(4)),
            Title:          row.Str(5),
            Link:           row.Str(6),
            Width:          row.Uint(7),
            Height:         row.Uint(8),
            Action:         row.Str(9),
        }
        activities = append(activities, activity)
    }
    c.JSON(http.StatusOK, activities)
}
