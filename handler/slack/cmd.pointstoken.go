package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	//"github.com/shopspring/decimal"
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/tokenme/tmm/common"
	"net/http"
	"strconv"
	"time"
)

func PointsTokenHandler(c *gin.Context) {
	exchangeRate, pointsPerTs, err := common.GetExchangeRate(Config, Service)
	if CheckErr(err, c) {
		return
	}
	tmm, err := common.GetTMMPerTs(Config, Service)
	if CheckErr(err, c) {
		return
	}
	msg := WebhookMessage{
		Text: "Points <-> Token ExchangeRate",
		Attachments: []slack.Attachment{
			{
				Fields: []slack.AttachmentField{
					{
						Title: "UC/Points",
						Value: exchangeRate.Rate.String(),
						Short: true,
					},
					{
						Title: "Min Points",
						Value: exchangeRate.MinPoints.String(),
						Short: true,
					},
					{
						Title: "Points/Ts",
						Value: pointsPerTs.String(),
						Short: true,
					},
					{
						Title: "UC/Ts",
						Value: tmm.String(),
						Short: true,
					},
				},
				Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
