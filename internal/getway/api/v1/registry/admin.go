package registry

import (
	ah "github.com/myxy99/ndisk/internal/getway/api/v1/handle/admin"
	"github.com/myxy99/ndisk/internal/getway/api/v1/middleware"
)

func init() {
	admin := V1().Group("/admin")
	// 管理员创建用户
	admin.POST("/users", middleware.Auth(), middleware.Authority(), ah.CreateUser)
	// 管理员用户信息修改
	admin.PUT("/user/:uid", middleware.Auth(), middleware.Authority(), ah.UpdateUser)
	// 管理员删除用户
	admin.DELETE("/users", middleware.Auth(), middleware.Authority(), ah.DeleteUser)
	// 管理员用户列表
	admin.GET("/users", middleware.Auth(), middleware.Authority(), ah.UserList)
	// 管理员修改用户状态
	admin.PATCH("/user/:uid/status", middleware.Auth(), middleware.Authority(), ah.UpdateStatusUser)
	// 管理员恢复已经删除用户
	admin.PATCH("/users/restore", middleware.Auth(), middleware.Authority(), ah.RestoreDeleteUser)
	// 管理员获取用户信息
	admin.GET("/user/:uid", middleware.Auth(), middleware.Authority(), ah.UserById)

	// 管理员获取全部角色RoleList
	admin.GET("/role", middleware.Auth(), middleware.Authority(), ah.RoleList)
	// 获取角色下的所有用户
	admin.GET("/role/user", middleware.Auth(), middleware.Authority(), ah.UserByRole)
	// 获取角色下的所有菜单权限
	admin.GET("/role/info", middleware.Auth(), middleware.Authority(), ah.GetPermissionAndMenuByRoles)
	// 更新角色下的所有菜单权限
	admin.POST("/role/info", middleware.Auth(), middleware.Authority(), ah.UpdateRolesMenuAndResources)
	// 添加角色
	admin.POST("/role", middleware.Auth(), middleware.Authority(), ah.AddRoles)
	// 删除角色
	admin.DELETE("/role", middleware.Auth(), middleware.Authority(), ah.DelRoles)
	// 更新角色
	admin.PUT("/role/:id", middleware.Auth(), middleware.Authority(), ah.UpdateRoles)

	// 管理员给用户添加角色
	admin.POST("/user/:uid/role", middleware.Auth(), middleware.Authority(), ah.UserAddRoles)
	// 管理员删除用户角色
	admin.DELETE("/user/:uid/role", middleware.Auth(), middleware.Authority(), ah.DeleteUserRole)

	//menu
	//add
	admin.POST("/menu", middleware.Auth(), middleware.Authority(), ah.AddMenu)
	//del
	admin.DELETE("/menu", middleware.Auth(), middleware.Authority(), ah.DelMenu)
	//get
	admin.GET("/menu", middleware.Auth(), middleware.Authority(), ah.MenuList)
	//update
	admin.PUT("/menu/:id", middleware.Auth(), middleware.Authority(), ah.UpdateMenu)

	//Resources
	//add
	admin.POST("/resources", middleware.Auth(), middleware.Authority(), ah.AddResources)
	//del
	admin.DELETE("/resources", middleware.Auth(), middleware.Authority(), ah.DelResources)
	//get
	admin.GET("/resources", middleware.Auth(), middleware.Authority(), ah.ResourcesList)
	//update
	admin.PUT("/resources/:id", middleware.Auth(), middleware.Authority(), ah.UpdateResources)

	//机构 	agency
	//todo
	//add
	admin.POST("/agency", middleware.Auth(), middleware.Authority(), ah.AddAgency)
	//del
	admin.DELETE("/agency", middleware.Auth(), middleware.Authority(), ah.DelAgency)
	//get
	admin.GET("/agency/:pid", middleware.Auth(), middleware.Authority(), ah.AgencyList)
	//update
	admin.PUT("/agency/:id", middleware.Auth(), middleware.Authority(), ah.UpdateAgency)
	//修改状态
	admin.PATCH("/agency/:id/status", middleware.Auth(), middleware.Authority(), ah.UpdateAgencyStatus)
	//恢复删除
	admin.PATCH("/agency_res", middleware.Auth(), middleware.Authority(), ah.RecoverDelAgency)
	//把用户剔除组织
	admin.POST("/agency/remove", middleware.Auth(), middleware.Authority(), ah.RemoveAgency)
	//获取加入的用户
	admin.GET("/agency_user", middleware.Auth(), middleware.Authority(), ah.ListUserByJoinAgency)
}
