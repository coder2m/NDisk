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
	ValidationErrCode

	//NUser
	AccountLoginErrCode
	LoginUserBanErrCode
	FrequentOperationErrCode
	SMSSendErrCode
	SMSLoginErrCode
	SendEmailErrCode
	UserRegisterErrCode
	RetrievePasswordErrCode
	GetUserByIdErrCode
	GetUserListErrCode
	UpdateUserStatusErrCode
	UpdateUserEmailStatusErrCode
	UpdateUserErrCode
	DelUsersErrCode
	RecoverDelUsersErrCode
	GetUserListByUidErrCode
	CreateUsersErrCode
	VerifyUsersTokenErrCode
	RefreshTokenErrCode
	MaximumNumberErrCode
	DataExistErrCode

	//Authority
	DeleteRoleErrCode
	GetRolesForUserErrCode
	AddRolesForUserErrCode
	HasRoleForUserErrCode
	DeleteRoleForUserErrCode
	DeleteUserErrCode
	DeleteRolesForUserErrCode
	GetUsersForRoleErrCode
	AddPermissionForUserErrCode
	GetPermissionsForUserErrCode
	DeletePermissionForUserErrCode
	DeletePermissionsForUserErrCode
	DeletePermissionErrCode
	HasPermissionForUserErrCode
	EnforceErrCode
)

var (
	systemErrDepot = []xcode.CodeInfo{
		{xcode.SystemType, MysqlErr, "mysql 错误"},
	}
	businessErrDepot = []xcode.CodeInfo{
		{xcode.BusinessType, EmptyData, "Empty Data"},
		{xcode.BusinessType, AccountLoginErrCode, "AccountLogin"},
		{xcode.BusinessType, LoginUserBanErrCode, "User Ban"},
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
		{xcode.BusinessType, UpdateUserEmailStatusErrCode, "UpdateUserEmailStatus Err"},
		{xcode.BusinessType, UpdateUserErrCode, "UpdateUser Err"},
		{xcode.BusinessType, DelUsersErrCode, "DelUsers Err"},
		{xcode.BusinessType, RecoverDelUsersErrCode, "RecoverDelUsers Err"},
		{xcode.BusinessType, CreateUsersErrCode, "CreateUsers Err"},
		{xcode.BusinessType, VerifyUsersTokenErrCode, "VerifyUsers Token Err"},
		{xcode.BusinessType, RefreshTokenErrCode, "RefreshToken Token Err"},
		{xcode.BusinessType, MaximumNumberErrCode, "Maximum number of operations exceeded error"},
		{xcode.BusinessType, DataExistErrCode, "data exist"},
		{xcode.BusinessType, GetUserListByUidErrCode, "GetUserListByUid error"},


		{xcode.BusinessType, DeleteRoleErrCode, "delete role error"},
		{xcode.BusinessType, GetRolesForUserErrCode, "get roles for user error"},
		{xcode.BusinessType, AddRolesForUserErrCode, "AddRolesForUse error"},
		{xcode.BusinessType, HasRoleForUserErrCode, "HasRoleForUserErrCode error"},
		{xcode.BusinessType, DeleteRoleForUserErrCode, "DeleteRoleForUserErrCode error"},
		{xcode.BusinessType, DeleteUserErrCode, "DeleteUser error"},
		{xcode.BusinessType, DeleteRolesForUserErrCode, "DeleteRolesForUser error"},
		{xcode.BusinessType, GetUsersForRoleErrCode, "GetUsersForRoleErr error"},
		{xcode.BusinessType, AddPermissionForUserErrCode, "AddPermissionForUser error"},
		{xcode.BusinessType, GetPermissionsForUserErrCode, "GetPermissionsForUser error"},
		{xcode.BusinessType, DeletePermissionForUserErrCode, "DeletePermissionForUser error"},
		{xcode.BusinessType, DeletePermissionsForUserErrCode, "DeletePermissionsForUser error"},
		{xcode.BusinessType, DeletePermissionErrCode, "DeletePermission error"},
		{xcode.BusinessType, HasPermissionForUserErrCode, "HasPermissionForUser error"},
		{xcode.BusinessType, EnforceErrCode, "Enforce error"},
	}
)

func init() {
	xcode.CodeAdds(systemErrDepot)
	xcode.CodeAdds(businessErrDepot)
}
