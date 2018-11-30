package user


import ("github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"net/http"
	"github.com/tokenme/tmm/common"
)

func GetUserInfoHandler(c *gin.Context){
	userContext, exists := c.Get("USER")
	if Check(!exists, `Need login`, c) {
		return
	}
	user:=userContext.(common.User)
	c.JSON(http.StatusOK,gin.H{
		`code`:http.StatusOK,
		`msg`:"",
		"data":user,})
}