/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:23
 **/
package registry

import (
	ah "github.com/myxy99/ndisk/internal/getway/api/v1/handle/auth"
	"github.com/myxy99/ndisk/internal/getway/api/v1/middleware"
)

func init() {
	user := V1().Group("/auth")
	//账号密码登录
	user.POST("/login", middleware.Recaptcha(), ah.AccountLogin)
	//手机验证码登录发送验证码
	user.POST("/login/sms/send", middleware.Recaptcha(), ah.LoginSMSSend)
	//手机验证码登录
	user.POST("/login/sms", middleware.Recaptcha(), ah.SMSLogin)
	//用户注册
	user.POST("/register", ah.Register)
	//用户注册发送验证码
	user.POST("/register/sms/send", middleware.Recaptcha(), ah.RegisterSendSMS)

	//找回密码发送短信
	user.POST("/retrieve/sms/send", middleware.Recaptcha(), ah.RetrieveSendSms)
	//找回密码 修改密码
	user.POST("/retrieve", ah.Retrieve)

	//刷新token
	user.POST("/refresh/token", middleware.Auth(), ah.RefreshToken)
	//获取用户信息
	user.GET("/info", middleware.Auth(), ah.Info)

	//认证邮箱 发送邮件验证码
	user.POST("/bind/email", middleware.Auth(), ah.BindEmailSend)
	//认证邮箱
	user.POST("/bind", middleware.Auth(), ah.BindEmail)
}
