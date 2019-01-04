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
            OnlineStatus:   int8(row.Int(3)),
            Title:          row.Str(4),
            Link:           row.Str(5),
            Width:          row.Uint(6),
            Height:         row.Uint(7),
            Action:         row.Str(8),
        }
        activities = append(activities, activity)
    }
    c.JSON(http.StatusOK, activities)
}
