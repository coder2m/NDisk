package xclient

import (
	xemail "github.com/coder2m/component/xinvoker/email"
	xsms "github.com/coder2m/component/xinvoker/sms"
	"github.com/coder2m/ndisk/internal/nuser/model"
	"github.com/coder2m/ndisk/internal/nuser/server/token"
	redisToken "github.com/coder2m/ndisk/internal/nuser/server/token/redis"
)

var (
	emailMain        *xemail.Email
	smsMain          *xsms.Client
	redisAccessToken token.AccessToken
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

func RedisToken() token.AccessToken {
	if redisAccessToken == nil {
		redisAccessToken = redisToken.NewAccessToken(model.MainRedis())
	}
	return redisAccessToken
}
