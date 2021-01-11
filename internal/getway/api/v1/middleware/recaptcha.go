/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 18:11
 **/
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/myxy99/ndisk/internal/nfile/error/httpError"
	"github.com/myxy99/ndisk/pkg/recaptcha"
)

func Recaptcha() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !recaptcha.Verify(ctx.GetHeader("captcha")).Success {
			ctx.Abort()
			httpError.HandleCaptchaError(ctx)
			return
		}
		ctx.Next()
		return
	}
}
