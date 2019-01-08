package redpacket

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	//"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/utils"
	"github.com/ziutek/mymysql/mysql"
	"net/http"
)

type SubmitRequest struct {
	UnionId     string `json:"union_id" form:"union_id"`
	Nick        string `json:"nick" form:"nick"`
	Avatar      string `json:"avatar" form:"avatar"`
	RecipientId string `json:"recipient_id" form:"recipient_id"`
	RedpacketId string `json:"redpacket_id" form:"redpacket_id" binding:"required"`
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
		log.Error(err.Error())
		return
	}

	userId, _ := utils.DecryptUint64(req.RecipientId, []byte(Config.LinkSalt))

	if CheckWithCode(packetId == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	if Check(req.UnionId == "" && userId == 0, "bad request", c) {
		return
	}

	var (
		ret  mysql.Result
		resp common.RedpacketRecipient
	)
	log.Warn("UID:%d, packetId: %s", userId, req.UnionId)
	db := Service.Db
	if userId > 0 {
		unionId := "NULL"
		var nick string
		{
			rows, _, err := db.Query(`SELECT u.id, u.country_code, u.mobile, wx.union_id, wx.nick, wx.avatar FROM ucoin.users AS u LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id) WHERE u.id=%d LIMIT 1`, userId)
			if err != nil {
				log.Error(err.Error())
			}
			if Check(len(rows) == 0, "bad request", c) {
				return
			}

			user = common.User{
				Id:          rows[0].Uint64(0),
				CountryCode: rows[0].Uint(1),
				Mobile:      rows[0].Str(2),
			}
			wxUnionId := rows[0].Str(3)
			if wxUnionId != "" {
				user.Wechat = &common.Wechat{
					Nick:   rows[0].Str(4),
					Avatar: rows[0].Str(4),
				}
				unionId = fmt.Sprintf("'%s'", db.Escape(wxUnionId))
				nick = user.GetShowName()
			} else {
				nick = utils.HideMobile(user.Mobile)
			}
		}
		resp.Nick = user.GetShowName()
		resp.Avatar = user.GetAvatar(Config.CDNUrl)
		_, ret, err = db.Query(`UPDATE tmm.redpacket_recipients SET user_id=%d, union_id=%s, nick='%s', avatar='%s', submited_at=NOW() WHERE redpacket_id=%d AND (user_id IS NULL AND union_id IS NULL) LIMIT 1`, user.Id, unionId, db.Escape(nick), db.Escape(resp.Avatar), packetId)
	} else {
		userId := "NULL"
		{
			rows, _, _ := db.Query(`SELECT u.id FROM tmm.wx AS wx INNER JOIN ucoin.users AS u ON (u.id=wx.user_id) WHERE wx.union_id='%s' LIMIT 1`, req.UnionId)
			if len(rows) > 0 {
				userId = fmt.Sprintf("%d", rows[0].Uint64(0))
			}
		}
		resp.Nick = req.Nick
		resp.Avatar = req.Avatar
		_, ret, err = db.Query(`UPDATE tmm.redpacket_recipients SET user_id=%s, union_id='%s', nick='%s', avatar='%s', submited_at=NOW() WHERE redpacket_id=%d AND (user_id IS NULL AND union_id IS NULL) LIMIT 1`, userId, db.Escape(req.UnionId), db.Escape(resp.Nick), db.Escape(resp.Avatar), packetId)
	}

	if err != nil && err.(*mysql.Error).Code == mysql.ER_DUP_ENTRY {
		//c.JSON(http.StatusOK, APIResponse{Msg: "submitted"})
		//return
		log.Warn("submitted")
	} else if CheckErr(err, c) {
		return
	} else if ret.AffectedRows() == 0 {
		//c.JSON(http.StatusOK, APIResponse{Msg: "unlucky"})
		//return
		log.Warn("unlucky")
	}
	/*
		if userId > 0 {
			rows, _, err := db.Query(`SELECT tmm FROM tmm.redpacket_recipients WHERE redpacket_id=%d AND user_id=%d LIMIT 1`, packetId, userId)
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
	*/
	recipients, err := getRecipients(packetId, userId, req.UnionId, Service)
	tmmPrice := common.GetTMMPrice(Service, Config, common.MarketPrice)
	rate := forex.Rate(Service, "USD", "CNY")
	tmmPrice = tmmPrice.Mul(rate)
	for _, receipt := range recipients {
		receipt.Cash = tmmPrice.Mul(receipt.Tmm).RoundBank(2)
	}
	c.JSON(http.StatusOK, recipients)
}
