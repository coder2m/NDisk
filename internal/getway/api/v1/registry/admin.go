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

	// 管理员获取所有权限列表
	admin.GET("/authority", middleware.Auth(), middleware.Authority(), ah.CompetenceList)
	// 管理员添加权限列表
	admin.POST("/authority", middleware.Auth(), middleware.Authority(), ah.AddCompetence)
	// 管理员删除权限
	admin.DELETE("/authority", middleware.Auth(), middleware.Authority(), ah.DeleteCompetence)

	// 管理员获取全部角色RoleList
	admin.GET("/role", middleware.Auth(), middleware.Authority(), ah.RoleList)
	// 管理员获取角色下的权限
	admin.GET("/role2authority", middleware.Auth(), middleware.Authority(), ah.CompetenceByRole)
	// 管理员给角色添加权限
	admin.POST("/role/authority", middleware.Auth(), middleware.Authority(), ah.RoleAddCompetence)
	// 获取角色下的所有用户
	admin.GET("/role/authority", middleware.Auth(), middleware.Authority(), ah.UserByRole)

	// 删除角色指定权限
	admin.DELETE("/role/authority", middleware.Auth(), middleware.Authority(), ah.DeleteRoleCompetence)

	// 管理员获取用户的角色
	admin.GET("/user/:uid/role", middleware.Auth(), middleware.Authority(), ah.RoleByUser)
	// 管理员给用户添加角色
	admin.POST("/user/:uid/role", middleware.Auth(), middleware.Authority(), ah.UserAddRoles)
	// 管理员删除用户角色
	admin.DELETE("/user/:uid/role", middleware.Auth(), middleware.Authority(), ah.DeleteUserRole)
}
