package user

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/govatar"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"image"
	"image/png"
	"net/http"
)

func AvatarGetHandler(c *gin.Context) {
	userContext, exists := c.Get("USER")
	var (
		img image.Image
		err error
	)
	if !exists {
		value := c.Param("key")
		img, err = govatar.GenerateFromUsername(govatar.MALE, value)
	} else {
		user := userContext.(common.User)
		img, err = govatar.GenerateFromUsername(govatar.MALE, user.ShowName)
	}
	if CheckErr(err, c) {
		return
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if CheckErr(err, c) {
		return
	}
	c.Writer.Header().Add("Cache-Control", "max-age=86400")
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}
