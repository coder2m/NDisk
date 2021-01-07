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
	Password string `validate:"required,min=8"`
}

type SMSLogin struct {
	Tel  string `validate:"e164,required"`
	Code string `validate:"required"`
}

type UserRegister struct {
	Name     string `validate:"required,alphanum"`
	Alias    string `validate:"required"`
	Email    string `validate:"required,email"`
	Tel      string `validate:"required,e164"`
	Password string `validate:"required,min=8"`
	Code     string `validate:"required"`
}

type RetrievePassword struct {
	Account  string `validate:"required"`
	Password string `validate:"required,min=8"`
	Code     string `validate:"required"`
}
