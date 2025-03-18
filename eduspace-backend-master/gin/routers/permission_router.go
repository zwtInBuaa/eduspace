package routers

import (
	"EDU_TH_2_backend/gin/controllers"
	"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/repositories"
	"EDU_TH_2_backend/gin/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var permissionController *controllers.PermissionController

var permissionRoutes Routes

func NewPermissionRouter(db *gorm.DB) {
	var permissionRepo = repositories.NewPermissionRepository(enforcer)
	var permissionService = services.NewPermissionService(permissionRepo)
	permissionController = controllers.NewPermissionController(permissionService)

	permissionRoutes = Routes{
		{
			"AddPermission",
			"POST",
			"/addpermission",
			permissionController.AddPermission,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeletePermission",
			"DELETE",
			"/deletepermission",
			permissionController.DeletePermission,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateRole",
			"PUT",
			"/updaterole",
			permissionController.UpdateRole,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
