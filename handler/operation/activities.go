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
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
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
          AND NOW() >= a.start_date 
          AND NOW() <= a.end_date
          AND IF (a.active_days = 0 OR (
              SELECT 1
              FROM tmm.devices AS d
              WHERE (
                EXISTS (SELECT 1 FROM tmm.device_app_tasks AS dat WHERE dat.device_id = d.id AND dat.inserted_at >= DATE_SUB(NOW(), INTERVAL a.active_days DAY) LIMIT 1)
                OR
                EXISTS (SELECT 1 FROM tmm.device_share_tasks AS dst WHERE dst.device_id = d.id AND dst.inserted_at >= DATE_SUB(NOW(), INTERVAL a.active_days DAY) LIMIT 1)
                OR
                EXISTS (SELECT 1 FROM tmm.reading_logs AS rl WHERE rl.user_id = d.user_id AND (rl.inserted_at >= DATE_SUB(NOW(), INTERVAL a.active_days DAY) OR rl.updated_at >= DATE_SUB(NOW(), INTERVAL a.active_days DAY)) LIMIT 1)
                OR
                EXISTS (SELECT 1 FROM tmm.daily_bonus_logs AS dbl WHERE dbl.user_id = d.user_id AND dbl.updated_on >= DATE_SUB(NOW(), INTERVAL a.active_days DAY) LIMIT 1)
              )
              AND d.user_id = %d
              LIMIT 1
            ), 1, 0)
          AND IF(a.registered_on IS NULL OR (
              a.registered_on <= (SELECT DATE(u.created) FROM ucoin.users AS u WHERE u.id = %d)
            ), 1, 0)
        ORDER BY a.row_id
    `
    rows, _, err := db.Query(query, user.Id, user.Id)
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
