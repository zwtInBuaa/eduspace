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

var examController *controllers.ExamController

var examRoutes Routes

func NewExamRouter(db *gorm.DB) {

	var examRepo = repositories.NewExamRepository(db)
	var examService = services.NewExamService(examRepo)
	examController = controllers.NewExamController(examService)

	examRoutes = Routes{
		{
			"CreateExam",
			http.MethodPost,
			"/addexam",
			examController.CreateExamWithQuestions, // TODO
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetExam",
			http.MethodGet,
			"/:exam_id",
			examController.GetExamWithQuestionsByExamID,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"UpdateExam",
			http.MethodPut,
			"/:exam_id",
			examController.UpdateExamWithQuestions,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteExam",
			http.MethodDelete,
			"/:exam_id",
			examController.DeleteExam,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"GetAllExams",
			http.MethodGet,
			"/getall",
			examController.GetAllExams,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"AddQuestionsToExam",
			http.MethodPost,
			"/:exam_id/questions",
			examController.AddQuestionsToExam,
			[]gin.HandlerFunc{middlewares.JWTAuthMiddleware(tokenUtil)},
		},
		{
			"DeleteQuestionFromExam",
			http.MethodDelete,
			"/:exam_id/questions/:question_id",
			examController.DeleteQuestionFromExam,
			[]gin.HandlerFunc{},
		},
		{
			"FindQuestionsByExamID",
			http.MethodGet,
			"/:exam_id/questions",
			examController.FindQuestionsByExamID,
			[]gin.HandlerFunc{},
		},
	}
}
