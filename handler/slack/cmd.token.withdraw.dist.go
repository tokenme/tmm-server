package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"github.com/olekukonko/tablewriter"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"time"
)

func TokenWithdrawDistHandler(c *gin.Context) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT
    COUNT(*) AS users,
    l
FROM (
    SELECT
        tx.user_id,
        FLOOR(LOG10(SUM(tx.cny))) AS l
    FROM tmm.withdraw_txs AS tx
    GROUP BY tx.user_id
) AS tmp
GROUP BY l ORDER BY l`)
	if CheckErr(err, c) {
		return
	}
	var bars []BarChartValue
	for _, row := range rows {
		bars = append(bars, BarChartValue{
			Value: row.ForceFloat(0),
			Label: fmt.Sprintf("10^%d", row.Int(1)),
		})
	}
	js, err := json.Marshal(bars)
	if CheckErr(err, c) {
		return
	}
	data := base64.URLEncoding.EncodeToString(js)
	var b bytes.Buffer
	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"CNY", "Users"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, bar := range bars {
		table.Append([]string{
			bar.Label,
			fmt.Sprintf("%v", bar.Value),
		})
	}
	table.Render()
	msg := WebhookMessage{
		Text: "User Withdraw Distribution",
		Attachments: []slack.Attachment{
			{
				Text:     b.String(),
				ImageURL: fmt.Sprintf("https://tmm.tokenmama.io/slack/chart/bar/%s", data),
				Ts:       json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
