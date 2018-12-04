package app

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
	"strings"
)

func DownloadHandler(c *gin.Context) {
	downloadLink := Config.App.AndroidLink
	parser, err := uaparser.New(Config.UAParserPath)
	var platform uint8
	if err != nil {
		log.Error(err.Error())
	} else {
		client := parser.Parse(c.Request.UserAgent())
		if strings.Contains(strings.ToLower(client.Os.Family), "ios") {
			downloadLink = Config.App.IOSLink
			platform = 1
		} else if strings.Contains(strings.ToLower(client.Os.Family), "android") {
			platform = 2
		}
	}
	utmz := c.Query("__utmz")
	if utmz != "" {
		db := Service.Db
		_, _, err := db.Query(`INSERT INTO tmm.app_download_stats (utmz, platform, inserted_on) VALUES ('%s', %d, NOW())`, db.Escape(utmz), platform)
		if err != nil {
			log.Error(err.Error())
		}
	}
	c.Redirect(http.StatusFound, downloadLink)
}
