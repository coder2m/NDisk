package xclient

import (
	xemail "github.com/myxy99/component/xinvoker/email"
	xsms "github.com/myxy99/component/xinvoker/sms"
)

var (
	emailMain *xemail.Email
	smsMain   *xsms.Client
)

func EmailMain() *xemail.Email {
	if emailMain == nil {
		emailMain = xemail.Invoker("main")
	}
	return emailMain
}

func SMSMain() *xsms.Client {
	if smsMain == nil {
		smsMain = xsms.Invoker("main")
	}
	return smsMain
}
