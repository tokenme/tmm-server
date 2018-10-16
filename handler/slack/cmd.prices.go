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

func PricesHandler(c *gin.Context) {
	pointPrice := common.GetPointPrice(Service, Config)
	marketPrice := common.GetTMMPrice(Service, Config, common.MarketPrice)
	recyclePrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	ethPrice := common.GetETHPrice(Service, Config)
	tokenRate := common.GetTMMRate(Service, Config)
	msg := WebhookMessage{
		Text: "Prices",
		Attachments: []slack.Attachment{
			{
				Fields: []slack.AttachmentField{
					{
						Title: "ETH/USD",
						Value: ethPrice.String(),
						Short: true,
					},
					{
						Title: "Points/USD",
						Value: pointPrice.String(),
						Short: true,
					},
					{
						Title: "UC Market/USD",
						Value: marketPrice.String(),
						Short: true,
					},
					{
						Title: "UC Recycle/USD",
						Value: recyclePrice.String(),
						Short: true,
					},
					{
						Title: "UC/ETH",
						Value: tokenRate.String(),
						Short: true,
					},
				},
				Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
