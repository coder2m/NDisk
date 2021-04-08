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
