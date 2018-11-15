package good

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strings"
	"time"
)

type GoodTxBonus struct {
	GoodId uint64
	UserId uint64
	Amount uint
}

func TxsHandler(c *gin.Context) {
	data := c.PostForm("data")
	txs, err := common.DecodeGoodTxs([]byte(Config.YktApiSecret), data)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	var (
		val      []string
		bonusMap = make(map[string]*GoodTxBonus)
	)
	for _, tx := range txs {
		income := decimal.New(int64(tx.Income), -4)
		key := fmt.Sprintf("%d-%d", tx.GoodId, tx.Uid)
		if bonus, found := bonusMap[key]; found {
			bonus.Amount += tx.Amount
		} else {
			bonusMap[key] = &GoodTxBonus{
				GoodId: tx.GoodId,
				UserId: tx.Uid,
				Amount: tx.Amount,
			}
		}
		createdAt, err := time.Parse(time.RFC3339, tx.CreatedAt)
		if CheckErr(err, c) {
			return
		}
		val = append(val, fmt.Sprintf("(%d, %d, %d, %d, %s, '%s')", tx.OrderId, tx.Uid, tx.GoodId, tx.Amount, income.String(), createdAt.UTC().Format("2006-01-02 15:04:05")))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO tmm.good_txs (oid, uid, good_id, amount, income, created_at) VALUES %s`, strings.Join(val, ","))
		if CheckErr(err, c) {
			return
		}
	}
	tspoints, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	for _, bonus := range bonusMap {
		points := decimal.New(int64(bonus.Amount)*Config.GoodCommissionPoints, 0)
		ts := points.Div(tspoints)
		_, _, err := db.Query(`UPDATE tmm.devices AS d, tmm.good_invests AS gi SET d.points=d.points+%s, d.total_ts=d.total_ts+%d, gi.bonus=gi.bonus+%s WHERE d.id=gi.device_id AND gi.user_id=%d AND gi.good_id=%d`, points.String(), ts.IntPart(), points.String(), bonus.UserId, bonus.GoodId)
		if err != nil {
			log.Error(err.Error())
		}
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
