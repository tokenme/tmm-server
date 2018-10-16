package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	//"github.com/shopspring/decimal"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/tokenme/tmm/common"
	"net/http"
	"strconv"
	"time"
)

func TsTokenHandler(c *gin.Context) {
	tmm, err := common.GetTMMPerTs(Config, Service)
	if CheckErr(err, c) {
		return
	}
	msg := WebhookMessage{
		Text: fmt.Sprintf("%s UC/ts", tmm.String()),
		Attachments: []slack.Attachment{
			{
				Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
