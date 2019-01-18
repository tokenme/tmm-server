package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tokenme/tmm/handler/admin/verify"
	"github.com/tokenme/tmm/handler/auth"
	"github.com/tokenme/tmm/middlewares/jwt"
	"time"
)

var ADMIN_AUTH_KEY = []byte("20eefe8d82asdfasdfvasdfeba3ca8a417e14a48d24632bc35bbd7")

var AdminAuthMiddleware = &jwt.GinJWTMiddleware{
	Realm:         AUTH_REALM,
	Key:           ADMIN_AUTH_KEY,
	Timeout:       AUTH_TIMEOUT,
	MaxRefresh:    AUTH_MAXREFRESH,
	Authenticator: auth.AuthenticatorFunc,
	Authorizator:  auth.AuthorizatorFunc,
	Unauthorized: func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	},
	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup: "header:Authorization",
	// TokenLookup: "query:token",
	// TokenLookup: "cookie:token",

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName: "Bearer",

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc: time.Now,
}

func AdminRouter(r *gin.Engine) {
	r.POST(`/admin/auth/login`, AdminAuthMiddleware.LoginHandler,verify.RecordLoginHandler)
	AricleRouter(r)
	TaskRouter(r)
	UserRouter(r)
	InfoRouter(r)
	AdRouter(r)
	AccountRouter(r)
	WithdrawRouter(r)
	AuditRouter(r)
}
