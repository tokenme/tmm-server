package qiniu

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/qiniu"
	"net/http"
	"strconv"
	"time"
)

func UpTokenHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	upToken, key, link := qiniu.UpToken(Config.Qiniu, Config.Qiniu.ImagePath, timestamp)
	c.JSON(http.StatusOK, gin.H{"uptoken": upToken, "key": key, "link": link})
}
