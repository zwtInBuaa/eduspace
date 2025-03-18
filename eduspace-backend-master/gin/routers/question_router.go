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

var questionController *controllers.QuestionController

var questionRoutes Routes

func NewQuestionRouter(db *gorm.DB) {
	// 初始化依赖实例
	var questionRepo = repositories.NewQuestionRepository(db)
	var questionService = services.NewQuestionService(questionRepo)
	questionController = controllers.NewQuestionController(questionService)

	// 定义路由
	questionRoutes = Routes{
		{
			"CreateQuestion",
			http.MethodPost,
			"/addquestion",
			questionController.CreateQuestionAndAddSource, // TODO
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetQuestion",
			http.MethodGet,
			"/:question_id",
			questionController.GetQuestion,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateQuestion",
			http.MethodPut,
			"/:question_id",
			questionController.UpdateQuestion,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteQuestionFromExam",
			http.MethodDelete,
			"/:question_id",
			questionController.DeleteQuestion,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllQuestions",
			http.MethodGet,
			"/getall",
			questionController.GetAllQuestions,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"Submit",
			http.MethodPost,
			"/:question_id/submit",
			questionController.Submit,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
	}

}
