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
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/pkg/xcode"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/component/xlog"

	xclient "github.com/myxy99/ndisk/internal/authority/client"
	_map "github.com/myxy99/ndisk/internal/authority/map"
	auth_server "github.com/myxy99/ndisk/internal/authority/server"
	"github.com/myxy99/ndisk/pkg/pb/authority"
	xrpc "github.com/myxy99/ndisk/pkg/rpc"
)

type Server struct{}

func (s Server) GetAllRoles(ctx context.Context, request *AuthorityPb.PageRequest) (*AuthorityPb.RolesInfoListResponse, error) {
	var req = _map.DefaultPageRequest
	req.Page = request.Page
	req.PageSize = request.Limit
	req.Keyword = request.Keyword
	req.IsDelete = request.IsDelete
	data, count, err := auth_server.GetRolesList(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.GetAllRolesErrCode)
	}
	var list = make([]*AuthorityPb.RolesInfo, len(data))
	for i, datum := range data {
		list[i] = &AuthorityPb.RolesInfo{
			Id:          xcast.ToUint32(datum.ID),
			Name:        datum.Name,
			Description: datum.Description,
			CreatedAt:   datum.CreatedAt,
			UpdatedAt:   datum.UpdatedAt,
			DeletedAt:   datum.DeletedAt,
		}
	}
	return &AuthorityPb.RolesInfoListResponse{
		List:  list,
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) DeleteRoles(ctx context.Context, ids *AuthorityPb.Ids) (*AuthorityPb.ChangeNumRes, error) {
	count, err := auth_server.DeleteRoles(ctx, _map.Ids{List: ids.To})
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.DeleteRolesErrCode)
	}
	return &AuthorityPb.ChangeNumRes{
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) AddRoles(ctx context.Context, info *AuthorityPb.RolesInfo) (*AuthorityPb.Empty, error) {
	req := _map.RolesReq{
		Name:        info.Name,
		Description: info.Description,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("AddRoles data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.AddRoles(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.AddRolesErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) UpdateRoles(ctx context.Context, info *AuthorityPb.RolesInfo) (*AuthorityPb.Empty, error) {
	req := _map.RolesReq{
		Name:        info.Name,
		Description: info.Description,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateRoles data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.UpdateRoles(ctx, xcast.ToUint(info.Id), req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.UpdateRolesErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) GetAllMenu(ctx context.Context, request *AuthorityPb.PageRequest) (*AuthorityPb.MenuInfoListResponse, error) {
	var req = _map.DefaultPageRequest
	req.Page = request.Page
	req.PageSize = request.Limit
	req.Keyword = request.Keyword
	req.IsDelete = request.IsDelete
	data, count, err := auth_server.GetMenuList(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.GetAllMenuErrCode)
	}
	var list = make([]*AuthorityPb.MenuInfo, len(data))
	for i, datum := range data {
		list[i] = &AuthorityPb.MenuInfo{
			Id:          xcast.ToUint32(datum.ID),
			ParentId:    xcast.ToUint32(datum.ParentId),
			Path:        datum.Path,
			Name:        datum.Name,
			Description: datum.Description,
			IconClass:   datum.IconClass,
			CreatedAt:   datum.CreatedAt,
			UpdatedAt:   datum.UpdatedAt,
			DeletedAt:   datum.DeletedAt,
		}
	}
	return &AuthorityPb.MenuInfoListResponse{
		List:  list,
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) DeleteMenu(ctx context.Context, ids *AuthorityPb.Ids) (*AuthorityPb.ChangeNumRes, error) {
	count, err := auth_server.DeleteMenu(ctx, _map.Ids{List: ids.To})
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.DeleteMenuErrCode)
	}
	return &AuthorityPb.ChangeNumRes{
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) AddMenu(ctx context.Context, info *AuthorityPb.MenuInfo) (*AuthorityPb.Empty, error) {
	req := _map.MenuReq{
		ParentId:    xcast.ToUint(info.ParentId),
		Path:        info.Path,
		Name:        info.Name,
		Description: info.Description,
		IconClass:   info.IconClass,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("AddMenu data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.AddMenu(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.AddMenuErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) UpdateMenu(ctx context.Context, info *AuthorityPb.MenuInfo) (*AuthorityPb.Empty, error) {
	req := _map.MenuReq{
		ParentId:    xcast.ToUint(info.ParentId),
		Path:        info.Path,
		Name:        info.Name,
		Description: info.Description,
		IconClass:   info.IconClass,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateMenu data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.UpdateMenu(ctx, xcast.ToUint(info.Id), req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.UpdateMenuErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) GetAllResources(ctx context.Context, request *AuthorityPb.PageRequest) (*AuthorityPb.ResourcesInfoListResponse, error) {
	var req = _map.DefaultPageRequest
	req.Page = request.Page
	req.PageSize = request.Limit
	req.Keyword = request.Keyword
	req.IsDelete = request.IsDelete
	data, count, err := auth_server.GetResourcesList(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.GetAllResourcesErrCode)
	}
	var list = make([]*AuthorityPb.ResourcesInfo, len(data))
	for i, datum := range data {
		list[i] = &AuthorityPb.ResourcesInfo{
			Id:          xcast.ToUint32(datum.ID),
			Name:        datum.Name,
			Path:        datum.Path,
			Action:      datum.Action,
			Description: datum.Description,
			CreatedAt:   datum.CreatedAt,
			UpdatedAt:   datum.UpdatedAt,
			DeletedAt:   datum.DeletedAt,
		}
	}
	return &AuthorityPb.ResourcesInfoListResponse{
		List:  list,
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) DeleteResources(ctx context.Context, ids *AuthorityPb.Ids) (*AuthorityPb.ChangeNumRes, error) {
	count, err := auth_server.DeleteResources(ctx, _map.Ids{List: ids.To})
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.DeleteResourcesErrCode)
	}
	return &AuthorityPb.ChangeNumRes{
		Count: xcast.ToUint32(count),
	}, err
}

func (s Server) AddResources(ctx context.Context, info *AuthorityPb.ResourcesInfo) (*AuthorityPb.Empty, error) {
	req := _map.ResourcesReq{
		Name:        info.Name,
		Path:        info.Path,
		Action:      info.Action,
		Description: info.Description,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("AddResources data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.AddResources(ctx, req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.AddResourcesErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) UpdateResources(ctx context.Context, info *AuthorityPb.ResourcesInfo) (*AuthorityPb.Empty, error) {
	req := _map.ResourcesReq{
		Name:        info.Name,
		Path:        info.Path,
		Action:      info.Action,
		Description: info.Description,
	}
	err := xvalidator.Struct(req)
	if !errors.Is(err, nil) {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("UpdateResources data validation error : %s", xvalidator.GetMsg(err).Error())
	}
	err = auth_server.UpdateResources(ctx, xcast.ToUint(info.Id), req)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.UpdateResourcesErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) UpdateRolesMenuAndResources(ctx context.Context, req *AuthorityPb.UpdateRolesMenuAndResourcesReq) (*AuthorityPb.Empty, error) {
	err := auth_server.UpdateRolesMenuAndResources(ctx, _map.UpdateRolesMenuAndResourcesReq{
		ID:        req.Id,
		Menus:     req.Menus,
		Resources: req.Resources,
	})
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.UpdateRolesMenuAndResourcesErrCode)
	}
	return new(AuthorityPb.Empty), err
}

func (s Server) GetPermissionAndMenuByRoles(ctx context.Context, target *AuthorityPb.Target) (*AuthorityPb.RolesInfo, error) {
	if target.To == "" {
		return nil, xcode.BusinessCode(xrpc.ValidationErrCode).SetMsgf("GetPermissionAndMenuByRoles data validation error : %s", "is not nil")
	}
	data, err := auth_server.GetPermissionAndMenuByRoles(ctx, target.To)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.GetPermissionAndMenuByRolesErrCode)
	}
	var menusList = make([]*AuthorityPb.MenuInfo, len(data.Menus))
	for i, menu := range data.Menus {
		menusList[i] = &AuthorityPb.MenuInfo{
			Id:          xcast.ToUint32(menu.ID),
			ParentId:    xcast.ToUint32(menu.ParentId),
			Path:        menu.Path,
			Name:        menu.Name,
			Description: menu.Description,
			IconClass:   menu.IconClass,
			CreatedAt:   menu.CreatedAt,
			UpdatedAt:   menu.UpdatedAt,
			DeletedAt:   menu.DeletedAt,
		}
	}

	var resourcesList = make([]*AuthorityPb.ResourcesInfo, len(data.Resources))
	for i, resources := range data.Resources {
		resourcesList[i] = &AuthorityPb.ResourcesInfo{
			Id:          xcast.ToUint32(resources.ID),
			Name:        resources.Name,
			Path:        resources.Path,
			Action:      resources.Action,
			Description: resources.Description,
			CreatedAt:   resources.CreatedAt,
			UpdatedAt:   resources.UpdatedAt,
			DeletedAt:   resources.DeletedAt,
		}
	}

	return &AuthorityPb.RolesInfo{
		Id:          xcast.ToUint32(data.ID),
		Name:        data.Name,
		Description: data.Description,
		Menus:       menusList,
		Resources:   resourcesList,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		DeletedAt:   data.DeletedAt,
	}, err
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

func (s Server) GetUsersRoles(ctx context.Context, req *AuthorityPb.Ids) (res *AuthorityPb.UsersRole, err error) {
	data, err := auth_server.GetUsersRoles(ctx, req.To)
	if !errors.Is(err, nil) {
		if err == auth_server.EmptyDataErr {
			return nil, xcode.BusinessCode(xrpc.EmptyData)
		}
		return nil, xcode.BusinessCode(xrpc.GetUsersRolesErrCode)
	}
	return &AuthorityPb.UsersRole{
		Data: data,
	}, err
}
