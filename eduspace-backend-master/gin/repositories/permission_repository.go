package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"github.com/casbin/casbin/v2"
	"strconv"
)

type PermissionRepository interface {
	AddPermission(permission *models.Permission) error
	DeletePermission(permission *models.Permission) error
	UpdateRole(rolePermission *models.RolePermission) error
}

type permissionRepository struct {
	enforcer *casbin.Enforcer
}

func NewPermissionRepository(enforcer *casbin.Enforcer) PermissionRepository {
	return &permissionRepository{enforcer: enforcer}
}

// casbin 向数据库中添加单个路由权限的操作
func (r *permissionRepository) AddPermission(permission *models.Permission) error {
	// 添加用户权限策略
	if _, err := r.enforcer.AddPolicy(strconv.FormatUint(uint64(permission.UserID), 10), permission.Path, permission.Method); err != nil {
		return err
	}

	return nil
}

// casbin 向数据库中删除单个路由权限的操作
func (r *permissionRepository) DeletePermission(permission *models.Permission) error {
	// 删除用户权限策略
	if _, err := r.enforcer.RemovePolicy(strconv.FormatUint(uint64(permission.UserID), 10), permission.Path, permission.Method); err != nil {
		return err
	}

	return nil
}

// casbin 向数据库中更改该用户角色的操作
func (r *permissionRepository) UpdateRole(rolePermission *models.RolePermission) error {
	// 删除用户权限策略
	if _, err := r.enforcer.RemoveGroupingPolicy(rolePermission.UserID, rolePermission.Role); err != nil {
		return err
	}

	// 添加用户权限策略
	if _, err := r.enforcer.AddGroupingPolicy(rolePermission.UserID, rolePermission.Role); err != nil {
		return err
	}

	return nil
}
