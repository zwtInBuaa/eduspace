package controllers

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/services"
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExamController struct {
	examService services.ExamService
}

func NewExamController(examService services.ExamService) *ExamController {
	return &ExamController{examService: examService}
}

func (ctrl *ExamController) CreateExam(c *gin.Context) {
	var exam models.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.examService.CreateExam(&exam); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户创建了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d create exam %d", userId, exam.ID))

	c.JSON(http.StatusOK, gin.H{"message": "create exam success", "exam": exam})
}

func (ctrl *ExamController) CreateExamWithQuestions(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam := models.Exam{
		Name:        reqBody["name"].(string),
		Description: reqBody["description"].(string),
	}

	// 将题目ID数组存储在 Exam 模型中
	questions, ok := reqBody["questions"].([]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid questions"})
		return
	}

	if err := ctrl.examService.CreateExamWithQuestions(&exam, questions); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户创建了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d create exam %d", userId, exam.ID))

	c.JSON(http.StatusOK, gin.H{"message": "create exam success", "exam": exam})
}

func (ctrl *ExamController) GetExam(c *gin.Context) {
	examId, _ := utils.ParseStringToUint(c.Param("exam_id"))

	exam, err := ctrl.examService.GetExamByID(examId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exam)
}

func (ctrl *ExamController) UpdateExam(c *gin.Context) {
	examID, err := utils.ParseStringToUint(c.Param("exam_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exam models.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam.ID = examID

	if err := ctrl.examService.UpdateExam(&exam); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户更新了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d update exam %d", userId, exam.ID))

	c.JSON(http.StatusOK, gin.H{"message": "update exam success", "exam": exam})
}

func (ctrl *ExamController) UpdateExamWithQuestions(c *gin.Context) {
	examID, err := utils.ParseStringToUint(c.Param("exam_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam := models.Exam{
		Name:        reqBody["name"].(string),
		Description: reqBody["description"].(string),
	}

	exam.ID = examID

	// 将题目ID数组存储在 Exam 模型中
	questions, ok := reqBody["questions"].([]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid questions"})
		return
	}

	if err := ctrl.examService.UpdateExamWithQuestions(&exam, questions); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户更新了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d update exam %d", userId, exam.ID))

	c.JSON(http.StatusOK, gin.H{"message": "update exam success", "exam": exam})
}

func (ctrl *ExamController) DeleteExam(c *gin.Context) {
	examId, _ := utils.ParseStringToUint(c.Param("exam_id"))

	if err := ctrl.examService.DeleteExam(examId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户删除了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete exam %d", userId, examId))

	c.JSON(http.StatusOK, gin.H{"message": "delete exam success"})
}

func (ctrl *ExamController) AddQuestionsToExam(c *gin.Context) {
	examID, err := utils.ParseStringToUint(c.Param("exam_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var questionIDs []uint
	if err := c.ShouldBindJSON(&questionIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, questionID := range questionIDs {
		if err := ctrl.examService.AddQuestion(examID, questionID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户添加了哪些题目到哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d add questions %v to exam %d", userId, questionIDs, examID))

	c.JSON(http.StatusOK, gin.H{"message": "Questions added to exam"})
}

func (ctrl *ExamController) DeleteQuestionFromExam(c *gin.Context) {
	examID, err := utils.ParseStringToUint(c.Param("exam_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionID, err := utils.ParseStringToUint(c.Param("question_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.examService.RemoveQuestion(examID, questionID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户从哪个题组删除了哪个题目
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete question %d from exam %d", userId, questionID, examID))

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted from exam"})
}

func (ctrl *ExamController) GetAllExams(c *gin.Context) {
	exams, err := ctrl.examService.GetAllExams()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exams)
}

func (ctrl *ExamController) FindQuestionsByExamID(c *gin.Context) {
	examId, _ := utils.ParseStringToUint(c.Param("exam_id"))

	questions, err := ctrl.examService.FindQuestionsByExamID(examId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (ctrl *ExamController) GetExamWithQuestionsByExamID(c *gin.Context) {
	examId, _ := utils.ParseStringToUint(c.Param("exam_id"))

	exam, err := ctrl.examService.GetExamByID(examId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questions, err := ctrl.examService.FindQuestionsByExamID(examId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var questionsInfo []map[string]interface{}

	for i := range questions {
		var questionInfo = map[string]interface{}{
			"id":    questions[i].ID,
			"title": questions[i].Title,
			"tags":  questions[i].Tags,
		}
		questionsInfo = append(questionsInfo, questionInfo)
	}

	examInfo := map[string]interface{}{
		"id":          exam.ID,
		"title":       exam.Name,
		"description": exam.Description,
		"questions":   questionsInfo,
	}

	c.JSON(http.StatusOK, examInfo)
}
