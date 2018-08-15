package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"time"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

func SharesHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	page, _ := strconv.ParseUint(c.Param("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(c.Param("pageSize"), 10, 64)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 || pageSize > DEFAULT_PAGE_SIZE {
		pageSize = DEFAULT_PAGE_SIZE
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
    st.points_left,
    st.inserted_at,
    st.updated_at
FROM share_tasks AS st
WHERE st.points_left > 0
ORDER BY st.bonus DESC, st.id DESC LIMIT %d, %d`
	rows, _, err := db.Query(query, (page-1)*pageSize, pageSize)
	if CheckErr(err, c) {
		return
	}
	var tasks []common.ShareTask
	for _, row := range rows {
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
			InsertedAt: row.ForceLocaltime(9).Format(time.RFC3339),
			UpdatedAt:  row.ForceLocaltime(10).Format(time.RFC3339),
		}
		task.ShareLink, _ = task.GetShareLink(user.Id, Config)
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, tasks)
}
