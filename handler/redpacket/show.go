package redpacket

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/tools/wechatmp"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strings"
)

const (
	WX_AUTH_GATEWAY = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
)

type ShowData struct {
	PacketId         string          `json:"packet_id"`
	Message          string          `json:"message"`
	Tmm              decimal.Decimal `json:"tmm"`
	Cash             decimal.Decimal `json:"cash"`
	Nick             string          `json:"nick"`
	Avatar           string          `json:"avatar"`
	RecipientUnionId string          `json:"recipient_union_id"`
	RecipientNick    string          `json:"recipient_nick"`
	RecipientAvatar  string          `json:"recipient_avatar"`
}

func ShowHandler(c *gin.Context) {
	var (
		user      common.User
		recipient common.User
	)
	packetId, userId, err := common.DecryptRedpacketLink(c.Param("encryptedId"), c.Param("encryptedDeviceId"), Config)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(packetId == 0 || userId == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	db := Service.Db
	isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")
	if isWx {
		mpClient := wechatmp.NewClient(Config.Wechat.AppId, Config.Wechat.AppSecret, Service, c)
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
				}
			}
		} else if code == "null" {
			redpacketUri := c.Request.URL.String()
			redirectUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, redpacketUri)
			mpClient.AuthRedirect(redirectUrl, "snsapi_base", "redpacket")
			return
		}
	} else {
		rows, _, _ := db.Query(`SELECT u.country_code, u.mobile, wx.nick, wx.avatar FROM ucoin.users AS u LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id) WHERE u.id=%d LIMIT 1`, userId)
		if len(rows) > 0 {
			user = common.User{
				CountryCode: rows[0].Uint(0),
				Mobile:      utils.HideMobile(rows[0].Str(1)),
			}
			wxNick := rows[0].Str(2)
			if wxNick != "" {
				user.Wechat = &common.Wechat{
					Nick:   wxNick,
					Avatar: rows[0].Str(3),
				}
			}
		}
	}
	if Check(user.Id == 0 && recipient.Wechat == nil, "unauthorized request", c) {
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
	data := ShowData{
		PacketId: c.Param("encryptedId"),
		Message:  row.Str(1),
		Tmm:      tmm,
	}
	if recipient.Wechat != nil {
		data.Nick = user.GetShowName()
		data.Avatar = user.GetAvatar(Config.CDNUrl)
		data.RecipientUnionId = recipient.Wechat.UnionId
		data.RecipientNick = recipient.GetShowName()
		data.RecipientAvatar = recipient.GetAvatar(Config.CDNUrl)
	} else if user.Id > 0 {
		data.RecipientNick = user.GetShowName()
		data.RecipientAvatar = user.GetAvatar(Config.CDNUrl)
	}
	tmmPrice := common.GetTMMPrice(Service, Config, common.MarketPrice)
	rate := forex.Rate(Service, "USD", "CNY")
	tmmPrice = tmmPrice.Mul(rate)
	data.Cash = tmmPrice.Mul(data.Tmm)
	c.HTML(http.StatusOK, "redpacket.tmpl", data)
}

func getRecipients(id uint64, service *common.Service) (recipients []*common.RedpacketRecipient, err error) {
	db := service.Db
	rows, _, err := db.Query(`SELECT
        rr.user_id, IFNULL(wx.union_id, rr.union_id), IFNULL(wx.nick, rr.nick), IFNULL(wx.avatar, rr.avatar), rr.tmm
    FROM tmm.redpacket_recipients AS rr
    LEFT JOIN ucoin.users AS u ON (u.id=rr.id)
    LEFT JOIN tmm.wx AS wx ON (wx.user_id=u.id)
    WHERE rr.id=%d ORDER BY rr.tmm DESC`)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	for _, row := range rows {
		tmm, _ := decimal.NewFromString(row.Str(4))
		rr := &common.RedpacketRecipient{
			UserId:  row.Uint64(0),
			UnionId: row.Str(1),
			Nick:    row.Str(2),
			Avatar:  row.Str(3),
			Tmm:     tmm,
		}
		recipients = append(recipients, rr)
	}
	return recipients, nil
}
