package auth_server

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/xlog"
	xclient "github.com/myxy99/ndisk/internal/authority/client"
	_map "github.com/myxy99/ndisk/internal/authority/map"
	"github.com/myxy99/ndisk/internal/authority/model"
	"gorm.io/gorm"
)

var (
	EmptyDataErr = errors.New("empty data")
	ExistDataErr = errors.New("data exist")
)

func GetRolesList(ctx context.Context, req _map.PageList) ([]_map.RolesResInfo, int64, error) {
	var (
		data    []model.Roles
		where   map[string][]interface{}
		resList []_map.RolesResInfo
	)

	if req.Keyword != "" {
		where = map[string][]interface{}{
			"id like ? or name like ?": {
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
			},
		}
	}
	count, err := new(model.Roles).
		Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete, false)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resList, 0, EmptyDataErr
		}
		xlog.Error("GetAllRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return resList, 0, errors.New("get role list error")
	}
	resList = make([]_map.RolesResInfo, len(data))
	for i, datum := range data {
		resList[i] = _map.RolesResInfo{
			ID:          datum.ID,
			Name:        datum.Name,
			Description: datum.Description,
			CreatedAt:   xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return resList, count, err
}

func DeleteRoles(ctx context.Context, req _map.Ids) (int64, error) {
	count, err := new(model.Roles).Del(ctx, req.List)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, EmptyDataErr
		}
		xlog.Error("DeleteRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("DeleteRoles error")
	}
	_ = xclient.CasbinClient().LoadPolicy()
	return count, err
}

func AddRoles(ctx context.Context, req _map.RolesReq) error {
	r := &model.Roles{
		Name:        req.Name,
		Description: req.Description,
	}
	err := r.Add(ctx)
	if !errors.Is(err, nil) {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				return ExistDataErr
			}
		}
		xlog.Error("AddRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("create roles error")
	}
	return err
}

func UpdateRoles(ctx context.Context, id uint, req _map.RolesReq) error {
	r := &model.Roles{
		Name:        req.Name,
		Description: req.Description,
	}
	err := r.UpdatesWhereById(ctx, id)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("update roles error")
	}
	return err
}

// 获取角色下的所有权限和菜单
func GetPermissionAndMenuByRoles(ctx context.Context, IdOrName string) (data _map.RolesResInfo, err error) {
	r := new(model.Roles)
	err = r.GetByWhere(ctx, map[string][]interface{}{
		"id = ? or name = ?": {
			IdOrName,
			IdOrName,
		},
	}, true)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return data, EmptyDataErr
		}
		xlog.Error("GetPermissionByRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return data, errors.New("GetPermissionByRoles error")
	}
	var resourcesList = make([]_map.ResourcesResInfo, len(r.Resources))
	for i, resource := range r.Resources {
		resourcesList[i] = _map.ResourcesResInfo{
			ID:          resource.ID,
			Name:        resource.Name,
			Path:        resource.Path,
			Action:      resource.Action,
			Description: resource.Description,
			CreatedAt:   xcast.ToUint64(resource.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(resource.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(resource.DeletedAt.Time.Unix()),
		}
	}
	var menuList = make([]_map.MenuResInfo, len(r.Menus))
	for i, resource := range r.Menus {
		menuList[i] = _map.MenuResInfo{
			ID:          resource.ID,
			ParentId:    resource.ParentId,
			Path:        resource.Path,
			Name:        resource.Name,
			Description: resource.Description,
			IconClass:   resource.IconClass,
			CreatedAt:   xcast.ToUint64(resource.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(resource.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(resource.DeletedAt.Time.Unix()),
		}
	}
	return _map.RolesResInfo{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Resources:   resourcesList,
		Menus:       menuList,
		CreatedAt:   xcast.ToUint64(r.CreatedAt.Unix()),
		UpdatedAt:   xcast.ToUint64(r.UpdatedAt.Unix()),
	}, err
}
