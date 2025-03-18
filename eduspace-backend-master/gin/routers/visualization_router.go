package routers

import (
	"EDU_TH_2_backend/gin/controllers"
	"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var visualizationController *controllers.VisualizationController

var visualizationRoutes Routes

func NewVisualizationRouter(db *gorm.DB) {
	// 初始化依赖实例
	var visualizationService = services.NewVisualizationService()
	visualizationController = controllers.NewVisualizationController(visualizationService)

	// 定义路由
	visualizationRoutes = Routes{
		{
			"BinarySearch_SubmitCode",
			http.MethodPost,
			"/binarySearch/submitCode",
			visualizationController.BinarySearchSubmitCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Sort_GenCode",
			http.MethodGet,
			"/sort/:type",
			visualizationController.SortGenCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Sort_SubmitCode",
			http.MethodPost,
			"/sort/submitCode/:type",
			visualizationController.SortSubmitCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Graph_GenCode",
			http.MethodGet,
			"/graph/:type",
			visualizationController.GraphGenCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Graph_SubmitCode",
			http.MethodPost,
			"/graph/submitCode/:type",
			visualizationController.GraphSubmitCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Custom_SubmitCode",
			http.MethodPost,
			"/custom/submitCode",
			visualizationController.CustomSubmitCode,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
