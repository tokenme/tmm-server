package feedback

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nlopes/slack"
	"github.com/panjf2000/ants"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"sync"
)

func ListHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	db := Service.Db
	rows, _, err := db.Query(`SELECT ts, channel, msg, image FROM tmm.feedbacks WHERE user_id=%d ORDER BY ts DESC LIMIT 5`, user.Id)
	if CheckErr(err, c) {
		return
	}
	var feedbacks []*common.Feedback

	var wg sync.WaitGroup
	pool, _ := ants.NewPoolWithFunc(10000, func(req interface{}) error {
		defer wg.Done()
		fb := req.(*common.Feedback)
		params := &slack.GetConversationRepliesParameters{
			ChannelID: fb.Channel,
			Timestamp: fb.Ts,
		}
		msgs, _, _, err := Service.Slack.GetConversationRepliesContext(c, params)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		if len(msgs) > 1 {
			var replies []common.Feedback
			for _, msg := range msgs[1:] {
				replies = append(replies, common.Feedback{
					Ts:  msg.Timestamp,
					Msg: msg.Text,
				})
			}
			fb.Replies = replies
		}
		return nil
	})

	for _, row := range rows {
		fb := &common.Feedback{
			Ts:      row.Str(0),
			Channel: row.Str(1),
			Msg:     row.Str(2),
			Image:   row.Str(3),
		}
		feedbacks = append(feedbacks, fb)
		wg.Add(1)
		pool.Serve(fb)
	}
	wg.Wait()
	c.JSON(http.StatusOK, feedbacks)
}
