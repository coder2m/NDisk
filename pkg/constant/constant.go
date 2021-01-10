package constant

import (
	"fmt"
	"time"
)

type Key string

func (k Key) Format(age ...interface{}) string {
	return fmt.Sprintf(string(k), age...)
}

// redis
const (
	SendVerificationCode Key = `NUser_SendSMS_%v_%v` //操作类型 用户电话号码
)

// 验证码
const (
	VerificationEffectiveTime = time.Minute //验证码持续时间
	VerificationCodeLength    = 6           //验证码长度
)

//server
const (
	DefaultNamespaces     = `NDisk`
	GRPCTargetEtcd    Key = `etcd://%s/%s` //命名空间 //servername
)
