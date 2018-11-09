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
	r.Use(helmet.NoSniff(), helmet.DNSPrefetchControl(), helmet.FrameGuard("ALLOW-FROM https://tmm.tokenmama.io"), helmet.SetHSTS(true), helmet.IENoOpen(), helmet.XSSFilter())
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
	r.GET("/ios/download", func(c *gin.Context) {
		c.Redirect(http.StatusFound, config.App.IOSLink)
		return
	})
	r.GET("/android/download", func(c *gin.Context) {
		c.Redirect(http.StatusFound, config.App.AndroidLink)
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
	redeemRouter(r)
	feedbackRouter(r)
	slackRouter(r)
	bonusRouter(r)
	blowupRouter(r)
	articleRouter(r)
	inviteRouter(r)
	goodRouter(r)
	return r
}
