/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/18 15:38
 **/
package middleware

import (
	"net/http"
	"time"

	"github.com/coder2m/component/xmonitor"
	"github.com/gin-gonic/gin"
)

func extractAID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("AID")
}

func XMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		beg := time.Now()
		c.Next()

		xmonitor.ServerHandleHistogram.WithLabelValues(
			xmonitor.TypeHTTP,
			c.Request.Method+"."+c.Request.URL.Path,
			extractAID(c),
		).Observe(time.Since(beg).Seconds())

		xmonitor.ServerHandleCounter.WithLabelValues(
			xmonitor.TypeHTTP,
			c.Request.Method+"."+c.Request.URL.Path,
			extractAID(c),
			http.StatusText(c.Writer.Status()),
		).Inc()
		return
	}
}
