package invite

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
	"strings"
)

func ShowHandler(c *gin.Context) {
	inviteCode := c.Param("code")
	submitted := c.Query("submitted")
	if submitted != "" {
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
		return
	}
	c.HTML(http.StatusOK, "invite.tmpl", gin.H{"code": inviteCode})
}
