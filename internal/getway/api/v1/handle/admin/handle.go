package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/myxy99/component/pkg/xcast"
	"github.com/myxy99/component/pkg/xvalidator"
	"github.com/myxy99/ndisk/internal/getway/error/httpError"
	_map "github.com/myxy99/ndisk/internal/getway/map"
	"github.com/myxy99/ndisk/internal/getway/server/admin_server"
	R "github.com/myxy99/ndisk/pkg/response"
)

// 管理员创建用户
func CreateUser(ctx *gin.Context) {
	var req _map.CreateUser
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.CreateUser(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员用户信息修改
func UpdateUser(ctx *gin.Context) {
	var req _map.UpdateUser
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UpdateUser(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 管理员删除用户
func DeleteUser(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DeleteUser(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员用户列表
func UserList(ctx *gin.Context) {
	var req = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.UserList(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

// 管理员修改用户状态
func UpdateStatusUser(ctx *gin.Context) {
	var req _map.UpdateStatus
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UpdateStatusUser(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 管理员恢复已经删除用户
func RestoreDeleteUser(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.RestoreDeleteUser(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员获取用户信息
func UserById(ctx *gin.Context) {
	var req _map.Uid
	if err := ctx.BindUri(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.UserById(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员获取全部角色
func RoleList(ctx *gin.Context) {
	var req = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.GetAllRoles(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

// 获取角色下的所有用户
func UserByRole(ctx *gin.Context) {
	var req _map.RoleReq
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.UserByRole(ctx, req.Content); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员获取用户的角色
func RoleByUser(ctx *gin.Context) {
	var req _map.Uid
	if err := ctx.BindUri(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.RoleByUser(ctx, xcast.ToString(req.Uid)); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 管理员给用户添加角色
func UserAddRoles(ctx *gin.Context) {
	var req _map.UserRolesReq
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UserAddRoles(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 管理员删除用户角色
func DeleteUserRole(ctx *gin.Context) {
	var req _map.UserRoleReq
	if err := ctx.ShouldBind(&req); err != nil {
		httpError.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		httpError.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UserDelRoles(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}
