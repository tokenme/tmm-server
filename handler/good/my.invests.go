package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	//"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"time"
)

func MyInvestsHandler(c *gin.Context) {
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
	defaultPageSize := uint64(DEFAULT_PAGE_SIZE)
	if pageSize == 0 || pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT
            g.id, g.name,
            g.pic,
            gi.points,
            SUM(IFNULL(tx.income, 0)),
            gi.redeem_status,
            gi.inserted_at
        FROM tmm.good_invests AS gi
        INNER JOIN tmm.goods AS g ON (g.id=gi.good_id)
        LEFT JOIN tmm.good_tx AS tx ON (tx.good_id=gi.good_id AND tx.created_at>=gi.inserted_at)
        WHERE
            gi.user_id=%d
        GROUP BY gi.good_id
        ORDER BY gi.inserted_at DESC LIMIT %d, %d`, user.Id, (page-1)*pageSize, pageSize)
	if CheckErr(err, c) {
		return
	}
	var invests []common.GoodInvest
	for _, row := range rows {
		points, _ := decimal.NewFromString(row.Str(3))
		income, _ := decimal.NewFromString(row.Str(4))
		inv := common.GoodInvest{
			GoodId:       row.Uint64(0),
			GoodName:     row.Str(1),
			GoodPic:      row.Str(2),
			Points:       points,
			Income:       income,
			RedeemStatus: row.Uint(5),
			InsertedAt:   row.ForceLocaltime(6).Format(time.RFC3339),
		}
		invests = append(invests, inv)
	}
	c.JSON(http.StatusOK, invests)
}
