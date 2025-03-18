package services

import (
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
)

type PermissionService interface {
	AddPermission(permission *models.Permission) error
	DeletePermission(permission *models.Permission) error
	UpdateRole(rolePermission *models.RolePermission) error
}

type permissionService struct {
	permissionRepo repositories.PermissionRepository
}

func NewPermissionService(permissionRepo repositories.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) AddPermission(permission *models.Permission) error {
	return s.permissionRepo.AddPermission(permission)
}

func (s *permissionService) DeletePermission(permission *models.Permission) error {
	return s.permissionRepo.DeletePermission(permission)
}

func (s *permissionService) UpdateRole(rolePermission *models.RolePermission) error {
	return s.permissionRepo.UpdateRole(rolePermission)
}
