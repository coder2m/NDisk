package admin_server

import (
	"context"
	"errors"
	"github.com/myxy99/component/pkg/xcast"
	xclient "github.com/myxy99/ndisk/internal/getway/client"
	xerror "github.com/myxy99/ndisk/internal/getway/error"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	AuthorityPb "github.com/myxy99/ndisk/pkg/pb/authority"
	NUserPb "github.com/myxy99/ndisk/pkg/pb/nuser"
)

func CreateUser(ctx context.Context, user _map.CreateUser) (data _map.Batch, errs *xerror.Err) {
	var list = make([]*NUserPb.UserInfo, len(user.Data))
	for i, datum := range user.Data {
		list[i] = &NUserPb.UserInfo{
			Name:     datum.Name,
			Alias:    datum.Alias,
			Tel:      datum.Tel,
			Email:    datum.Email,
			Password: datum.Password,
		}
	}
	rep, err := xclient.NUserServer.CreateUsers(ctx, &NUserPb.UserList{
		List: list,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UpdateUser(ctx context.Context, user _map.UpdateUser) *xerror.Err {
	_, err := xclient.NUserServer.UpdateUser(ctx, &NUserPb.UserInfo{
		Uid:      user.Uid,
		Name:     user.Name,
		Alias:    user.Alias,
		Tel:      user.Tel,
		Email:    user.Email,
		Password: user.Password,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func DeleteUser(ctx context.Context, user _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.NUserServer.DelUsers(ctx, &NUserPb.UidList{
		Uid: user.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UserList(ctx context.Context, req _map.PageList) (data _map.UserList, errs *xerror.Err) {
	rep, err := xclient.NUserServer.GetUserList(ctx, &NUserPb.PageRequest{
		Keyword:  req.Keyword,
		Page:     xcast.ToUint32(req.Page),
		Limit:    xcast.ToUint32(req.PageSize),
		IsDelete: req.IsDelete,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	list := make([]_map.UserInfo, len(rep.List))
	for i, info := range rep.List {
		list[i] = _map.UserInfo{
			Uid:         info.Uid,
			Name:        info.Name,
			Alias:       info.Alias,
			Tel:         info.Tel,
			Email:       info.Email,
			Status:      info.Status,
			EmailStatus: info.EmailStatus,
			CreatedAt:   info.CreatedAt,
			UpdatedAt:   info.UpdatedAt,
			DeletedAt:   info.DeletedAt,
		}
	}
	return _map.UserList{
		Count: rep.Count,
		Data:  list,
	}, nil
}

func UpdateStatusUser(ctx context.Context, req _map.UpdateStatus) *xerror.Err {
	_, err := xclient.NUserServer.UpdateUserStatus(ctx, &NUserPb.UserInfo{
		Uid:    req.Uid,
		Status: req.Status,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func RestoreDeleteUser(ctx context.Context, req _map.UidList) (data _map.Batch, errs *xerror.Err) {
	rep, err := xclient.NUserServer.RecoverDelUsers(ctx, &NUserPb.UidList{
		Uid: req.List,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.Batch{Count: rep.Count}, nil
}

func UserById(ctx context.Context, req _map.Uid) (data _map.UserInfo, errs *xerror.Err) {
	rep, err := xclient.NUserServer.GetUserById(ctx, &NUserPb.UserInfo{
		Uid: req.Uid,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return _map.UserInfo{
		Uid:         rep.Uid,
		Name:        rep.Name,
		Alias:       rep.Alias,
		Tel:         rep.Tel,
		Email:       rep.Email,
		Status:      rep.Status,
		EmailStatus: rep.EmailStatus,
		CreatedAt:   rep.CreatedAt,
		UpdatedAt:   rep.UpdatedAt,
		DeletedAt:   rep.DeletedAt,
	}, nil
}

// 获取角色下的权限
func CompetenceList(ctx context.Context, role string) (data _map.CompetenceList, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetPermissionsForUser(ctx, &AuthorityPb.Single{
		To: role,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	var list = make([]_map.Competence, len(rep.List))
	for i, array := range rep.List {
		list[i] = _map.Competence{
			Objective: array.Data[1],
			Action:    array.Data[2],
		}
	}
	return _map.CompetenceList{
		Data: list,
	}, nil
}

// 给角色添加权限
func RoleAddCompetence(ctx context.Context, role string, data _map.CompetenceReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddPermissionForUser(ctx, &AuthorityPb.Resources{
		Role:   role,
		Obj:    data.Objective,
		Action: data.Action,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

// 给角色删除权限
func RoleDelCompetence(ctx context.Context, role string, data _map.CompetenceReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.DeletePermissionForUser(ctx, &AuthorityPb.Batch{
		To:      role,
		Operate: []string{data.Objective, data.Action},
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func GetAllRoles(ctx context.Context) (data []string, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetAllRoles(ctx, &AuthorityPb.Empty{})
	if !errors.Is(err, nil) {
		return nil, xerror.NewErrRPC(err)
	}
	return rep.Data, nil
}

func CompetenceByRole(ctx context.Context, role string) (data _map.CompetenceList, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetPermissionsForUser(ctx, &AuthorityPb.Single{
		To: role,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	var list = make([]_map.Competence, len(rep.List))
	for i, array := range rep.List {
		list[i] = _map.Competence{
			Objective: array.Data[1],
			Action:    array.Data[2],
		}
	}
	return _map.CompetenceList{
		Data: list,
	}, nil
}

func UserByRole(ctx context.Context, role string) (data _map.UserList, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetUsersForRole(ctx, &AuthorityPb.Single{
		To: role,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	var uidList = make([]uint32, len(rep.Data))
	for i, datum := range rep.Data {
		uidList[i] = xcast.ToUint32(datum)
	}
	reqUser, err := xclient.NUserServer.GetUserListByUid(ctx, &NUserPb.UidList{
		Uid: uidList,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	list := make([]_map.UserInfo, len(reqUser.List))
	for i, info := range reqUser.List {
		list[i] = _map.UserInfo{
			Uid:       info.Uid,
			Name:      info.Name,
			Alias:     info.Alias,
			Tel:       info.Tel,
			Email:     info.Email,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
			DeletedAt: info.DeletedAt,
		}
	}
	return _map.UserList{
		Count: reqUser.Count,
		Data:  list,
	}, nil
}

func RoleByUser(ctx context.Context, uid string) (data []string, errs *xerror.Err) {
	rep, err := xclient.AuthorityServer.GetRolesForUser(ctx, &AuthorityPb.Target{
		To: uid,
	})
	if !errors.Is(err, nil) {
		return data, xerror.NewErrRPC(err)
	}
	return rep.Data, nil
}

func UserAddRoles(ctx context.Context, req _map.UserRolesReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.AddRolesForUser(ctx, &AuthorityPb.Batch{
		To:      xcast.ToString(req.Uid),
		Operate: req.Role,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}

func UserDelRoles(ctx context.Context, req _map.UserRoleReq) (errs *xerror.Err) {
	_, err := xclient.AuthorityServer.DeleteRoleForUser(ctx, &AuthorityPb.Single{
		To:      xcast.ToString(req.Uid),
		Operate: req.Role,
	})
	if !errors.Is(err, nil) {
		return xerror.NewErrRPC(err)
	}
	return nil
}
