/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/15 13:50
 **/
package rpc

import (
	"context"
	"errors"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xlog"
	xclient "github.com/myxy99/ndisk/internal/authority/client"
	_map "github.com/myxy99/ndisk/internal/authority/map"
	"github.com/myxy99/ndisk/pkg/pb/authority"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

type Server struct{}

func (s Server) GetAllRoles(ctx context.Context, empty *AuthorityPb.Empty) (*AuthorityPb.Array, error) {
	return &AuthorityPb.Array{
		Data: xclient.CasbinClient().GetAllRoles(),
	}, nil
}

func (s Server) DeleteRole(ctx context.Context, target *AuthorityPb.Target) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Target{
		To: target.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeleteRole data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeleteRole(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("delete role", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeleteRoleErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) GetRolesForUser(ctx context.Context, target *AuthorityPb.Target) (rep *AuthorityPb.Array, err error) {
	var req = _map.Target{
		To: target.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetRolesForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	data, err := xclient.CasbinClient().GetRolesForUser(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("get roles for user", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.GetRolesForUserErrCode)
	}
	return &AuthorityPb.Array{
		Data: data,
	}, nil
}

func (s Server) AddRolesForUser(ctx context.Context, batch *AuthorityPb.Batch) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Batch{
		To:      batch.To,
		Operate: batch.Operate,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("AddRolesForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().AddRolesForUser(req.To, req.Operate)
	if !errors.Is(err, nil) {
		xlog.Error("add roles for user", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.AddRolesForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) HasRoleForUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Determine, err error) {
	var req = _map.Single{
		To:      single.To,
		Operate: single.Operate,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("HasRoleForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	ok, err := xclient.CasbinClient().HasRoleForUser(req.To, req.Operate)
	if !errors.Is(err, nil) {
		xlog.Error("HasRoleForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.HasRoleForUserErrCode)
	}
	return &AuthorityPb.Determine{
		Ok: ok,
	}, nil
}

func (s Server) DeleteRoleForUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Single{
		To:      single.To,
		Operate: single.Operate,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeleteRoleForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeleteRoleForUser(req.To, req.Operate)
	if !errors.Is(err, nil) {
		xlog.Error("DeleteRoleForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeleteRoleForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) DeleteUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Target{
		To: single.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeleteUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeleteUser(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("DeleteUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeleteRoleForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) DeleteRolesForUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Target{
		To: single.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeleteRolesForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeleteRolesForUser(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("DeleteRolesForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeleteRolesForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) GetUsersForRole(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Array, err error) {
	var req = _map.Target{
		To: single.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUsersForRole data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	data, err := xclient.CasbinClient().GetUsersForRole(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("DeleteRolesForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeleteRolesForUserErrCode)
	}
	return &AuthorityPb.Array{
		Data: data,
	}, nil
}

func (s Server) AddPermissionForUser(ctx context.Context, resources *AuthorityPb.Resources) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Resources{
		Role:   resources.Role,
		Obj:    resources.Obj,
		Action: resources.Action,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("AddPermissionForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().AddPermissionForUser(req.Role, req.Obj, req.Action)
	if !errors.Is(err, nil) {
		xlog.Error("AddPermissionForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.AddPermissionForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) GetPermissionsForUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Arrays, err error) {
	var req = _map.Target{
		To: single.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetUsersForRole data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	data := xclient.CasbinClient().GetPermissionsForUser(req.To)
	var list []*AuthorityPb.Array
	for _, datum := range data {
		list = append(list, &AuthorityPb.Array{
			Data: datum,
		})
	}
	return &AuthorityPb.Arrays{
		List: list,
	}, nil
}

func (s Server) DeletePermissionForUser(ctx context.Context, single *AuthorityPb.Batch) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Batch{
		To:      single.To,
		Operate: single.Operate,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeletePermissionForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeletePermissionForUser(req.To, req.Operate...)
	if !errors.Is(err, nil) {
		xlog.Error("DeletePermissionForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeletePermissionForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) DeletePermissionsForUser(ctx context.Context, single *AuthorityPb.Single) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Target{
		To: single.To,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeletePermissionsForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeletePermissionsForUser(req.To)
	if !errors.Is(err, nil) {
		xlog.Error("DeletePermissionsForUser", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeletePermissionsForUserErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) DeletePermission(ctx context.Context, array *AuthorityPb.Array) (rep *AuthorityPb.Empty, err error) {
	var req = _map.Array{
		Data: array.Data,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("DeletePermission data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	_, err = xclient.CasbinClient().DeletePermission(req.Data...)
	if !errors.Is(err, nil) {
		xlog.Error("DeletePermission", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.DeletePermissionErrCode)
	}
	return new(AuthorityPb.Empty), nil
}

func (s Server) HasPermissionForUser(ctx context.Context, batch *AuthorityPb.Batch) (rep *AuthorityPb.Determine, err error) {
	var req = _map.Batch{
		To:      batch.To,
		Operate: batch.Operate,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("HasPermissionForUser data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	ok := xclient.CasbinClient().HasPermissionForUser(req.To, req.Operate...)
	return &AuthorityPb.Determine{
		Ok: ok,
	}, nil
}

func (s Server) Enforce(ctx context.Context, resources *AuthorityPb.Resources) (rep *AuthorityPb.Determine, err error) {
	var req = _map.Resources{
		Role:   resources.Role,
		Obj:    resources.Obj,
		Action: resources.Action,
	}
	err = xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return rep, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("Enforce data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	ok, err := xclient.CasbinClient().Enforce(req.Role, req.Obj, req.Action)
	if !errors.Is(err, nil) {
		xlog.Error("Enforce", xlog.FieldErr(err), xlog.FieldName(xapp.Name()))
		return rep, xcode.BusinessCode(xrpc.EnforceErrCode)
	}
	return &AuthorityPb.Determine{
		Ok: ok,
	}, nil
}
