package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ykt"
	"net/http"
	"strconv"
)

func ItemHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	itemId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	yktReq := ykt.GoodInfoRequest{
		Id:  itemId,
		Uid: user.Id,
	}
	res, err := yktReq.Run()
	if CheckErr(err, c) {
		return
	}
	good := res.Data.Data
	good.CommissionPoints = decimal.New(Config.GoodCommissionPoints, 0)
	db := Service.Db
	rows, _, err := db.Query(`SELECT points FROM tmm.good_invests WHERE good_id=%d AND user_id=%d LIMIT 1`, good.Id, user.Id)
	if err == nil && len(rows) > 0 {
		row := rows[0]
		points, _ := decimal.NewFromString(row.Str(0))
		good.InvestPoints = points
	} else if err != nil {
		log.Error(err.Error())
	}
	rows, _, err = db.Query(`SELECT SUM(points), COUNT(*) FROM tmm.good_invests WHERE good_id=%d`, good.Id)
	if err == nil && len(rows) > 0 {
		row := rows[0]
		points, _ := decimal.NewFromString(row.Str(0))
		good.TotalInvest = points
		good.TotalInvestors = row.Uint(1)
	} else if err != nil {
		log.Error(err.Error())
	}
	c.JSON(http.StatusOK, good)
}
