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

func BlockReasonHandler(c *gin.Context, mobile string) {
	db := Service.Db
	rows, _, err := db.Query(`SELECT id FROM ucoin.users WHERE mobile='%s' LIMIT 1`, db.Escape(mobile))
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(len(rows) == 0, NOTFOUND_ERROR, "not found", c) {
		return
	}
	user := common.User{
		Id: rows[0].Uint64(0),
	}
	err = user.BlockReason(Service)
	var reason string
	if err == nil {
		reason = "User not blocked"
	} else {
		reason = err.Error()
	}
	msg := WebhookMessage{
		Text: reason,
		Attachments: []slack.Attachment{
			{
				Ts: json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
			},
		},
	}
	c.JSON(http.StatusOK, msg)
}
