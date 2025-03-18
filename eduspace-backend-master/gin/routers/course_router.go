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

var courseController *controllers.CourseController

var courseRoutes Routes

func NewCourseRouter(db *gorm.DB) {
	// 初始化依赖实例
	var courseRepo = repositories.NewCourseRepository(db)
	var courseService = services.NewCourseService(courseRepo)
	courseController = controllers.NewCourseController(courseService)

	courseRoutes = Routes{
		{
			"CreateCourse",
			http.MethodPost,
			"/signup",
			courseController.CreateCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetCourse",
			http.MethodGet,
			"/:course_id",
			courseController.GetCourseByID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllCourses",
			http.MethodGet,
			"/getall",
			courseController.GetAllCourses,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateCourse",
			http.MethodPut,
			"/:course_id",
			courseController.UpdateCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteCourse",
			http.MethodDelete,
			"/:course_id",
			courseController.DeleteCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetTeachersByCourseID",
			http.MethodGet,
			"/:course_id/teachers",
			courseController.GetTeachersByCourseID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetStudentsByCourseID",
			http.MethodGet,
			"/:course_id/students",
			courseController.GetStudentsByCourseID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"AddTeachersToCourse",
			http.MethodPost,
			"/:course_id/add_teachers",
			courseController.AddTeachersToCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"AddStudentsToCourse",
			http.MethodPost,
			"/:course_id/add_students",
			courseController.AddStudentsToCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteStudentFromCourse",
			http.MethodDelete,
			"/:course_id/students/:student_id",
			courseController.DeleteStudentFromCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteTeacherFromCourse",
			http.MethodDelete,
			"/:course_id/teachers/:teacher_id",
			courseController.DeleteTeacherFromCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"AddQuestionsToCourse",
			http.MethodPost,
			"/:course_id/questions",
			courseController.AddQuestionsToCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteQuestionFromCourse",
			http.MethodDelete,
			"/:course_id/questions/:question_id",
			courseController.DeleteQuestionFromCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"AddExamsToCourse",
			http.MethodPost,
			"/:course_id/exams",
			courseController.AddExamsToCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteExamFromCourse",
			http.MethodDelete,
			"/:course_id/exams/:exam_id",
			courseController.DeleteExamFromCourse,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"FindQuestionsByCourseID",
			http.MethodGet,
			"/:course_id/questions",
			courseController.FindQuestionsByCourseID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"FindExamsByCourseID",
			http.MethodGet,
			"/:course_id/exams",
			courseController.FindExamsByCourseID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}

}
