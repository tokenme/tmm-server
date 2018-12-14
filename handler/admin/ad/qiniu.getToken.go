package ad

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"strconv"
	"time"
	"github.com/tokenme/tmm/tools/qiniu"
	"net/http"
	"github.com/tokenme/tmm/common"
	"fmt"
)

func GetUpdateToken(c *gin.Context) {
	userContext, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	user := userContext.(common.User)
	_type := c.Query(`type`)
	if Check(_type != "1" && _type != "2", `Invalid param`, c) {
		return
	}
	var path string
	fileName := fmt.Sprintf("%d-%s", user.Id, strconv.FormatInt(time.Now().UnixNano(), 10))
	if _type == "1" {
		path = "ad/creative"
	} else {
		path = "ad/creative/share"
	}
	upToken, key, link := qiniu.UpToken(Config.Qiniu, path, fileName)

	c.JSON(http.StatusOK, gin.H{"uptoken": upToken, "key": key, "link": link})
}
