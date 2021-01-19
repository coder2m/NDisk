/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 19:14
 **/
package _map

import NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"

type Phone struct {
	Number string `validate:"required,phone"`
}

type Email struct {
	Email string `validate:"required,email"`
}

type AccountLogin struct {
	Account  string `validate:"required"`
	Password string `validate:"required,min=8"`
}

type SMSLogin struct {
	Tel  string `validate:"required,phone"`
	Code string `validate:"required"`
}

type UserRegister struct {
	Name       string `validate:"required,alphanum,min=6"`
	Alias      string `validate:"required"`
	Email      string `validate:"required,email"`
	Tel        string `validate:"required,phone"`
	Password   string `validate:"required,max=20,min=8"`
	RePassword string `validate:"required,max=20,min=8,eqfield=Password" json:"re_password"`
	Code       string `validate:"required"`
}

type RetrievePassword struct {
	Account  string `validate:"required"`
	Password string `validate:"required,min=8"`
	Code     string `validate:"required"`
}

type UpdateUserStatus struct {
	Uid    uint64 `validate:"required,number,min=1"`
	Status uint32 `validate:"required,number,min=1"`
}

type UpdateUser struct {
	Uid      uint64 `validate:"required,number,min=1"`
	Name     string `validate:"required,alphanum"`
	Alias    string `validate:"required"`
	Email    string `validate:"required,email"`
	Tel      string `validate:"required,phone"`
	Password string
}

type UserToken struct {
	Token string `validate:"required"`
}

type CheckCode struct {
	Account string             `validate:"required"`
	Code    string             `validate:"required"`
	Type    NUserPb.ActionType `validate:"required,number,min=1"` //0为注册邮件验证发送邮件；1为找回密码发送邮件验证 2为登录 3为邮箱验证
}

type UserInfo struct {
	AUId        uint   `json:"auid,omitempty"` //关联表id
	Uid         uint   `json:"uid,omitempty"`
	Name        string `json:"name,omitempty"`
	Alias       string `json:"alias,omitempty"`
	Tel         string `json:"tel,omitempty"`
	Email       string `json:"email,omitempty"`
	Status      uint32 `json:"status,omitempty"`
	EmailStatus uint32 `json:"email_status,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
	DeletedAt   int64  `json:"deleted_at,omitempty"`
}
