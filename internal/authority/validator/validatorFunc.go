/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:47
 */
package validator

import (
	"github.com/go-playground/validator/v10"
)

func myValidationFunc(fl validator.FieldLevel) bool {
	return true
}
