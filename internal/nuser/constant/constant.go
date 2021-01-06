/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 15:47
 **/
package constant

import "fmt"

type RedisKey string

func (k RedisKey) Format(age ...interface{}) string {
	return fmt.Sprintf(string(k), age...)
}

const (
	SendSMS RedisKey = `NUser_SendSMS_%v_%v` //操作类型 用户电话号码
)
