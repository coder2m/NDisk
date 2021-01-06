/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 19:14
 **/
package _map

type Phone struct {
	Number string `validate:"e164,required"`
}

type Email struct {
	Email string `validate:"email,required"`
}

type AccountLogin struct {
	Account  string `validate:"required"`
	Password string `validate:"required"`
}

type SMSLogin struct {
	Tel  string `validate:"e164,required"`
	Code string `validate:"required"`
}
