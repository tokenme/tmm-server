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

func TsPointsHandler(c *gin.Context) {
	points, err := common.GetPointsPerTs(Service)
	if CheckErr(err, c) {
		return
	}
	msg := WebhookMessage{
		Text: fmt.Sprintf("%s points/ts", points.String()),
		Attachments: []slack.Attachment{
			{
				Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
