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

func DailyWithdrawHandler(c *gin.Context, days int64) {
	if days < 0 {
		days = 30
	}
	if days > 180 {
		days = 180
	}
	startDate := time.Now().Add(-1 * time.Duration(days)).Format("2006-06-02")
	db := Service.Db
	rows, _, err := db.Query(`SELECT record_on, SUM(cny) AS cny
    FROM (
    SELECT
            DATE(tx.inserted_at) AS record_on, SUM(tx.cny) AS cny
    FROM tmm.withdraw_txs AS tx
    WHERE tx.inserted_at>='%s'
    GROUP BY record_on
    UNION ALL
    SELECT
            DATE(pw.inserted_at) AS record_on, SUM(pw.cny) AS cny
    FROM tmm.point_withdraws AS pw
    WHERE pw.inserted_at >= '%s'
    GROUP BY record_on) AS t GROUP BY record_on ORDER BY record_on ASC`, startDate, startDate)
	if CheckErr(err, c) {
		return
	}
	var bars []BarChartValue
	for _, row := range rows {
		bars = append(bars, BarChartValue{
			Value: row.ForceFloat(1),
			Label: row.ForceLocaltime(0).Format("2006-01-02"),
		})
	}
	js, err := json.Marshal(bars)
	if CheckErr(err, c) {
		return
	}
	data := base64.URLEncoding.EncodeToString(js)
	var b bytes.Buffer
	table := tablewriter.NewWriter(&b)
	table.SetHeader([]string{"Date", "Cash"})
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
		Text: "Daily Withdraws",
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
