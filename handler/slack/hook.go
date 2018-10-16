package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	//. "github.com/tokenme/tmm/handler"
	//"github.com/shopspring/decimal"
	//"github.com/tokenme/tmm/common"
	//"net/http"
)

type WebhookMessage struct {
	Text        string             `json:"text,omitempty"`
	Attachments []slack.Attachment `json:"attachments,omitempty"`
}

func HookHandler(c *gin.Context) {
	cmd := c.PostForm("command")
	switch cmd {
	case "/tspoints":
		TsPointsHandler(c)
	case "/tstoken":
		TsTokenHandler(c)
	case "/pointstoken":
		PointsTokenHandler(c)
	case "/prices":
		PricesHandler(c)
	}
}
