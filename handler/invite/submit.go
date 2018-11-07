package invite

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
	"strings"
)

func SubmitHandler(c *gin.Context) {
	inviteCode := c.Param("code")
	tel := c.PostForm("tel")
	inviteToken, err := tokenUtils.Decode(inviteCode)
	if err != nil {
		log.Error(err.Error())
	} else {
		db := Service.Db
		_, _, err = db.Query(`INSERT IGNORE INTO tmm.invite_submissions (tel, code) VALUES ('%s', %d)`, db.Escape(tel), inviteToken)
		if err != nil {
			log.Error(err.Error())
		}
	}
	downloadLink := Config.App.AndroidLink
	parser, err := uaparser.New(Config.UAParserPath)
	if err != nil {
		log.Error(err.Error())
	} else {
		client := parser.Parse(c.Request.UserAgent())
		if strings.Contains(strings.ToLower(client.Os.Family), "ios") {
			downloadLink = Config.App.IOSLink
		}
	}
	c.HTML(http.StatusOK, "invite.tmpl", gin.H{"code": inviteCode, "submitted": true, "link": downloadLink})
}
