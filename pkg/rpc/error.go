/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 18:37
 **/
package xrpc

import (
	"github.com/myxy99/component/pkg/xcode"
)

const (
	//system error
	MysqlErr = iota + 100

	//业务错误 高于 10000
	EmptyData = iota + 10000
	AccountLoginErrCode
	ValidationErrCode
	FrequentOperationErrCode
	SMSSendErrCode
	SMSLoginErrCode
	SendEmailErrCode
	UserRegisterErrCode
	RetrievePasswordErrCode
	GetUserByIdErrCode
	GetUserListErrCode
	UpdateUserStatusErrCode
	UpdateUserErrCode
	DelUsersErrCode
	RecoverDelUsersErrCode
	CreateUsersErrCode
	VerifyUsersTokenErrCode
	RefreshTokenErrCode
)

var (
	systemErrDepot = []xcode.CodeInfo{
		{xcode.SystemType, MysqlErr, "mysql 错误"},
	}
	businessErrDepot = []xcode.CodeInfo{
		{xcode.BusinessType, EmptyData, "数据未找到"},
		{xcode.BusinessType, AccountLoginErrCode, "AccountLogin"},
		{xcode.BusinessType, ValidationErrCode, "validation error"},
		{xcode.BusinessType, FrequentOperationErrCode, "frequent operation"},
		{xcode.BusinessType, SMSSendErrCode, "SMS Send Err"},
		{xcode.BusinessType, SMSLoginErrCode, "SMS Login Err"},
		{xcode.BusinessType, SendEmailErrCode, "Send Email Err"},
		{xcode.BusinessType, UserRegisterErrCode, "User Register Err"},
		{xcode.BusinessType, RetrievePasswordErrCode, "User Register Err"},
		{xcode.BusinessType, GetUserByIdErrCode, "GetUserByIdErrCode Err"},
		{xcode.BusinessType, GetUserListErrCode, "GetUserListErrCode Err"},
		{xcode.BusinessType, UpdateUserStatusErrCode, "UpdateUserStatus Err"},
		{xcode.BusinessType, UpdateUserErrCode, "UpdateUser Err"},
		{xcode.BusinessType, DelUsersErrCode, "DelUsers Err"},
		{xcode.BusinessType, RecoverDelUsersErrCode, "RecoverDelUsers Err"},
		{xcode.BusinessType, CreateUsersErrCode, "CreateUsers Err"},
		{xcode.BusinessType, VerifyUsersTokenErrCode, "VerifyUsers Token Err"},
		{xcode.BusinessType, RefreshTokenErrCode, "RefreshToken Token Err"},
	}
)

func init() {
	xcode.CodeAdds(systemErrDepot)
	xcode.CodeAdds(businessErrDepot)
}
