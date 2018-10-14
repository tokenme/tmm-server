package feedback

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AddRequest struct {
	Message      string `json:"message" form:"message" binding:"required"`
	Image        string `json:"image" form:"image"`
	Attachements string `json:"attachements" form:"attachements"`
}

func AddHandler(c *gin.Context) {
	var req AddRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	params := slack.PostMessageParameters{Parse: "full", UnfurlMedia: true, Markdown: true}
	var fields []slack.AttachmentField
	attachements := strings.Split(req.Attachements, "\n")
	for _, attr := range attachements {
		arr := strings.Split(attr, "\t")
		if len(arr) != 2 {
			continue
		}
		fields = append(fields, slack.AttachmentField{
			Title: arr[0],
			Value: arr[1],
			Short: true,
		})
	}
	fields = append(fields, slack.AttachmentField{
		Title: "UserID",
		Value: strconv.FormatUint(user.Id, 10),
		Short: true,
	})
	fields = append(fields, slack.AttachmentField{
		Title: "CountryCode",
		Value: strconv.FormatUint(uint64(user.CountryCode), 10),
		Short: true,
	})
	feedbackTitle := fmt.Sprintf("Feedback from user:%d", user.Id)
	attachmentBase := slack.Attachment{
		Color:      "warning",
		AuthorName: user.ShowName,
		AuthorIcon: user.Avatar,
		Title:      feedbackTitle,
		Pretext:    req.Image,
		Text:       req.Message,
		ImageURL:   req.Image,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	attachmentFields := slack.Attachment{
		Fallback: "Fallback message",
		Fields:   fields,
	}
	params.Attachments = []slack.Attachment{attachmentBase, attachmentFields}
	channel, ts, err := Service.Slack.PostMessage(Config.Slack.FeedbackChannel, feedbackTitle, params)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	image := "NULL"
	if req.Image != "" {
		image = fmt.Sprintf("'%s'", db.Escape(req.Image))
	}
	_, _, err = db.Query(`INSERT INTO tmm.feedbacks (ts, user_id, channel, msg, image) VALUES ('%s', %d, '%s', '%s', %s)`, ts, user.Id, db.Escape(channel), db.Escape(req.Message), image)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, APIResponse{Msg: "ok"})
}
