package router

import (
	"github.com/danielkov/gin-helmet"
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(templatePath string) *gin.Engine {
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
	authRouter(r)
	userRouter(r)
	deviceRouter(r)
	appRouter(r)
	taskRouter(r)
	exchangeRouter(r)
	tokenRouter(r)
	return r
}
