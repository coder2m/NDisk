package middleware

import (
	"crypto/subtle"
	R "github.com/coder2z/ndisk/pkg/response"

	"github.com/coder2z/ndisk/pkg/constant"
	xrand "github.com/coder2z/ndisk/pkg/rand"
	"github.com/gin-gonic/gin"
)

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(constant.CSRF)
		var token string
		if err != nil {
			token = xrand.CreateRandomString(32)
		} else {
			token = cookie
		}
		switch c.Request.Method {
		case "GET", "OPTIONS", "HEAD":
		default:
			clientToken := c.Request.Header.Get("X-CSRF-TOKEN")
			if clientToken == "" {
				R.HandleForbidden(c)
				c.Abort()
				return
			}
			if !validateCSRFToken(token, clientToken) {
				R.HandleForbidden(c)
				c.Abort()
				return
			}
		}
		domain := c.Request.Header.Get("X-Server-Domain")
		c.SetCookie(constant.CSRF, token, 86400, "/", "."+domain, false, false)
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
