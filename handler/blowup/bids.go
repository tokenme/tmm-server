package blowup

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/blowup"
	"github.com/tokenme/tmm/utils"
	"net/http"
)

func BidsHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	db := Service.Db
	rows, _, err := db.Query(`SELECT bb.user_id, u.country_code, u.mobile, bb.points, bb.rate/100, bb.escaped FROM tmm.blowup_bids AS bb INNER JOIN ucoin.users AS u ON (u.id=bb.user_id) WHERE NOT EXISTS (SELECT 1 FROM tmm.blowup_sessions AS bs WHERE bs.id=bb.session_id LIMIT 1) ORDER BY bb.inserted_at DESC LIMIT 30`)
	if CheckErr(err, c) {
		return
	}
	var events []*blowup.Event
	for _, row := range rows {
		var (
			countryCode = row.Uint(1)
			mobile      = utils.HideMobile(row.Str(2))
			points, _   = decimal.NewFromString(row.Str(3))
			rate, _     = decimal.NewFromString(row.Str(4))
			escaped     = row.Int(5) == 1
			eventType   = blowup.BiddingEvent
		)

		if escaped {
			eventType = blowup.EscapeEvent
		}

		ev := &blowup.Event{
			Nick:  fmt.Sprintf("+%d%s", countryCode, mobile),
			Type:  eventType,
			Value: points,
			Rate:  rate,
		}
		events = append(events, ev)
	}
	c.JSON(http.StatusOK, events)
}
