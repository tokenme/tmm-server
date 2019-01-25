package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/qiniu"
	"golang.org/x/crypto/sha3"
	"net/http"
	"time"
)

func UpTokenHandler(c *gin.Context) {
	_, exists := c.Get("USER")
	if CheckWithCode(!exists, UNAUTHORIZED_ERROR, "need login", c) {
		return
	}
	hash := sha3.Sum256([]byte(time.Now().String()))
	filename := fmt.Sprintf("%x", hash[:5]) + ".png"
	GetFileName := `tmm/task/image`
	upToken, key, link := qiniu.UpToken(Config.Qiniu, GetFileName, filename)
	c.JSON(http.StatusOK, gin.H{"uptoken": upToken, "key": key, "link": link})
}
