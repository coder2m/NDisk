/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:23
 **/
package registry

import (
	ah "github.com/coder2m/ndisk/internal/getway/api/v1/handle/auth"
	"github.com/coder2m/ndisk/internal/getway/api/v1/middleware"
)

func init() {
	auth := V1().Group("/auth")
	//账号密码登录
	auth.POST("/login", middleware.Recaptcha(), ah.AccountLogin)
	//手机验证码登录发送验证码
	auth.GET("/login/sms", middleware.Recaptcha(), ah.LoginSMSSend)
	//手机验证码登录
	auth.POST("/login/sms", middleware.Recaptcha(), ah.SMSLogin)
	//用户注册
	auth.POST("/register", ah.Register)
	//用户注册发送验证码
	auth.GET("/register/sms", middleware.Recaptcha(), ah.RegisterSendSMS)
	//找回密码发送短信
	auth.GET("/retrieve/sms", middleware.Recaptcha(), ah.RetrieveSendSms)
	//找回密码 修改密码
	auth.POST("/retrieve", ah.Retrieve)
	//刷新token
	auth.POST("/refresh/token", middleware.Auth(), ah.RefreshToken)
	//获取用户信息
	auth.GET("/info", middleware.Auth(), ah.Info)
	//认证邮箱 发送邮件验证码
	auth.POST("/bind/email", middleware.Auth(), ah.BindEmailSend)
	//认证邮箱
	auth.POST("/bind", middleware.Auth(), ah.BindEmail)
}
