/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:20
 **/
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/ndisk/internal/getway/error/httpError"
	_map "github.com/myxy99/ndisk/internal/getway/map"
)

//  账号登录
func AccountLogin(ctx *gin.Context) {
	var page _map.AccountLogin
	if err := ctx.ShouldBind(&page); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(page); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}

	return
}

//  发送 短信

//  账号登录-验证码

//  注册

//  注册邮件验证发送邮件
//  找回密码发送邮件验证

//  找回密码 邮件或者电话 都可

//  根据id获取用户信息

//  批量获取用户信息 分页

//  修改用户状态

//  修改用户信息

//  批量删除用户

//  批量恢复用户

//  批量添加用户

//  用户校验

//  刷新token
