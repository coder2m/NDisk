/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 18:37
 **/
package xrpc

import (
	"github.com/coder2m/component/xcode"
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
	CreateManyAgencyErrCode
	DelManyAgencyErrCode
	ListAgencyErrCode
	UpdateAgencyErrCode
	GetAgencyByIdErrCode
	UpdateAgencyStatusErrCode
	RecoverDelAgencyErrCode
	ListAgencyByCreateUIdErrCode
	ListAgencyByJoinUIdErrCode
	ListUserByJoinAgencyErrCode
	UpdateStatusAgencyUserErrCode
	DelManyAgencyUserErrCode

	//Authority
	GetRolesForUserErrCode
	AddRolesForUserErrCode
	HasRoleForUserErrCode
	DeleteRoleForUserErrCode
	DeleteRolesForUserErrCode
	EnforceErrCode
	GetAllRolesErrCode
	DeleteRolesErrCode
	AddRolesErrCode
	UpdateRolesErrCode
	GetAllMenuErrCode
	DeleteMenuErrCode
	AddMenuErrCode
	UpdateMenuErrCode
	GetAllResourcesErrCode
	DeleteResourcesErrCode
	AddResourcesErrCode
	UpdateResourcesErrCode
	UpdateRolesMenuAndResourcesErrCode
	GetPermissionAndMenuByRolesErrCode
	GetUsersRolesErrCode
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
		{xcode.BusinessType, CreateManyAgencyErrCode, "CreateManyAgency error"},
		{xcode.BusinessType, DelManyAgencyErrCode, "DelManyAgency error"},
		{xcode.BusinessType, ListAgencyErrCode, "ListAgencyErr error"},
		{xcode.BusinessType, UpdateAgencyErrCode, "UpdateAgency error"},
		{xcode.BusinessType, GetAgencyByIdErrCode, "GetAgencyByIdErr error"},
		{xcode.BusinessType, UpdateAgencyStatusErrCode, "UpdateAgencyStatus error"},
		{xcode.BusinessType, RecoverDelAgencyErrCode, "RecoverDelAgency error"},
		{xcode.BusinessType, ListAgencyByCreateUIdErrCode, "ListAgencyByCreateUId error"},
		{xcode.BusinessType, ListAgencyByJoinUIdErrCode, "ListAgencyByJoinUId error"},
		{xcode.BusinessType, ListUserByJoinAgencyErrCode, "ListUserByJoinAgency error"},
		{xcode.BusinessType, UpdateStatusAgencyUserErrCode, "UpdateStatusAgency error"},
		{xcode.BusinessType, DelManyAgencyUserErrCode, "DelManyAgencyUser error"},

		{xcode.BusinessType, GetRolesForUserErrCode, "get roles for user error"},
		{xcode.BusinessType, AddRolesForUserErrCode, "AddRolesForUse error"},
		{xcode.BusinessType, HasRoleForUserErrCode, "HasRoleForUserErrCode error"},
		{xcode.BusinessType, DeleteRoleForUserErrCode, "DeleteRoleForUserErrCode error"},
		{xcode.BusinessType, DeleteRolesForUserErrCode, "DeleteRolesForUser error"},
		{xcode.BusinessType, EnforceErrCode, "Enforce error"},
		{xcode.BusinessType, GetAllRolesErrCode, "GetAllRoles error"},
		{xcode.BusinessType, DeleteRolesErrCode, "DeleteRoles error"},
		{xcode.BusinessType, AddRolesErrCode, "AddRoles error"},
		{xcode.BusinessType, UpdateRolesErrCode, "UpdateRoles error"},
		{xcode.BusinessType, GetAllMenuErrCode, "GetAllMenu error"},
		{xcode.BusinessType, DeleteMenuErrCode, "DeleteMenu error"},
		{xcode.BusinessType, AddMenuErrCode, "AddMenuErr error"},
		{xcode.BusinessType, UpdateMenuErrCode, "UpdateMenu error"},
		{xcode.BusinessType, GetAllResourcesErrCode, "GetAllResources error"},
		{xcode.BusinessType, DeleteResourcesErrCode, "DeleteResources error"},
		{xcode.BusinessType, AddResourcesErrCode, "AddResources error"},
		{xcode.BusinessType, UpdateResourcesErrCode, "UpdateResources error"},
		{xcode.BusinessType, UpdateRolesMenuAndResourcesErrCode, "UpdateRolesMenuAndResources error"},
		{xcode.BusinessType, GetPermissionAndMenuByRolesErrCode, "GetPermissionAndMenuByRoles error"},
		{xcode.BusinessType, GetUsersRolesErrCode, "GetUsersRoles error"},
	}
)

func init() {
	xcode.CodeAdds(systemErrDepot)
	xcode.CodeAdds(businessErrDepot)
}
