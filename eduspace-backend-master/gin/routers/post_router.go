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

var postController *controllers.PostController

var postRoutes Routes

func NewPostRouter(db *gorm.DB) {
	// 初始化依赖实例
	var postRepo = repositories.NewPostRepository(db)

	// user 依赖注入
	var userRepo = repositories.NewUserRepository(db)

	var postService = services.NewPostService(postRepo, userRepo)
	postController = controllers.NewPostController(postService)

	// 定义路由
	postRoutes = Routes{
		{
			"GetPost",
			http.MethodGet,
			"/:post_id",
			postController.GetPost,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllPosts",
			http.MethodGet,
			"/getall",
			postController.GetAllPosts,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"CreatePost",
			http.MethodPost,
			"/postblog",
			postController.CreatePost,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeletePost",
			http.MethodDelete,
			"/:post_id",
			postController.DeletePost,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdatePost",
			http.MethodPut,
			"/:post_id",
			postController.UpdatePost,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
