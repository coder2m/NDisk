/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/12 14:48
 **/
package middleware

import (
	"github.com/gin-gonic/gin"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	"github.com/myxy99/ndisk/internal/getway/error/httpError"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			httpError.HandleForbidden(ctx, nil)
			ctx.Abort()
			return
		}
		ctx.Set("token", token)
		userInfo, err := xclient.NUserServer.VerifyUsers(ctx, &NUserPb.Token{
			AccountToken: token,
		})
		if err != nil {
			httpError.HandleForbidden(ctx, nil)
			ctx.Abort()
			return
		}
		var info = _map.UserInfo{
			Uid:         userInfo.Uid,
			Name:        userInfo.Name,
			Alias:       userInfo.Alias,
			Tel:         userInfo.Tel,
			Email:       userInfo.Email,
			Status:      userInfo.Status,
			EmailStatus: userInfo.EmailStatus,
			CreatedAt:   userInfo.CreatedAt,
			UpdatedAt:   userInfo.UpdatedAt,
		}
		ctx.Set("user", info)
		ctx.Next()
		return
	}
}
