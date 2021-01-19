package agency_server

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/xlog"
	_map "github.com/myxy99/ndisk/internal/nuser/map"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"gorm.io/gorm"
	"strings"
)

var (
	EmptyDataErr = errors.New("empty data")
	ExistDataErr = errors.New("data exist")
)

//创建机构
func CreateManyAgency(ctx context.Context, req _map.CreateManyAgencyReq) (count int64, err error) {
	var list = make([]model.Agency, len(req.Agency))
	for i, agencyReq := range req.Agency {
		list[i] = model.Agency{
			ParentId:  agencyReq.ParentId,
			Name:      agencyReq.Name,
			Remark:    agencyReq.Remark,
			CreateUId: xcast.ToUint(req.Uid),
		}
	}
	count, err = new(model.Agency).Adds(ctx, &list)
	if !errors.Is(err, nil) {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Number == 1062 {
				return 0, ExistDataErr
			}
		}
		xlog.Error("CreateManyAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("create agency error")
	}
	return
}

//	批量删除机构
func DelManyAgency(ctx context.Context, req _map.Ids) (count int64, err error) {
	var list = make([]string, len(req.List))
	for i, u := range req.List {
		list[i] = xcast.ToString(u)
	}
	count, err = new(model.Agency).Del(ctx, map[string][]interface{}{
		"id in (?)": {strings.Join(list, ",")},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, EmptyDataErr
		}
		xlog.Error("DelManyAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("DelManyAgency error")
	}
	return
}

//	机构列表
func ListAgency(ctx context.Context, parentId uint, req _map.PageList) ([]_map.AgencyInf, int64, error) {
	var (
		data    []model.Agency
		where   map[string][]interface{}
		repList []_map.AgencyInf
	)
	where = map[string][]interface{}{
		"parent_id=?": {parentId},
	}
	if req.Keyword != "" {
		where["id like ? or name like ? or create_uid like ?"] = []interface{}{
			"%" + req.Keyword + "%",
			"%" + req.Keyword + "%",
			"%" + req.Keyword + "%",
		}
	}
	total, err := new(model.Agency).Get(ctx, xcast.ToInt(req.PageSize*(req.Page-1)), xcast.ToInt(req.PageSize), &data, where, req.IsDelete, true)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repList, 0, EmptyDataErr
		}
		xlog.Error("ListAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return repList, 0, errors.New("ListAgency error")
	}
	repList = make([]_map.AgencyInf, len(data))
	for i, datum := range data {
		repList[i] = _map.AgencyInf{
			ID:       datum.ID,
			ParentId: datum.ParentId,
			Name:     datum.Name,
			Remark:   datum.Remark,
			Status:   datum.Status,
			CreateUser: _map.UserInfo{
				Uid:   datum.CreateUser.ID,
				Name:  datum.CreateUser.Name,
				Alias: datum.CreateUser.Alias,
				Tel:   datum.CreateUser.Tel,
				Email: datum.CreateUser.Email,
			},
			CreatedAt: datum.CreatedAt.Unix(),
			UpdatedAt: datum.UpdatedAt.Unix(),
			DeletedAt: datum.DeletedAt.Time.Unix(),
		}
	}
	return repList, total, err
}

//	修改机构信息
func UpdateAgency(ctx context.Context, req _map.UpdateAgency) (err error) {
	a := &model.Agency{
		ID:       req.ID,
		ParentId: req.ParentId,
		Name:     req.Name,
		Remark:   req.Remark,
	}
	err = a.UpdatesWhere(ctx, map[string][]interface{}{
		"id=?": {req.ID},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("UpdateAgency error")
	}
	return err
}

func UpdateStatusAgency(ctx context.Context, id uint, status uint32) (err error) {
	a := &model.Agency{
		ID: id,
	}
	err = a.UpdateStatus(ctx, status)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return EmptyDataErr
		}
		xlog.Error("UpdateStatusAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return errors.New("UpdateStatusAgency error")
	}
	return err
}

//	恢复删除后的机构
func RegainDelAgency(ctx context.Context, req _map.Ids) (count int64, err error) {
	var list = make([]string, len(req.List))
	for i, u := range req.List {
		list[i] = xcast.ToString(u)
	}
	count, err = new(model.Agency).DelRes(ctx, map[string][]interface{}{
		"id in (?)": {strings.Join(list, ",")},
	})
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, EmptyDataErr
		}
		xlog.Error("RegainDelAgency", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return 0, errors.New("RegainDelAgency error")
	}
	return count, err
}

// 获取指定用户创建的所有机构
func ListAgencyByCreateUId(ctx context.Context, req _map.Id, status uint) ([]_map.AgencyInf, error) {
	var data []model.Agency
	err := new(model.Agency).GetAll(ctx, &data, map[string][]interface{}{
		"id = ?":     {req.Id},
		"status = ?": {status},
	}, false)
	if !errors.Is(err, nil) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, EmptyDataErr
		}
		xlog.Error("ListAgencyByCreateUId", xlog.FieldErr(err), xlog.FieldName(xapp.Name()), xlog.FieldType("mysql"))
		return nil, errors.New("ListAgencyByCreateUId error")
	}
	var list = make([]_map.AgencyInf, len(data))
	for i, datum := range data {
		list[i] = _map.AgencyInf{
			ID:        datum.ID,
			ParentId:  datum.ParentId,
			Name:      datum.Name,
			Remark:    datum.Remark,
			Status:    datum.Status,
			CreatedAt: datum.CreatedAt.Unix(),
			UpdatedAt: datum.UpdatedAt.Unix(),
			DeletedAt: datum.DeletedAt.Time.Unix(),
		}
	}
	return list, err
}

//	获取用户加入的所有机构
func ListAgencyByJoinUId(ctx context.Context, req _map.Id) ([]_map.AgencyInf, error) {
	panic("TODO")
}

// 获取机构下的所有用户
func ListUserByJoinAgency(ctx context.Context, req _map.Id) ([]_map.UserInfo, error) {
	panic("TODO")
}
