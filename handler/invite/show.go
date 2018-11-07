package invite

import (
	//"github.com/davecgh/go-spew/spew"
	//"github.com/mkideal/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowHandler(c *gin.Context) {
	inviteCode := c.Param("code")
	c.HTML(http.StatusOK, "invite.tmpl", gin.H{"code": inviteCode})
}
