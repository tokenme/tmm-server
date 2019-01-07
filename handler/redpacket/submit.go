package redpacket

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/utils"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
)

type SubmitRequest struct {
	UnionId     string `json:"union_id"`
	Nick        string `json:"nick"`
	Avatar      string `json:"avatar"`
	RedpacketId string `json:"redpacket_id" binding:"required"`
}

func SubmitHandler(c *gin.Context) {
	var (
		user common.User
		req  SubmitRequest
	)
	if CheckErr(c.Bind(&req), c) {
		return
	}
	packetId, err := utils.DecryptUint64(req.RedpacketId, []byte(Config.LinkSalt))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(packetId == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	userContext, exists := c.Get("USER")
	if exists {
		user = userContext.(common.User)
	} else if Check(req.UnionId == "", "bad request", c) {
		return
	}

	var (
		ret  mysql.Result
		resp common.RedpacketRecipient
	)
	db := Service.Db
	if user.Id > 0 {
		unionId := "NULL"
		nick := utils.HideMobile(user.Mobile)
		{
			rows, _, _ := db.Query(`SELECT union_id, nick, avatar FROM tmm.wx WHERE user_id=%d LIMIT 1`, user.Id)
			if len(rows) > 0 {
				user.Wechat = &common.Wechat{
					UnionId: rows[0].Str(0),
					Nick:    rows[0].Str(1),
					Avatar:  rows[0].Str(2),
				}
				unionId = fmt.Sprintf("'%s'", db.Escape(user.Wechat.UnionId))
				nick = user.GetShowName()
			}
		}
		resp.Nick = user.GetShowName()
		resp.Avatar = user.GetAvatar(Config.CDNUrl)
		_, ret, err = db.Query(`UPDATE tmm.redpacket_recipients SET user_id=%d, union_id=%s, nick='%s', avatar='%s', submited_at=NOW() WHERE redpacket_id=%d AND (user_id IS NULL OR union_id IS NULL) LIMIT 1`, user.Id, unionId, db.Escape(nick), db.Escape(resp.Avatar), packetId)
	} else {
		userId := "NULL"
		{
			rows, _, _ := db.Query(`SELECT u.id FROM tmm.wx AS wx INNER JOIN ucoin.users AS u ON (u.id=wx.id) WHERE wx.union_id='%s' LIMIT 1`, req.UnionId)
			if len(rows) > 0 {
				userId = fmt.Sprintf("%d", rows[0].Uint64(0))
			}
		}
		resp.Nick = req.Nick
		resp.Avatar = req.Avatar
		_, ret, err = db.Query(`UPDATE tmm.redpacket_recipients SET user_id=%s, union_id='%s', nick='%s', avatar='%s', submited_at=NOW() WHERE redpacket_id=%d AND (user_id IS NULL OR union_id IS NULL) LIMIT 1`, userId, db.Escape(req.UnionId), db.Escape(resp.Nick), db.Escape(resp.Avatar), packetId)
	}

	if err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY {
		c.JSON(http.StatusOK, APIResponse{Msg: "submitted"})
		return
	}
	if CheckErr(err, c) {
		return
	}
	if ret.AffectedRows() == 0 {
		c.JSON(http.StatusOK, APIResponse{Msg: "unlucky"})
		return
	}

	if user.Id > 0 {
		rows, _, err := db.Query(`SELECT tmm FROM tmm.redpacket_recipients WHERE redpacket_id=%d AND user_id=%d LIMIT 1`, packetId, user.Id)
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			resp.Tmm, _ = decimal.NewFromString(rows[0].Str(0))
		}
	} else {
		rows, _, err := db.Query(`SELECT tmm FROM tmm.redpacket_recipients WHERE red_packet_id=%d AND union_id='%s' LIMIT 1`, packetId, db.Escape(req.UnionId))
		if CheckErr(err, c) {
			return
		}
		if len(rows) > 0 {
			resp.Tmm, _ = decimal.NewFromString(rows[0].Str(0))
		}
	}
	tmmPrice := common.GetTMMPrice(Service, Config, common.MarketPrice)
	rate := forex.Rate(Service, "USD", "CNY")
	tmmPrice = tmmPrice.Mul(rate)
	resp.Cash = tmmPrice.Mul(resp.Tmm)
	c.JSON(http.StatusOK, resp)
}
