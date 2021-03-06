/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 16:28
 **/
package _map

import NUserPb "github.com/coder2z/ndisk/pkg/pb/nuser"

type (
	/*-----------------请求------------------*/
	SMSSend struct {
		Tel  string             `validate:"required,phone" json:"tel"`
		Type NUserPb.ActionType `validate:"required" json:"type"`
	}

	AccountLogin struct {
		Account  string `validate:"required" json:"account"`
		Password string `validate:"required,min=8" json:"password"`
	}

	EmailSend struct {
		Email string             `validate:"required,email" json:"email"`
		Type  NUserPb.ActionType `validate:"required" json:"type"`
	}

	SMSLogin struct {
		Tel  string `validate:"required,phone" json:"tel"`
		Code string `validate:"required" json:"code"`
	}

	UserRegister struct {
		Name       string `validate:"required,alphanum,min=6" json:"name"`
		Alias      string `validate:"required" json:"alias"`
		Email      string `validate:"required,email" json:"email"`
		Tel        string `validate:"required,phone" json:"tel"`
		Password   string `validate:"required,max=20,min=8" json:"password"`
		RePassword string `validate:"required,max=20,min=8,eqfield=Password" json:"re_password"`
		Code       string `validate:"required" json:"code"`
	}

	RetrievePassword struct {
		Account    string `validate:"required" json:"account"`
		Password   string `validate:"required,min=8" json:"password"`
		RePassword string `validate:"required,max=20,min=8,eqfield=Password" json:"re_password"`
		Code       string `validate:"required" json:"code"`
	}

	UserToken struct {
		Token string `validate:"required" json:"token"`
	}

	BindEmail struct {
		Uid   uint64 `validate:"required" json:"uid"`
		Email string `validate:"required,email" json:"email"`
		Code  string `validate:"required" json:"code"`
	}

	/*-----------------响应------------------*/
	UserInfo struct {
		AuId        uint32 `json:"auid,omitempty"`
		Uid         uint64 `json:"uid,omitempty"`
		Name        string `json:"name,omitempty"`
		Alias       string `json:"alias,omitempty"`
		Tel         string `json:"tel,omitempty"`
		Email       string `json:"email,omitempty"`
		Status      uint32 `json:"status,omitempty"`
		Authority   string `json:"authority,omitempty"`
		EmailStatus uint32 `json:"email_status,omitempty"`
		Password    string `json:"password,omitempty"`
		CreatedAt   uint64 `json:"created_at,omitempty"`
		UpdatedAt   uint64 `json:"updated_at,omitempty"`
		DeletedAt   uint64 `json:"deleted_at,omitempty"`
	}

	Token struct {
		AccountToken string `json:"account_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}
)
