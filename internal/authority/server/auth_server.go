package auth_server

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	xapp "github.com/coder2m/component"
	"github.com/coder2m/g-saber/xcast"
	"github.com/coder2m/g-saber/xlog"
	"gorm.io/gorm"

	xclient "github.com/coder2m/ndisk/internal/authority/client"
	_map "github.com/coder2m/ndisk/internal/authority/map"
	"github.com/coder2m/ndisk/internal/authority/model"
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
	_ = xclient.CasbinClient().LoadPolicy()
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

// 更新角色添加权限和菜单
func UpdateRolesMenuAndResources(ctx context.Context, req _map.UpdateRolesMenuAndResourcesReq) error {
	var menuList = make([]model.Menu, len(req.Menus))
	for i, menu := range req.Menus {
		menuList[i] = model.Menu{
			ID: xcast.ToUint(menu),
		}
	}

	var resourcesList = make([]model.Resources, len(req.Resources))
	for i, resources := range req.Resources {
		resourcesList[i] = model.Resources{
			ID: xcast.ToUint(resources),
		}
	}

	r := model.Roles{
		ID:        xcast.ToUint(req.ID),
		Menus:     menuList,
		Resources: resourcesList,
	}
	err := r.UpdateRolesMenuAndResources(ctx)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateRolesMenuAndResources", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("UpdateRolesMenuAndResources error")
	}
	_ = xclient.CasbinClient().LoadPolicy()
	return err
}

func GetMenuList(ctx context.Context, req _map.PageList) ([]_map.MenuResInfo, int64, error) {
	var (
		data    []model.Menu
		where   map[string][]interface{}
		resList []_map.MenuResInfo
	)

	if req.Keyword != "" {
		where = map[string][]interface{}{
			"id like ? or path like ? or name like ?": {
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
			},
		}
	}
	count, err := new(model.Menu).
		Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resList, 0, EmptyDataErr
		}
		xlog.Error("GetMenuList", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return resList, 0, errors.New("get menu list error")
	}
	resList = make([]_map.MenuResInfo, len(data))
	for i, datum := range data {
		resList[i] = _map.MenuResInfo{
			ID:          datum.ID,
			ParentId:    datum.ID,
			Path:        datum.Path,
			Name:        datum.Name,
			Description: datum.Description,
			IconClass:   datum.IconClass,
			CreatedAt:   xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return resList, count, err
}

func DeleteMenu(ctx context.Context, req _map.Ids) (int64, error) {
	var list = make([]string, len(req.List))
	for i, u := range req.List {
		list[i] = xcast.ToString(u)
	}
	count, err := new(model.Menu).Del(ctx, map[string][]interface{}{
		"id in (?)": {strings.Join(list, ",")},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, EmptyDataErr
		}
		xlog.Error("DeleteMenu", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("DeleteMenu error")
	}
	return count, err
}

func AddMenu(ctx context.Context, req _map.MenuReq) error {
	r := &model.Menu{
		ParentId:    req.ParentId,
		Path:        req.Path,
		Name:        req.Name,
		Description: req.Description,
		IconClass:   req.IconClass,
	}
	err := r.Add(ctx)
	if !errors.Is(err, nil) {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				return ExistDataErr
			}
		}
		xlog.Error("AddMenu", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("create menu error")
	}
	return err
}

func UpdateMenu(ctx context.Context, id uint, req _map.MenuReq) error {
	r := &model.Menu{
		ParentId:    req.ParentId,
		Path:        req.Path,
		Name:        req.Name,
		Description: req.Description,
		IconClass:   req.IconClass,
	}
	err := r.UpdatesWhere(ctx, map[string][]interface{}{
		"id=?": {id},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateMenu", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("update menu error")
	}
	return err
}

// api 资源
func GetResourcesList(ctx context.Context, req _map.PageList) ([]_map.ResourcesResInfo, int64, error) {
	var (
		data    []model.Resources
		where   map[string][]interface{}
		resList []_map.ResourcesResInfo
	)

	if req.Keyword != "" {
		where = map[string][]interface{}{
			"id like ? or name like ? or path like ?": {
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
				"%" + req.Keyword + "%",
			},
		}
	}
	count, err := new(model.Resources).
		Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete, false)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resList, 0, EmptyDataErr
		}
		xlog.Error("GetResourcesList", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return resList, 0, errors.New("get resources list error")
	}
	resList = make([]_map.ResourcesResInfo, len(data))
	for i, datum := range data {
		resList[i] = _map.ResourcesResInfo{
			ID:          datum.ID,
			Name:        datum.Name,
			Path:        datum.Path,
			Action:      datum.Action,
			Description: datum.Description,
			CreatedAt:   xcast.ToUint64(datum.CreatedAt.Unix()),
			UpdatedAt:   xcast.ToUint64(datum.UpdatedAt.Unix()),
			DeletedAt:   xcast.ToUint64(datum.DeletedAt.Time.Unix()),
		}
	}
	return resList, count, err
}

func DeleteResources(ctx context.Context, req _map.Ids) (int64, error) {
	count, err := new(model.Resources).Del(ctx, req.List)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, EmptyDataErr
		}
		xlog.Error("DeleteResources", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("delete resources error")
	}
	_ = xclient.CasbinClient().LoadPolicy()
	return count, err
}

func AddResources(ctx context.Context, req _map.ResourcesReq) error {
	r := &model.Resources{
		Name:        req.Name,
		Path:        req.Path,
		Action:      strings.ToUpper(req.Action),
		Description: req.Description,
	}
	err := r.Add(ctx)
	if !errors.Is(err, nil) {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				return ExistDataErr
			}
		}
		xlog.Error("AddResources", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("create resources error")
	}
	return err
}

func UpdateResources(ctx context.Context, id uint, req _map.ResourcesReq) error {
	r := &model.Resources{
		Name:        req.Name,
		Path:        req.Path,
		Action:      strings.ToUpper(req.Action),
		Description: req.Description,
	}
	err := r.UpdatesWhereById(ctx, id)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateResources", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("update resources error")
	}
	_ = xclient.CasbinClient().LoadPolicy()
	return err
}

func GetUsersRoles(ctx context.Context, ids []uint32) (map[uint32]string, error) {
	data, err := new(model.CasbinRule).GetUsersRoles(ctx, ids)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, EmptyDataErr
		}
		xlog.Error("GetUsersRoles", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, errors.New("get UsersRoles error")
	}
	var list = make(map[uint32]string)
	for _, datum := range data {
		list[datum.Uid] = datum.Roles
	}
	return list, err
}
