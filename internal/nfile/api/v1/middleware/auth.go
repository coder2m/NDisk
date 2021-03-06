/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/12 14:48
 **/
package middleware

import (
	R "github.com/coder2z/ndisk/pkg/response"

	xclient "github.com/coder2z/ndisk/internal/getway/client"
	NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		userInfo, err := xclient.NUserServer.VerifyUsers(ctx, &NUserPb.Token{
			AccountToken: token,
		})
		if err != nil {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("Uid", userInfo.Uid)
		ctx.Next()
		return
	}
}
