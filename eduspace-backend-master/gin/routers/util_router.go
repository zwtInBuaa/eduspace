package routers

import (
	"EDU_TH_2_backend/gin/controllers"
	"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var utilController *controllers.UtilController

var utilRoutes Routes

func NewUtilRouter(db *gorm.DB) {
	// 初始化依赖实例
	var utilService = services.NewUtilService()
	utilController = controllers.NewUtilController(utilService)

	// 定义路由
	utilRoutes = Routes{
		{
			"GetQRCode",
			http.MethodGet,
			"/getqrcode",
			utilController.GetQRCode,
			[]gin.HandlerFunc{},
		},
		{
			"Imgsave",
			http.MethodPost,
			"/imgdbs",
			utilController.ImgSave,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
