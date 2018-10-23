package feedback

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

type ReplyRequest struct {
	Message string `json:"message" form:"message" binding:"required"`
	Ts      string `json:"ts" form:"ts" binding:"required"`
}

func ReplyHandler(c *gin.Context) {
	var req ReplyRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	params := slack.PostMessageParameters{
		Username:        user.ShowName,
		IconURL:         user.Avatar,
		ThreadTimestamp: req.Ts,
		Parse:           "full",
		UnfurlMedia:     true,
		Markdown:        true,
	}
	_, _, err := Service.Slack.PostMessage(Config.Slack.FeedbackChannel, req.Message, params)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
