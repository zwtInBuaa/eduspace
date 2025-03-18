package controllers

import (
	"EDU_TH_2_backend/gin/config"
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type VisualizationController struct {
	visualizationService services.VisualizationService
}

func NewVisualizationController(visualizationService services.VisualizationService) *VisualizationController {
	return &VisualizationController{visualizationService: visualizationService}
}

var lastSubmitTime int64 = 0

func (ctrl *VisualizationController) BinarySearchSubmitCode(c *gin.Context) {

	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "You can only submit your code once per 5 seconds.",
			})
			return
		}
	}

	// 获取参数值
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断参数中是否包含了名为"codeBlanks"的数组
	codeBlanks, ok := reqBody["codeBlanks"].([]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing or invalid 'codeBlanks' parameter"})
		return
	}

	// 将数组转换为字符串数组
	codeBlankString := make([]string, len(codeBlanks))
	for i, v := range codeBlanks {
		if s, ok := v.(string); ok {
			codeBlankString[i] = s
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'codeBlanks' parameter"})
			return
		}
	}

	data, err := ctrl.visualizationService.BinarySearchSubmitCode(codeBlankString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info("User submitted code for binary search visualization.")

	c.JSON(http.StatusOK, data)
}

func (ctrl *VisualizationController) SortGenCode(c *gin.Context) {

	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "You can only submit your code once per 5 seconds.",
			})
			return
		}
	}

	// 获取分类形式
	sortType := c.Param("type")

	data, err := ctrl.visualizationService.SortGenCode(sortType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data2, _ := data["inputDefault"].([]interface{})

	defaultSubmit := make([]string, len(data2))
	for i, v := range data2 {
		defaultSubmit[i] = v.(string)
	}

	defaultData, err := ctrl.visualizationService.SortSubmitCodeNotInDocker(sortType, defaultSubmit)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info(fmt.Sprintf("User generated code for %s sort visualization.", sortType))

	// 拼接data全部内容和defaultData的全部内容
	for k, v := range defaultData {
		data[k] = v
	}

	c.JSON(http.StatusOK, data)
}

func (ctrl *VisualizationController) SortSubmitCode(c *gin.Context) {

	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "You can only submit your code once per 5 seconds.",
			})
			return
		}
	}

	// 获取分类形式
	sortType := c.Param("type")

	// 获取参数值
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断参数中是否包含了名为"codeBlanks"的数组
	codeBlanks, ok := reqBody["codeBlanks"].([]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing or invalid 'codeBlanks' parameter"})
		return
	}

	// 将数组转换为字符串数组
	codeBlankString := make([]string, len(codeBlanks))
	for i, v := range codeBlanks {
		if s, ok := v.(string); ok {
			codeBlankString[i] = s
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'codeBlanks' parameter"})
			return
		}
	}

	data, err := ctrl.visualizationService.SortSubmitCode(sortType, codeBlankString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info(fmt.Sprintf("User submitted code for %s sort visualization.", sortType))

	c.JSON(http.StatusOK, data)
}

func (ctrl *VisualizationController) GraphGenCode(c *gin.Context) {
	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "You can only submit your code once per 5 seconds.",
			})
			return
		}
	}

	// 获取分类形式
	graphType := c.Param("type")

	data, err := ctrl.visualizationService.GraphGenCode(graphType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//data2, _ := data["inputDefault"].([]interface{})
	//
	//defaultSubmit := make([]string, len(data2))
	//for i, v := range data2 {
	//	defaultSubmit[i] = v.(string)
	//}
	//
	//defaultData, err := ctrl.visualizationService.SortSubmitCode(sortType, defaultSubmit)
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info(fmt.Sprintf("User generated code for %s graph visualization.", graphType))

	//c.JSON(http.StatusOK, gin.H{"code": data["code"], "inputDefault": data["inputDefault"],
	//	"values": defaultData["values"], "heights": defaultData["heights"],
	//	"steps": defaultData["steps"], "detail": defaultData["detail"]})
	c.JSON(http.StatusOK, data)
}

func (ctrl *VisualizationController) GraphSubmitCode(c *gin.Context) {

	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "你只能每5秒提交一次代码",
			})
			return
		}
	}

	// 获取分类形式
	graphType := c.Param("type")

	// 获取参数值
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 判断参数中是否包含了名为"codeBlanks"的数组
	codeBlanks, ok := reqBody["codeBlanks"].([]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing or invalid 'codeBlanks' parameter"})
		return
	}

	// 从 reqBody 中获取 graphInfo
	graphInfo, ok := reqBody["graphInfo"].(map[string]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid graphInfo"})
		return
	}

	// 将数组转换为字符串数组
	codeBlankString := make([]string, len(codeBlanks))
	for i, v := range codeBlanks {
		if s, ok := v.(string); ok {
			codeBlankString[i] = s
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid 'codeBlanks' parameter"})
			return
		}
	}

	graphInfoBytes, err := json.Marshal(graphInfo)
	if err != nil {
		fmt.Println("Failed to marshal graphInfo:", err)
		return
	}
	graphInfoStr := string(graphInfoBytes)

	data, err := ctrl.visualizationService.GraphSubmitCode(graphType, codeBlankString, graphInfoStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info(fmt.Sprintf("User submitted code for %s graph visualization.", graphType))

	c.JSON(http.StatusOK, data)

}

func (ctrl *VisualizationController) CustomSubmitCode(c *gin.Context) {
	// 限制 5s 内只能访问一次
	currentTime := time.Now().Unix()

	limitSubmissionTime := config.GetString("security.limit_submission_time") == "true"
	if limitSubmissionTime {
		if currentTime-lastSubmitTime < 5 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "你只能每5秒提交一次代码.",
			})
			return
		}
	}

	var customSubmitCodeForm *models.CustomSubmitCodeForm

	if err := c.ShouldBindJSON(&customSubmitCodeForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := ctrl.visualizationService.CustomSubmitCode(customSubmitCodeForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录上一次完成时间
	lastSubmitTime = currentTime

	// logs:哪个用户干了些什么（从auth中获取用户信息）(英文)
	logger.Info(fmt.Sprintf("User submitted code for custom graph visualization."))

	c.JSON(http.StatusOK, data)
}
