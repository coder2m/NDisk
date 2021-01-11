/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:47
 */
package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func phoneValidationFunc(fl validator.FieldLevel) bool {
	reg := `^(0|86|17951)?(13[0-9]|15[012356789]|166|17[3678]|18[0-9]|14[57])[0-9]{8}$`
	return regexp.MustCompile(reg).MatchString(fl.Field().String())
}