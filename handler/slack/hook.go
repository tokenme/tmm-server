package slack

import (
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"strconv"
)

type WebhookMessage struct {
	Text        string             `json:"text,omitempty"`
	Attachments []slack.Attachment `json:"attachments,omitempty"`
}

func HookHandler(c *gin.Context) {
	cmd := c.PostForm("command")
	txt := c.PostForm("text")
	switch cmd {
	case "/tspoints":
		TsPointsHandler(c)
	case "/tstoken":
		TsTokenHandler(c)
	case "/pointstoken":
		PointsTokenHandler(c)
	case "/prices":
		PricesHandler(c)
	case "/points.dist":
		PointsDistHandler(c)
	case "/withdraw.dist":
		WithdrawDistHandler(c)
	case "/daily.withdraw":
		num, _ := strconv.ParseInt(txt, 10, 64)
		DailyWithdrawHandler(c, num)
	case "/token.withdraw.dist":
		TokenWithdrawDistHandler(c)
	case "/point.withdraw.dist":
		PointWithdrawDistHandler(c)
	case "/invites.dist":
		InvitesDistHandler(c)
	case "/top.points.users":
		num, _ := strconv.ParseInt(txt, 10, 64)
		TopPointsUsersHandler(c, num)
	case "/top.withdraw.users":
		num, _ := strconv.ParseInt(txt, 10, 64)
		TopWithdrawUsersHandler(c, num)
	case "/top.token.withdraw.users":
		num, _ := strconv.ParseInt(txt, 10, 64)
		TopTokenWithdrawUsersHandler(c, num)
	case "/top.point.withdraw.users":
		num, _ := strconv.ParseInt(txt, 10, 64)
		TopPointWithdrawUsersHandler(c, num)
	case "/top.invites.users":
		num, _ := strconv.ParseInt(txt, 10, 64)
		TopInvitesUsersHandler(c, num)
	case "/block.reason":
		BlockReasonHandler(c, txt)
	}
}
