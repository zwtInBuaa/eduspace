package controllers

import (
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PermissionController struct {
	service services.PermissionService
}

func NewPermissionController(service services.PermissionService) *PermissionController {
	return &PermissionController{service: service}
}

// 添加用户权限
func (ctrl *PermissionController) AddPermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := ctrl.service.AddPermission(&permission); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "add permission success"})
}

// 禁用用户权限
func (ctrl *PermissionController) DeletePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := ctrl.service.DeletePermission(&permission); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete permission success"})
}

// 更新用户权限
func (ctrl *PermissionController) UpdateRole(c *gin.Context) {
	var rolePermission models.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := ctrl.service.UpdateRole(&rolePermission); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update role success"})
}
