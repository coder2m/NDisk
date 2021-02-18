/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:20
 **/
package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/ndisk/internal/getway/error/httpError"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	"github.com/myxy99/ndisk/internal/getway/server/auth_server"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
	R "github.com/myxy99/ndisk/pkg/response"
)

//  账号登录
func AccountLogin(ctx *gin.Context) {
	var login _map.AccountLogin
	if err := ctx.ShouldBind(&login); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(login); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := auth_server.AccountLogin(ctx, login); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//  短信验证码发送验证码
func LoginSMSSend(ctx *gin.Context) {
	var send _map.SMSSend
	if err := ctx.ShouldBind(&send); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	send.Type = NUserPb.ActionType_Login_Type
	if err := xvalidator.Struct(send); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := auth_server.SMSSend(ctx, send); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//  账号登录-验证码
func SMSLogin(ctx *gin.Context) {
	var login _map.SMSLogin
	if err := ctx.ShouldBind(&login); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(login); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := auth_server.SMSLogin(ctx, login); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//  注册
func Register(ctx *gin.Context) {
	var register _map.UserRegister
	if err := ctx.ShouldBind(&register); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(register); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := auth_server.Register(ctx, register); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 注册发送验证码
func RegisterSendSMS(ctx *gin.Context) {
	var send _map.SMSSend
	if err := ctx.ShouldBind(&send); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	send.Type = NUserPb.ActionType_Register_Type
	if err := xvalidator.Struct(send); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := auth_server.SMSSend(ctx, send); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//  找回密码发送sms
func RetrieveSendSms(ctx *gin.Context) {
	var send _map.SMSSend
	if err := ctx.ShouldBind(&send); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	send.Type = NUserPb.ActionType_Retrieve_Type
	if err := xvalidator.Struct(send); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := auth_server.SMSSend(ctx, send); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//  找回密码 邮件或者电话 都可
func Retrieve(ctx *gin.Context) {
	var retrieve _map.RetrievePassword
	if err := ctx.ShouldBind(&retrieve); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(retrieve); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := auth_server.Retrieve(ctx, retrieve); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//  刷新token
func RefreshToken(ctx *gin.Context) {
	var token _map.UserToken
	if t, ok := ctx.Get("token"); ok {
		token.Token = t.(string)
	} else {
		httpError.HandleForbidden(ctx, nil)
		return
	}
	if err := xvalidator.Struct(token); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := auth_server.RefreshToken(ctx, token); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 获取用户信息
func Info(ctx *gin.Context) {
	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		menu, _ := auth_server.GetPermissionAndMenuByRoles(ctx, strings.Split(info.Authority, ","))
		R.Ok(ctx, gin.H{
			"info": info,
			"menu": menu,
		})
	} else {
		httpError.HandleForbidden(ctx, nil)
	}
	return
}

// 绑定邮件 发送邮件验证码
func BindEmailSend(ctx *gin.Context) {
	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		if err := auth_server.SendEmail(ctx, _map.EmailSend{
			Email: info.Email,
			Type:  NUserPb.ActionType_EmailAttest_Type,
		}); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, nil)
		}
	} else {
		httpError.HandleForbidden(ctx, nil)
	}
	return
}

//绑定邮箱
func BindEmail(ctx *gin.Context) {
	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)

		var bind _map.BindEmail
		if err := ctx.ShouldBind(&bind); err != nil {
			httpError.HandleBadRequest(ctx, nil)
			return
		}
		bind.Uid = info.Uid
		bind.Email = info.Email
		if err := xvalidator.Struct(bind); err != nil {
			httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
			return
		}

		if err := auth_server.BindEmail(ctx, _map.BindEmail{
			Uid:   bind.Uid,
			Email: bind.Email,
			Code:  bind.Code,
		}); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, nil)
		}
	} else {
		httpError.HandleForbidden(ctx, nil)
	}
	return
}
