package directory

import (
	"github.com/coder2z/g-saber/xcast"
	"github.com/coder2z/g-saber/xvalidator"
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
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		if data, err := directory_server.List(ctx.Request.Context(), info.Uid, req, pageReq); err != nil {
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
func Add(ctx *gin.Context) {
	var req _map.DirectoryPost
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}

	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}

	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		req.Uid = xcast.ToUint(info.Uid)
		if data, err := directory_server.Add(ctx.Request.Context(), req); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, data)
		}
	} else {
		R.HandleForbidden(ctx)
	}
	return
}

//删除文件夹 以及文件
func Del(ctx *gin.Context) {
	var req _map.Id
	if err := ctx.BindUri(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}

	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}

	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		if err := directory_server.Del(ctx.Request.Context(), xcast.ToUint(info.Uid), req.Id); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, nil)
		}
	} else {
		R.HandleForbidden(ctx)
	}

	return
}

func Update(ctx *gin.Context) {
	var (
		reqId  _map.Id
		reqDir _map.DirectoryUpdate
	)
	if err := ctx.BindUri(&reqId); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}

	if err := xvalidator.Struct(reqId); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}

	if err := ctx.ShouldBindJSON(&reqDir); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}

	if err := xvalidator.Struct(reqDir); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}

	reqDir.Id = reqId.Id

	if i, ok := ctx.Get("user"); ok {
		info := i.(_map.UserInfo)
		reqDir.Uid = xcast.ToUint(info.Uid)
		if err := directory_server.Update(ctx.Request.Context(), reqDir); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, nil)
		}
	} else {
		R.HandleForbidden(ctx)
	}

	return
}
