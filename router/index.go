package router

import (
	"github.com/danielkov/gin-helmet"
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/common"
	"net/http"
)

func NewRouter(templatePath string, config common.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(helmet.Default())
	xssMdlwr := &xss.XssMw{
		FieldsToSkip: []string{"password", "start_date", "end_date", "token"},
		BmPolicy:     "UGCPolicy",
	}
	r.Use(xssMdlwr.RemoveXss())
	r.LoadHTMLGlob(templatePath)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "tokenmama.io"})
		return
	})
	r.GET("/contact/list", func(c *gin.Context) {
		c.JSON(http.StatusOK, []gin.H{{
			"name":  "telegram",
			"value": config.Contact.Telegram,
		}, {
			"name":  "wechat",
			"value": config.Contact.Wechat,
		}, {
			"name":  "website",
			"value": config.Contact.Website,
		}})
		return
	})
	authRouter(r)
	userRouter(r)
	deviceRouter(r)
	appRouter(r)
	taskRouter(r)
	exchangeRouter(r)
	tokenRouter(r)
	qiniuRouter(r)
	orderbookRouter(r)
	return r
}
