package redpacket

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/tools/wechatmp"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"sort"
	"strings"
)

const (
	WX_AUTH_GATEWAY = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
)

type ShowData struct {
	AppId            string
	ShareTitle       string
	ShareLink        string
	ShareImage       string
	JSConfig         wechatmp.JSConfig
	PacketId         string                       `json:"packet_id"`
	Title            string                       `json:"title"`
	Message          string                       `json:"message"`
	Tmm              decimal.Decimal              `json:"tmm"`
	Cash             decimal.Decimal              `json:"cash"`
	Nick             string                       `json:"nick"`
	Avatar           string                       `json:"avatar"`
	RecipientUserId  string                       `json:"recipient_user_id"`
	RecipientUnionId string                       `json:"recipient_union_id"`
	RecipientNick    string                       `json:"recipient_nick"`
	RecipientAvatar  string                       `json:"recipient_avatar"`
	Recipients       []*common.RedpacketRecipient `json:"recipients"`
	NotSubmitted     bool                         `json:"not_submitted"`
}

func ShowHandler(c *gin.Context) {
	var (
		user      common.User
		recipient common.User
	)
	packetId, userId, err := common.DecryptRedpacketLink(c.Param("encryptedId"), c.Param("encryptedUserId"), Config)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(packetId == 0 || userId == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	db := Service.Db
	isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")
	var (
		recipientUserId  uint64
		recipientUnionId string
	)
	mpClient := wechatmp.NewClient(Config.Wechat.AppId, Config.Wechat.AppSecret, Service, c)
	if isWx {
		code := c.DefaultQuery("code", "null")
		if len(code) > 0 && code != "null" {
			oauthAccessToken, err := mpClient.GetOAuthAccessToken(code)
			if err != nil {
				log.Error(err.Error())
			} else if oauthAccessToken.Openid != "" {
				userInfo, err := mpClient.GetUserInfo(oauthAccessToken.AccessToken, oauthAccessToken.Openid)
				if err != nil {
					log.Error(err.Error())
				} else {
					recipient.Wechat = &common.Wechat{
						UnionId: userInfo.UnionId,
						Nick:    userInfo.Nick,
						Avatar:  userInfo.Avatar,
					}
					recipientUnionId = userInfo.UnionId
					rows, _, err := db.Query(`SELECT user_id FROM tmm.wx AS wx WHERE union_id='%s' LIMIT 1`, db.Escape(recipientUnionId))
					if err != nil {
						log.Error(err.Error())
					}
					if len(rows) > 0 {
						recipientUserId = rows[0].Uint64(0)
					}
				}
			}
		} else if code == "null" {
			redpacketUri := c.Request.URL.String()
			redirectUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, redpacketUri)
			mpClient.AuthRedirect(redirectUrl, "snsapi_base", "redpacket")
			return
		}
	} else {
		recipientUserId = userId
	}
	{
		rows, _, err := db.Query(`SELECT u.id, u.country_code, u.mobile, wx.nick, wx.avatar FROM ucoin.users AS u LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id) WHERE u.id=%d LIMIT 1`, userId)
		if err != nil {
			log.Error(err.Error())
		}
		if len(rows) > 0 {
			user = common.User{
				Id:          rows[0].Uint64(0),
				CountryCode: rows[0].Uint(1),
				Mobile:      utils.HideMobile(rows[0].Str(2)),
			}
			wxNick := rows[0].Str(3)
			if wxNick != "" {
				user.Wechat = &common.Wechat{
					Nick:   wxNick,
					Avatar: rows[0].Str(4),
				}
			}
		}
	}
	if Check(user.Id == 0 || isWx && recipient.Wechat == nil, "unauthorized request", c) {
		return
	}
	rows, _, err := db.Query(`SELECT
        rp.id, rp.message, rp.tmm, rp.recipients, rp.creator, IFNULL(u.country_code, 0), IFNULL(u.mobile, ""), IFNULL(u.avatar, ""), wx.nick, wx.avatar
        FROM tmm.redpackets AS rp
        LEFT JOIN ucoin.users AS u ON (u.id=rp.creator)
        LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id)
        WHERE rp.id=%d LIMIT 1`, packetId)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	row := rows[0]
	tmm, _ := decimal.NewFromString(row.Str(2))
	creator := common.User{
		Id:          row.Uint64(4),
		CountryCode: row.Uint(5),
		Mobile:      utils.HideMobile(row.Str(6)),
		Avatar:      row.Str(7),
	}
	wxNick := row.Str(8)
	if wxNick != "" {
		creator.Wechat = &common.Wechat{
			Nick:   wxNick,
			Avatar: row.Str(9),
		}
	}
	currentUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, c.Request.URL.String())
	jsConfig, err := mpClient.GetJSConfig(currentUrl)
	if err != nil {
		log.Error(err.Error())
	}
	data := ShowData{
		AppId:      Config.Wechat.AppId,
		JSConfig:   jsConfig,
		Title:      "UCoin红包，快来抢",
		PacketId:   c.Param("encryptedId"),
		ShareImage: "https://ucoin.tianxi100.com/redpacket.jpeg",
		ShareLink:  currentUrl,
		Message:    row.Str(1),
		Tmm:        tmm,
	}
	if creator.Id > 0 {
		data.Nick = creator.GetShowName()
		data.Avatar = creator.GetAvatar(Config.CDNUrl)
	} else {
		data.Nick = user.GetShowName()
		data.Avatar = user.GetAvatar(Config.CDNUrl)
	}

	if recipient.Wechat != nil {
		data.Title = fmt.Sprintf("快来抢%s的红包", data.Nick)
		data.RecipientUnionId = recipient.Wechat.UnionId
		data.RecipientNick = recipient.GetShowName()
		data.RecipientAvatar = recipient.GetAvatar(Config.CDNUrl)
		data.ShareTitle = fmt.Sprintf("快来抢%s的红包", data.RecipientNick)
	} else if user.Id > 0 {
		data.ShareTitle = fmt.Sprintf("快来抢%s的红包", data.Nick)
		data.RecipientUserId = c.Param("encryptedUserId")
		data.RecipientNick = user.GetShowName()
		data.RecipientAvatar = user.GetAvatar(Config.CDNUrl)
	}

	if recipientUserId > 0 {
		redpacket := common.Redpacket{
			Id: packetId,
		}
		data.ShareLink, _ = redpacket.GetLink(Config, recipientUserId)
	}
	if data.Message == "" {
		data.Message = "大吉大利，恭喜发财"
	}
	tmmPrice := common.GetTMMPrice(Service, Config, common.MarketPrice)
	rate := forex.Rate(Service, "USD", "CNY")
	tmmPrice = tmmPrice.Mul(rate)
	data.Cash = tmmPrice.Mul(data.Tmm).RoundBank(2)
	recipients, err := getRecipients(packetId, recipientUserId, recipientUnionId, Service)
	var submitted bool
	for _, receipt := range recipients {
		receipt.Cash = tmmPrice.Mul(receipt.Tmm).RoundBank(2)
		if recipientUserId > 0 && recipientUserId == receipt.UserId || recipientUnionId != "" && receipt.UnionId == recipientUnionId {
			submitted = true
		}
	}
	if submitted {
		data.Recipients = recipients
	}
	data.NotSubmitted = !submitted
	c.HTML(http.StatusOK, "redpacket.tmpl", data)
}

func getRecipients(id uint64, userId uint64, unionId string, service *common.Service) (recipients []*common.RedpacketRecipient, err error) {
	db := service.Db
	rows, _, err := db.Query(`SELECT
        rr.user_id, IFNULL(wx.union_id, rr.union_id), IFNULL(wx.nick, rr.nick), IFNULL(wx.avatar, rr.avatar), rr.tmm
    FROM tmm.redpacket_recipients AS rr
    LEFT JOIN ucoin.users AS u ON (u.id=rr.user_id)
    LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id)
    WHERE rr.redpacket_id=%d AND (rr.user_id IS NOT NULL OR rr.union_id IS NOT NULL) ORDER BY rr.tmm DESC`, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	for _, row := range rows {
		tmm, _ := decimal.NewFromString(row.Str(4))
		uid := row.Uint64(0)
		unid := row.Str(1)
		rr := &common.RedpacketRecipient{
			UserId:  uid,
			UnionId: unid,
			Nick:    row.Str(2),
			Avatar:  row.Str(3),
			Tmm:     tmm,
			IsSelf:  userId == uid || unionId == unid,
		}
		recipients = append(recipients, rr)
	}
	sorter := common.NewRecipientsSorter(recipients)
	sort.Sort(sort.Reverse(sorter))
	return sorter, nil
}
