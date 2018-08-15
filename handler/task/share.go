package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
)

func ShareHandler(c *gin.Context) {

	encryptedTaskId := c.Param("encryptedTaskId")
	encryptedUserId := c.Param("encryptedUserId")
	taskId, err := utils.DecryptUint64(encryptedTaskId, []byte(Config.LinkSalt))
	if CheckErr(err, c) {
		return
	}
	userId, err := utils.DecryptUint64(encryptedUserId, []byte(Config.LinkSalt))
	if CheckErr(err, c) {
		return
	}
	if Check(taskId == 0 || userId == 0, "not found", c) {
		return
	}
	db := Service.Db
	query := `SELECT
    st.id,
    st.title,
    st.summary,
    st.link,
    st.image,
    st.max_viewers,
    st.bonus,
    st.points,
    st.points_left
    ust.points,
    ust.viewers
FROM share_tasks AS st
LEFT JOIN user_share_tasks AS ust ON (ust.task_id=st.id AD ust.user_id=%d)
WHERE st.id=%d
LIMIT 1`
	rows, _, err := db.Query(query, userId, taskId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		return
	}
	row := rows[0]
	bonus, _ := decimal.NewFromString(row.Str(6))
	points, _ := decimal.NewFromString(row.Str(7))
	pointsLeft, _ := decimal.NewFromString(row.Str(8))
	task := common.ShareTask{
		Id:         row.Uint64(0),
		Title:      row.Str(1),
		Summary:    row.Str(2),
		Link:       row.Str(3),
		Image:      row.Str(4),
		MaxViewers: row.Uint(5),
		Bonus:      bonus,
		Points:     points,
		PointsLeft: pointsLeft,
	}
	userViewers := row.Uint(9)
	if _, err := c.Cookie(task.CookieKey()); err == nil && (task.PointsLeft.GreaterThanOrEqual(bonus) && task.MaxViewers > userViewers) {
		_, _, err := db.Query(`INSERT INTO tmm.user_share_tasks (user_id, task_id, points) VALUES (%d, %d, %s) ON DUPLICATE KEY UPDATE points=points+VALUES(points), viewers=viewers+1`, task.Id, userId, bonus.StringFixed(9))
		if err == nil {
			_, _, err = db.Query(`UPDATE tmm.share_tasks SET points_left=points_left-bonus, viewers=viewers+1 WHERE id=%d`, task.Id)
			if err != nil {
				log.Error(err.Error())
			}
			c.SetCookie(task.CookieKey(), "1", 86400, "/", Config.Domain, true, true)
		} else {
			log.Error(err.Error())
		}
	}
	c.HTML(http.StatusOK, "share.tmpl", task)
}
