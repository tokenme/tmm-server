package ad

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"net/http"
)

func ClkHandler(c *gin.Context) {
	code := c.Param("code")
	creative, err := common.DecodeCreative([]byte(Config.LinkSalt), code)
	if CheckErr(err, c) {
		return
	}
	db := Service.Db
	_, _, err = db.Query(`INSERT INTO tmm.creative_stats (creative_id, record_on, clk) VALUES (%d, NOW(), 1) ON DUPLICATE KEY UPDATE clk=clk+1`, creative.Id)
	if err != nil {
		log.Error(err.Error())
	}
	c.Redirect(http.StatusFound, creative.Link)
}
