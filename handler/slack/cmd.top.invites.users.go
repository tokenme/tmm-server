package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"time"
)

func TopInvitesUsersHandler(c *gin.Context, num int64) {
	if num < 0 {
		num = 10
	}
	if num > 100 {
		num = 100
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT
    u.id AS id,
    u.country_code AS country_code,
    u.mobile AS mobile,
    u.nickname AS nick,
    u.avatar AS avatar,
    wx.nick AS wx_nick,
    wx.avatar AS wx_avatar,
    COUNT( 0 ) AS invites
FROM tmm.invite_codes AS ic
INNER JOIN ucoin.users AS u ON (u.id = ic.parent_id)
LEFT JOIN tmm.wx AS wx ON ( wx.user_id = u.id )
GROUP BY u.id
ORDER BY invites DESC LIMIT %d`, num)
	if CheckErr(err, c) {
		return
	}
	var data [][]string
	for _, row := range rows {
		u := common.User{
			Id:          row.Uint64(0),
			CountryCode: row.Uint(1),
			Mobile:      row.Str(2),
			Nick:        row.Str(3),
			Avatar:      row.Str(4),
		}
		wxNick := row.Str(5)
		if wxNick != "" {
			wechat := &common.Wechat{
				Nick: wxNick,
			}
			u.Wechat = wechat
		}
		u.ShowName = u.GetShowName()
		invites, _ := decimal.NewFromString(row.Str(6))
		data = append(data, []string{strconv.FormatUint(u.Id, 10), u.ShowName, invites.StringFixed(9)})
	}

	var b bytes.Buffer
	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"Id", "Nick", "Invites"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
	msg := WebhookMessage{
		Text: fmt.Sprintf("Top %d Invites Users", num),
		Attachments: []slack.Attachment{
			{
				Text: b.String(),
				Ts:   json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
