package routers

import (
	"EDU_TH_2_backend/gin/controllers"
	"EDU_TH_2_backend/gin/middlewares"
	"EDU_TH_2_backend/gin/repositories"
	"EDU_TH_2_backend/gin/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var commentController *controllers.CommentController

var commentRoutes Routes

func NewCommentRouter(db *gorm.DB) {
	// 初始化依赖实例
	var commentRepo = repositories.NewCommentRepository(db)
	var commentService = services.NewCommentService(commentRepo)
	commentController = controllers.NewCommentController(commentService)

	// 定义路由
	commentRoutes = Routes{
		{
			"CreateComment",
			http.MethodPost,
			"/addcomment",
			commentController.CreateComment,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetComment",
			http.MethodGet,
			"/:comment_id",
			commentController.GetComment,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllComments",
			http.MethodGet,
			"/getall/:post_id",
			commentController.GetCommentsByPost,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateComment",
			http.MethodPut,
			"/:comment_id",
			commentController.UpdateComment,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteComment",
			http.MethodDelete,
			"/:comment_id",
			commentController.DeleteComment,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
