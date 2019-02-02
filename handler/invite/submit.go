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
    isFamily := c.PostForm("is_family")
	inviteToken, err := tokenUtils.Decode(inviteCode)
	if err != nil {
		log.Error(err.Error())
	} else {
		db := Service.Db
        rows, _, err := db.Query(`SELECT COUNT(1) FROM tmm.invite_submissions AS iss WHERE iss.code=%d AND iss.is_family=1`, inviteToken)
        if err != nil {
            log.Error(err.Error())
        }
        var familyInvites uint
        if len(rows) > 0 {
            familyInvites = rows[0].Uint(0)
            if familyInvites >= 10 {
                isFamily = "0"
            }
        }
		_, _, err = db.Query(`INSERT IGNORE INTO tmm.invite_submissions (tel, code, is_family) VALUES ('%s', %d, '%s')`, db.Escape(tel), inviteToken, db.Escape(isFamily))
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
    if len(isFamily) > 0 {
        c.JSON(http.StatusOK, gin.H{"code": inviteCode, "submitted": true, "link": downloadLink})
    } else {
	    c.HTML(http.StatusOK, "invite.tmpl", gin.H{"code": inviteCode, "submitted": true, "link": downloadLink})
    }
}
