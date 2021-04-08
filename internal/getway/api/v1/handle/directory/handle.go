package directory

import (
	"github.com/coder2z/g-saber/xcast"
	_map "github.com/coder2z/ndisk/internal/getway/map"
	"github.com/coder2z/ndisk/internal/getway/server/directory_server"
	R "github.com/coder2z/ndisk/pkg/response"
	"github.com/gin-gonic/gin"
)

//文件夹列表 文件夹+文件
func List(ctx *gin.Context) {
	var req _map.Id
	if err := ctx.BindUri(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}

	var pageReq = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&pageReq); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		if data, err := directory_server.List(ctx, info.Uid, req, pageReq); err != nil {
			R.Error(ctx, err)
		} else {
			R.Page(ctx, xcast.ToInt64(data.Count), pageReq.Page, pageReq.PageSize, data.Data)
		}
	} else {
		R.HandleForbidden(ctx)
	}

	return
}

//添加文件夹 以及文件

//删除文件夹 以及文件
