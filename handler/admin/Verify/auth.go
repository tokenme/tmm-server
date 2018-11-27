package Verify

import (
	"github.com/gin-gonic/gin"
	"net/http"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
)

func VerifyAdminFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext, exists := c.Get("USER")
		if Check(!exists, `Need login`, c) {
			return
		}
		user := userContext.(common.User)
		if user.IsAdmin() {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				`code`: `401`,
				`msg`:  `The User Must Be Admin`,
			})
			return
		}
		c.Next()
		return
	}
}

