package constant

import (
	"fmt"
	"time"
)

type RedisKey string

func (k RedisKey) Format(age ...interface{}) string {
	return fmt.Sprintf(string(k), age...)
}

const (
	SendVerificationCode RedisKey = `NUser_SendSMS_%v_%v` //操作类型 用户电话号码
)

const (
	VerificationEffectiveTime = time.Minute //验证码持续时间
	VerificationCodeLength    = 6           //验证码长度
)
const (
	DefaultNamespaces = `NDisk`
)
