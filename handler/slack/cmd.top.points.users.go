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

func TopPointsUsersHandler(c *gin.Context, num int64) {
	if num < 0 {
		num = 10
	}
	if num > 100 {
		num = 100
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT
        id,
        country_code,
        mobile,
        nick,
        wx_nick,
        total_points
    FROM tmm.top_points_users
    ORDER BY total_points DESC LIMIT %d`, num)
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
		}
		wxNick := row.Str(4)
		if wxNick != "" {
			wechat := &common.Wechat{
				Nick: wxNick,
			}
			u.Wechat = wechat
		}
		u.ShowName = u.GetShowName()
		points, _ := decimal.NewFromString(row.Str(5))
		data = append(data, []string{strconv.FormatUint(u.Id, 10), u.ShowName, points.StringFixed(9)})
	}

	var b bytes.Buffer
	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"Id", "Nick", "Points"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(data)
	table.Render()
	msg := WebhookMessage{
		Text: fmt.Sprintf("Top %d Points Users", num),
		Attachments: []slack.Attachment{
			{
				Text: b.String(),
				Ts:   json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
