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

var userController *controllers.UserController

var userRoutes Routes

func NewUserRouter(db *gorm.DB) {
	// 初始化依赖实例
	var userRepo = repositories.NewUserRepository(db)
	var userService = services.NewUserService(userRepo, tokenUtil)
	userController = controllers.NewUserController(userService)

	// 定义路由
	userRoutes = Routes{
		{
			"SignUp",
			http.MethodPost,
			"/signup",
			userController.Signup,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Login",
			http.MethodPost,
			"/login",
			userController.Login,
			[]gin.HandlerFunc{},
		},
		{
			"Logout",
			http.MethodPost,
			"/logout",
			userController.Logout,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetUser",
			http.MethodGet,
			"/:user_id",
			userController.GetUser,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllUser",
			http.MethodGet,
			"/getall",
			userController.GetAllUsers,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateUser",
			http.MethodPut,
			"/:user_id",
			userController.UpdateUser,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteUser",
			http.MethodDelete,
			"/:user_id",
			userController.DeleteUser,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UploadAvatar",
			http.MethodPost,
			"/upload-avatar/:user_id",
			userController.UploadAvatar,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAvatar",
			http.MethodGet,
			"/get-avatar/:user_id",
			userController.GetAvatar,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetTeacherCourses",
			http.MethodGet,
			"/:user_id/teacher_courses",
			userController.GetTeacherCourses,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetStudentCourses",
			http.MethodGet,
			"/:user_id/student_courses",
			userController.GetStudentCourses,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdatePassword",
			http.MethodPut,
			"/:user_id/changePassword",
			userController.UpdatePassword,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetWeakness",
			http.MethodGet,
			"/:user_id/weakness",
			userController.GetWeaknessLabel,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"RecQuestion",
			http.MethodGet,
			"/:user_id/recQuestion",
			userController.RecQuestion,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"QuestionOverview",
			http.MethodGet,
			"/:user_id/questionOverview",
			userController.QuestionOverview,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"ResetPasswordByBuaaId",
			http.MethodPost,
			"/resetPassword",
			userController.ResetPasswordByBuaaId,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}
}
