/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:46
 */
package validator

import (
	"github.com/coder2z/g-saber/xvalidator"
)

var RegisterValidation = map[string]*xvalidator.Register{
	"phone": {
		phoneValidationFunc,
		"电话号码验证不通过",
	},
	"oneOrTwo": {
		oneOrTwo,
		"status只能是1或者2",
	},
}
