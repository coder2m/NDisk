package admin

import (
	"github.com/coder2z/g-saber/xcast"
	"github.com/coder2z/g-saber/xvalidator"
	_map "github.com/coder2z/ndisk/internal/getway/map"
	"github.com/coder2z/ndisk/internal/getway/server/admin_server"
	"github.com/coder2z/ndisk/internal/getway/server/auth_server"
	R "github.com/coder2z/ndisk/pkg/response"
	"github.com/gin-gonic/gin"
)

// 管理员创建用户
func CreateUser(ctx *gin.Context) {
	var req _map.CreateUser
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.CreateUser(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UpdateUser(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DeleteUser(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.UserList(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

// 管理员修改用户状态
func UpdateStatusUser(ctx *gin.Context) {
	var req _map.UpdateStatus
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UpdateStatusUser(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.RestoreDeleteUser(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.UserById(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.GetAllRoles(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

// 获取角色下的所有用户
func UserByRole(ctx *gin.Context) {
	var req _map.RoleReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.UserByRole(ctx.Request.Context(), req.Content); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UserAddRoles(ctx.Request.Context(), req); err != nil {
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
		R.HandleBadRequest(ctx, nil)
		return
	}
	req.Uid = xcast.ToUint64(ctx.Param("uid"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UserDelRoles(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//菜单
//列表
func MenuList(ctx *gin.Context) {
	var req = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.MenuList(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

//删除
func DelMenu(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DelMenu(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//更新
func UpdateMenu(ctx *gin.Context) {
	var req _map.MenuReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	id := xcast.ToUint32(ctx.Param("id"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if id <= 0 {
		R.HandleBadRequest(ctx, "id is not cant 0")
		return
	}
	if err := admin_server.UpdateMenu(ctx.Request.Context(), id, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//添加
func AddMenu(ctx *gin.Context) {
	var req _map.MenuReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.AddMenu(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//api资源
func ResourcesList(ctx *gin.Context) {
	var req = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := admin_server.ResourcesList(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

func DelResources(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DelResources(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

func UpdateResources(ctx *gin.Context) {
	var req _map.ResourcesReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	id := xcast.ToUint32(ctx.Param("id"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if id <= 0 {
		R.HandleBadRequest(ctx, "id is not cant 0")
		return
	}
	if err := admin_server.UpdateResources(ctx.Request.Context(), id, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

func AddResources(ctx *gin.Context) {
	var req _map.ResourcesReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.AddResources(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 获取角色下的所有菜单权限
func GetPermissionAndMenuByRoles(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if data, err := auth_server.GetPermissionAndMenuByRoles(ctx.Request.Context(), []string{id}); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
}

// 更新角色下的所有菜单权限
func UpdateRolesMenuAndResources(ctx *gin.Context) {
	var req _map.UpdateRolesMenuAndResourcesReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.UpdateRolesMenuAndResources(ctx, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 添加角色
func AddRoles(ctx *gin.Context) {
	var req _map.RoleInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if err := admin_server.AddRoles(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

// 删除角色
func DelRoles(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DelRoles(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

// 更新角色
func UpdateRoles(ctx *gin.Context) {
	var req _map.RoleInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	id := xcast.ToUint32(ctx.Param("id"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if id <= 0 {
		R.HandleBadRequest(ctx, "id is not cant 0")
		return
	}
	if err := admin_server.UpdateRoles(ctx.Request.Context(), id, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//Agency
//列表
func AgencyList(ctx *gin.Context) {
	var req = _map.DefaultPageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	parentId := xcast.ToUint32(ctx.Param("pid"))
	if parentId < 0 {
		R.HandleBadRequest(ctx, "pid is greater than 0")
		return
	}
	if data, err := admin_server.AgencyList(ctx.Request.Context(), parentId, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Page(ctx, xcast.ToInt64(data.Count), req.Page, req.PageSize, data.Data)
	}
}

//删除
func DelAgency(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.DelAgency(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//更新
func UpdateAgency(ctx *gin.Context) {
	var req _map.AgencyInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	id := xcast.ToUint32(ctx.Param("id"))
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if id <= 0 {
		R.HandleBadRequest(ctx, "id is not cant 0")
		return
	}
	if err := admin_server.UpdateAgency(ctx.Request.Context(), id, req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//添加
func AddAgency(ctx *gin.Context) {
	var req _map.AgencyInfoReq
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if get, ok := ctx.Get("Uid"); ok {
		uid := xcast.ToUint32(get)
		if err := admin_server.AddAgency(ctx.Request.Context(), uid, req); err != nil {
			R.Error(ctx, err)
		} else {
			R.Ok(ctx, nil)
		}
	} else {
		R.HandleForbidden(ctx)
	}
	return
}

//恢复删除 RecoverDelAgency
func RecoverDelAgency(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.RecoverDelAgency(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//修改状态 UpdateAgencyStatus
func UpdateAgencyStatus(ctx *gin.Context) {
	var req _map.UpdateStatus
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if req.Uid <= 0 || (req.Status != 1 && req.Status != 2) {
		R.HandleBadRequest(ctx, "id >= 0 or status error")
		return
	}
	if err := admin_server.UpdateAgencyStatus(ctx.Request.Context(), xcast.ToUint32(req.Uid), req.Status); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, nil)
	}
	return
}

//退出 RemoveAgency
func RemoveAgency(ctx *gin.Context) {
	var req _map.UidList
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if data, err := admin_server.RemoveAgency(ctx.Request.Context(), req); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}

//加入的用户 ListUserByJoinAgency
func ListUserByJoinAgency(ctx *gin.Context) {
	var req _map.UpdateStatus
	if err := ctx.ShouldBind(&req); err != nil {
		R.HandleBadRequest(ctx, nil)
		return
	}
	if err := xvalidator.Struct(req); err != nil {
		R.HandleBadRequest(ctx, xvalidator.GetMsg(err).Error())
		return
	}
	if req.Uid <= 0 || (req.Status != 1 && req.Status != 2) {
		R.HandleBadRequest(ctx, "id >= 0 or status error")
		return
	}
	if data, err := admin_server.ListUserByJoinAgency(ctx.Request.Context(), xcast.ToUint32(req.Uid), req.Status); err != nil {
		R.Error(ctx, err)
	} else {
		R.Ok(ctx, data)
	}
	return
}
