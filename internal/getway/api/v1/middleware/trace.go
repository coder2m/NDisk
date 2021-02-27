/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/18 15:45
 **/
package middleware

import (
	"github.com/coder2m/component/xtrace"
	"github.com/gin-gonic/gin"
)

func XTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := xtrace.StartSpanFromContext(
			c.Request.Context(),
			c.Request.Method+" "+c.Request.URL.Path,
			xtrace.TagComponent("http"),
			xtrace.TagSpanKind("server"),
			xtrace.HeaderExtractor(c.Request.Header),
			xtrace.CustomTag("http.url", c.Request.URL.Path),
			xtrace.CustomTag("http.method", c.Request.Method),
			xtrace.CustomTag("peer.ipv4", c.ClientIP()),
		)
		c.Request = c.Request.WithContext(ctx)
		defer span.Finish()
		c.Next()
	}
}
