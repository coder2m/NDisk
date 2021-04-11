package directory_server

import (
	"context"
	"github.com/coder2z/g-saber/xcast"
	_map "github.com/coder2z/ndisk/internal/getway/map"
	"github.com/coder2z/ndisk/internal/getway/model"
)

func List(ctx context.Context, uid uint64, req _map.Id, page _map.PageList) (_map.DirList, error) {
	var (
		data  []model.Directory
		where = make(map[string][]interface{})
	)
	where["parent_id = ?"] = []interface{}{
		req.Id,
	}
	where["uid = ?"] = []interface{}{
		uid,
	}
	if page.Keyword != "" {
		where["name like ?"] = []interface{}{
			"%" + page.Keyword + "%",
		}
	}
	count, err := new(model.Directory).Get(ctx, page.PageSize*(page.Page-1), page.PageSize, &data, where, page.IsDelete)
	if err != nil {
		return _map.DirList{}, err
	}
	res := make([]_map.DirectoryInfo, len(data))
	for i, datum := range data {
		res[i] = _map.DirectoryInfo{
			Id:        datum.ID,
			Uid:       datum.Uid,
			FileId:    datum.FileId,
			IsDir:     datum.IsDir,
			Name:      datum.Name,
			ParentID:  datum.ParentID,
			CreatedAt: datum.CreatedAt.Unix(),
			UpdatedAt: datum.UpdatedAt.Unix(),
			DeletedAt: datum.DeletedAt.Time.Unix(),
		}
	}
	return _map.DirList{
		Count: xcast.ToUint32(count),
		Data:  res,
	}, err
}

func Add(ctx context.Context, post _map.DirectoryPost) (_map.DirectoryInfo, error) {
	dir := &model.Directory{
		Uid:      post.Uid,
		FileId:   post.FileId,
		IsDir:    post.IsDir,
		Name:     post.Name,
		ParentID: post.ParentID,
		Type:     post.Type,
	}
	err := dir.Add(ctx)
	if err != nil {
		return _map.DirectoryInfo{}, err
	}
	return _map.DirectoryInfo{
		Id:        dir.ID,
		Uid:       dir.Uid,
		FileId:    dir.FileId,
		IsDir:     dir.IsDir,
		Name:      dir.Name,
		ParentID:  dir.ParentID,
		Type:      dir.Type,
		CreatedAt: dir.CreatedAt.Unix(),
		UpdatedAt: dir.UpdatedAt.Unix(),
		DeletedAt: dir.DeletedAt.Time.Unix(),
	}, err
}

func Del(ctx context.Context, uid, id uint) error {
	var (
		where map[string][]interface{}
	)
	where = map[string][]interface{}{
		"id=? and uid =?": {
			id, uid,
		},
	}
	_, err := new(model.Directory).Del(ctx, where)
	return err
}

func Update(ctx context.Context, req _map.DirectoryUpdate) error {
	var (
		where map[string][]interface{}
	)
	where = map[string][]interface{}{
		"id=? and uid =?": {
			req.Id, req.Uid,
		},
	}
	dir := &model.Directory{
		Name:     req.Name,
		ParentID: req.ParentID,
	}
	err := dir.UpdatesWhere(ctx, where)
	return err
}
