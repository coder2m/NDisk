/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/18 15:30
 **/
package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/myxy99/component/xlog"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func RecoverMiddleware(slowQueryThresholdInMilli time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var beg = time.Now()
		var fields = make([]xlog.Field, 0, 8)
		var brokenPipe bool
		defer func() {
			fields = append(fields, zap.Float64("cost", time.Since(beg).Seconds()))
			if slowQueryThresholdInMilli > 0 {
				if cost := time.Since(beg); cost > slowQueryThresholdInMilli {
					fields = append(fields, zap.Float64("slow", cost.Seconds()))
				}
			}
			if rec := recover(); rec != nil {
				if ne, ok := rec.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				var err = rec.(error)
				fields = append(fields, zap.ByteString("stack", stack(3)))
				fields = append(fields, zap.String("err", err.Error()))
				xlog.Error("access", fields...)
				if brokenPipe {
					_ = c.Error(err)
					c.Abort()
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			fields = append(fields,
				zap.String("method", c.Request.Method),
				zap.Int("code", c.Writer.Status()),
				zap.Int("size", c.Writer.Size()),
				zap.String("host", c.Request.Host),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
				zap.String("err", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
			xlog.Info("access", fields...)
		}()
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
