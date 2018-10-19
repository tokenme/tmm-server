package blowup

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/nu7hatch/gouuid"
	. "github.com/tokenme/tmm/handler"
	"io"
)

func NotifyHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	uid, err := uuid.NewV4()
	if CheckErr(err, c) {
		return
	}
	BlowupService.NewChannel(uid.String())
	c.Stream(func(w io.Writer) bool {
		ch := BlowupService.Channel(uid.String())
		if ch == nil {
			log.Warn("Not found ev")
			return false
		}
		select {
		case ev := <-ch:
			c.SSEvent("message", ev)
		case <-ExitCh:
			return false
		}
		return true
	})
	BlowupService.CloseChannel(uid.String())
}
