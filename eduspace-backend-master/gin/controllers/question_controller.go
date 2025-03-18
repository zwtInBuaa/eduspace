package controllers

import (
	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/services"
	"EDU_TH_2_backend/gin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type QuestionController struct {
	questionService services.QuestionService
}

func NewQuestionController(questionService services.QuestionService) *QuestionController {
	return &QuestionController{questionService: questionService}
}

func (ctrl *QuestionController) CreateQuestion(c *gin.Context) {
	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.questionService.CreateQuestion(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户创建了什么问题（英文）
	userId := c.GetInt("user_id")
	logger.Info(fmt.Sprintf("user %d created question %d", userId, question.ID))

	c.JSON(http.StatusOK, gin.H{"message": "create question success", "question": question})
}

func (ctrl *QuestionController) CreateQuestionAndAddSource(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(reqBody["course_id"])
	courseIDFloat64, ok := reqBody["course_id"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
		return
	}
	courseID, _ := utils.ParseFloat64ToUint(courseIDFloat64)

	// 将请求体转换为Question结构体
	var question models.Question
	if err := mapstructure.Decode(reqBody, &question); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.questionService.CreateQuestionAndAddSource(courseID, &question); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户创建了什么问题（英文）
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d created question %d", userId, question.ID))

	c.JSON(http.StatusOK, gin.H{"message": "create question and add source success", "question": question})
}

func (ctrl *QuestionController) GetQuestion(c *gin.Context) {
	questionId, _ := utils.ParseStringToUint(c.Param("question_id"))

	question, err := ctrl.questionService.GetQuestionByID(questionId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (ctrl *QuestionController) UpdateQuestion(c *gin.Context) {
	questionId, _ := utils.ParseStringToUint(c.Param("question_id"))

	var question models.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question.ID = questionId

	if err := ctrl.questionService.UpdateQuestion(&question); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户更新了哪个问题
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d update question %d", userId, questionId))

	c.JSON(http.StatusOK, gin.H{"message": "update question success"})
}

func (ctrl *QuestionController) DeleteQuestion(c *gin.Context) {
	questionId, _ := utils.ParseStringToUint(c.Param("question_id"))

	if err := ctrl.questionService.DeleteQuestion(questionId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户删除了哪个问题
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete question %d", userId, questionId))

	c.JSON(http.StatusOK, gin.H{"message": "delete question success"})
}

func (ctrl *QuestionController) GetAllQuestions(c *gin.Context) {
	questions, err := ctrl.questionService.GetAllQuestions()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}

func (ctrl *QuestionController) Submit(c *gin.Context) {
	// 限制 1s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 1 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "你只能每1秒提交一次答案",
			})
			return
		}
	}

	userId := c.GetInt("userId")

	questionId, _ := utils.ParseStringToUint(c.Param("question_id"))

	var submitForm *models.SubmitFrom

	if err := c.ShouldBindJSON(&submitForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	submitInfo, err := ctrl.questionService.Submit(uint(userId), questionId, submitForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	c.JSON(http.StatusOK, submitInfo)
}
