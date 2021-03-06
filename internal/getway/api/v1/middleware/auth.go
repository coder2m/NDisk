/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/12 14:48
 **/
package middleware

import (
	"errors"
	R "github.com/coder2m/ndisk/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/coder2m/g-saber/xcast"
	xclient "github.com/coder2m/ndisk/internal/getway/client"
	_map "github.com/coder2m/ndisk/internal/getway/map"
	AuthorityPb "github.com/coder2m/ndisk/pkg/pb/authority"
	NUserPb "github.com/coder2m/ndisk/pkg/pb/nuser"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		ctx.Set("token", token)
		userInfo, err := xclient.NUserServer.VerifyUsers(ctx, &NUserPb.Token{
			AccountToken: token,
		})
		if err != nil {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		rolesData, err := xclient.AuthorityServer.GetUsersRoles(ctx, &AuthorityPb.Ids{
			To: []uint32{xcast.ToUint32(userInfo.Uid)},
		})
		if !errors.Is(err, nil) {
			R.HandleForbidden(ctx)
			ctx.Abort()
			return
		}
		var info = _map.UserInfo{
			Uid:         userInfo.Uid,
			Name:        userInfo.Name,
			Alias:       userInfo.Alias,
			Tel:         userInfo.Tel,
			Email:       userInfo.Email,
			Authority:   rolesData.Data[xcast.ToUint32(userInfo.Uid)],
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

func Authority() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if i, ok := ctx.Get("user"); ok {
			info := i.(_map.UserInfo)
			ctx.Set("Uid", info.Uid)
			rep, _ := xclient.AuthorityServer.Enforce(ctx, &AuthorityPb.Resources{
				Role:   xcast.ToString(info.Uid),
				Obj:    ctx.FullPath(),
				Action: ctx.Request.Method,
			})
			if rep != nil && rep.Ok {
				ctx.Next()
			} else {
				R.HandleForbidden(ctx)
				ctx.Abort()
			}
		} else {
			R.HandleForbidden(ctx)
			ctx.Abort()
		}
		return
	}
}
