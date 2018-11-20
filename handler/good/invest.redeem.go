package good

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/wechatpay"
	"github.com/tokenme/tmm/tools/ykt"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strconv"
	"strings"
)

type InvestRedeemRequest struct {
	GoodIds string `json:"ids" form:"ids"`
}

func InvestRedeemHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)

	var req InvestRedeemRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	idArr := strings.Split(req.GoodIds, ",")
	var goodIds []string
	for _, idStr := range idArr {
		idVal, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || idVal == 0 {
			continue
		}
		goodIds = append(goodIds, idStr)
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT wx.union_id, oi.open_id FROM tmm.wx LEFT JOIN tmm.wx_openids AS oi ON (oi.union_id=wx.union_id AND oi.app_id='%s') WHERE wx.user_id=%d LIMIT 1`, db.Escape(Config.Wechat.AppId), user.Id)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, WECHAT_UNAUTHORIZED_ERROR, "Wechat Unauthorized", c) {
		return
	}
	wxUnionId := rows[0].Str(0)
	wxOpenId := rows[0].Str(1)
	if wxOpenId == "" {
		openIdReq := ykt.OpenIdRequest{
			UnionId: wxUnionId,
		}
		openIdRes, err := openIdReq.Run()
		if CheckWithCode(err != nil, WECHAT_OPENID_ERROR, "need openid", c) {
			return
		}
		wxOpenId = openIdRes.Data.OpenId
		_, _, err = db.Query(`INSERT INTO tmm.wx_openids (app_id, open_id, union_id) VALUES ('%s', '%s', '%s')`, db.Escape(Config.Wechat.AppId), db.Escape(wxOpenId), db.Escape(wxUnionId))
		if CheckErr(err, c) {
			return
		}
	}
	query := `SELECT good_id, SUM(income) AS income
FROM (
SELECT gi.good_id AS good_id, IFNULL(tx.income, 0) * gi.points/SUM(gi2.points) AS income
FROM tmm.good_invests AS gi
INNER JOIN tmm.good_txs AS tx ON (tx.good_id=gi.good_id AND tx.created_at>=gi.inserted_at)
INNER JOIN tmm.good_invests AS gi2 ON (gi2.good_id=gi.good_id AND gi2.inserted_at<=tx.created_at)
WHERE gi.user_id=%d%s
GROUP BY tx.oid) AS tmp
GROUP BY good_id
HAVING income > 0`
	var where string
	if len(goodIds) > 0 {
		where = fmt.Sprintf(" AND gi.good_id IN (%s)", strings.Join(goodIds, ","))
	}
	rows, _, err = db.Query(query, user.Id, where, user.Id, where)
	if CheckErr(err, c) {
		return
	}
	if len(rows) == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
		return
	}
	payClient := wechatpay.NewClient(Config.Wechat.AppId, Config.Wechat.MchId, Config.Wechat.Key, Config.Wechat.CertCrt, Config.Wechat.CertKey)
	var (
		totalIncome   decimal.Decimal
		updateGoodIds []string
	)
	for _, row := range rows {
		goodId := row.Uint64(0)
		income, _ := decimal.NewFromString(row.Str(1))
		totalIncome = totalIncome.Add(income)
		updateGoodIds = append(updateGoodIds, strconv.FormatUint(goodId, 10))
		log.Warn("Good:%d, income:%s, total:%s", goodId, income, totalIncome)
	}
	_, _, err = db.Query(`UPDATE tmm.good_invests SET redeem_status=1, redeem_at=NOW() WHERE user_id=%d AND good_id IN (%s)`, user.Id, strings.Join(updateGoodIds, ","))
	if CheckErr(err, c) {
		return
	}
	tradeNumToken, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	payParams := &wechatpay.Request{
		TradeNum:    utils.Md5(tradeNumToken.String()),
		Amount:      totalIncome.Mul(decimal.New(100, 0)).IntPart(),
		CallbackURL: fmt.Sprintf("%s/wechat/pay/callback", Config.BaseUrl),
		OpenId:      wxOpenId,
		Ip:          ClientIP(c),
		Desc:        "UCoin 商城投资提现",
	}
	payParams.Nonce = utils.Md5(payParams.TradeNum)
	payRes, err := payClient.Pay(payParams)
	if CheckErr(err, c) {
		log.Error(err.Error())
		_, _, err = db.Query(`UPDATE tmm.good_invests SET redeem_status=0, redeem_at=NULL WHERE user_id=%d AND good_id IN (%s) AND redeem_status=1`, user.Id, strings.Join(updateGoodIds, ","))
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	if CheckWithCode(payRes.ErrCode != "", WECHAT_PAYMENT_ERROR, payRes.ErrCodeDesc, c) {
		_, _, err = db.Query(`UPDATE tmm.good_invests SET redeem_status=0, redeem_at=NULL WHERE user_id=%d AND good_id IN (%s) AND redeem_status=1`, user.Id, strings.Join(updateGoodIds, ","))
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
