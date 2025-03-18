package controllers

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/services"
	"EDU_TH_2_backend/gin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	userService services.UserService
	tokenUtil   utils.TokenUtil
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) Signup(c *gin.Context) {
	var reqBody []map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ids := make([]uint, len(reqBody))

	for _, user := range reqBody {
		username, ok := user["username"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
			return
		}

		buaaId, ok := user["buaa_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid buaaId"})
			return
		}

		password, ok := user["password"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
			return
		}

		role_string, ok := user["role"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}

		role, _ := utils.ParseFloat64ToInt64(role_string)

		id, err := ctrl.userService.Signup(c, username, buaaId, password, role)
		ids = append(ids, id)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "signup success", "IDs": ids})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buaaId, ok := reqBody["buaa_id"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid buaaId"})
		return
	}

	password, ok := reqBody["password"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token, userId, username, buaaId, role, err := ctrl.userService.Login(buaaId, password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var role_string string
	if role == 0 {
		role_string = "管理员"
	} else if role == 1 {
		role_string = "老师"
	} else if role == 2 {
		role_string = "助教"
	} else if role == 3 {
		role_string = "学生"
	}

	courses, err := ctrl.userService.GetStudentCourses(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	logger.Info("Login: " + username + " " + buaaId + " " + c.ClientIP())

	avatar, err := utils.GetAvatarURL(c, strconv.Itoa(int(userId)))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user_id": userId, "username": username, "buaa_id": buaaId, "role": role_string,
		"avatar": avatar, "courses": courses})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	// 从请求头中获取token
	token := c.GetString("tokenString")

	// 废弃token
	utils.InvalidateToken(token)

	// 设置响应头，删除token
	c.Header("Authorization", "")

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout success",
	})
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	user, err := ctrl.userService.GetUserByID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	user, err := ctrl.userService.GetAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, ok := reqBody["username"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid username"})
		return
	}

	buaaId, ok := reqBody["buaa_id"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid buaaId"})
		return
	}
	changeUserId := c.GetInt("userId")
	if err := ctrl.userService.UpdateUser(userId, username, buaaId, uint(changeUserId)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update success"})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	changeUserId := c.GetInt("userId")
	changeUserRole := c.GetInt("role")
	if changeUserRole != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "你没有权限删除用户"})
		return
	}

	if err := ctrl.userService.DeleteUser(c, userId, uint(changeUserId)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete success"})
}

func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("user_id")

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fileBase64
	image, ok := reqBody["image"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid image(base64)"})
		return
	}

	// 调用Service层的方法完成头像上传
	savePath, err := ctrl.userService.SaveAvatar(c, userID, image)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log:哪个用户
	logger.Info("SaveAvatar: " + userID)

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"message": "upload avatar successfully",
		"path":    savePath,
	})
}

func (ctrl *UserController) GetAvatar(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("user_id")

	// 调用service方法处理业务逻辑
	avatarURL, err := ctrl.userService.GetAvatar(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}

func (ctrl *UserController) GetTeacherCourses(c *gin.Context) {
	teacherID, err := utils.ParseStringToUint(c.Param("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	courses, err := ctrl.userService.GetTeacherCourses(teacherID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (ctrl *UserController) GetStudentCourses(c *gin.Context) {
	studentID, err := utils.ParseStringToUint(c.Param("user_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	courses, err := ctrl.userService.GetStudentCourses(studentID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (ctrl *UserController) UpdatePassword(c *gin.Context) {
	// 获取用户ID
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))
	changeUserId := c.GetInt("userId")
	if userId != uint(changeUserId) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "你没有权限修改别人的密码"})
		return
	}

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldPassword, ok := reqBody["old_password"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "旧密码格式错误"})
		return
	}

	newPassword, ok := reqBody["new_password"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "新密码格式错误"})
		return
	}

	if err := ctrl.userService.UpdatePassword(userId, oldPassword, newPassword); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update password success"})
}

func (ctrl *UserController) GetWeakness(c *gin.Context) {
	// 获取用户ID
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	weakness, err := ctrl.userService.GetWeakness(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weakness": weakness})
}

func (ctrl *UserController) ResetPasswordByBuaaId(c *gin.Context) {
	var resetPasswordFrom *models.ResetPasswordForm

	changeUserId := c.GetInt("userId")

	if err := c.ShouldBindJSON(&resetPasswordFrom); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
		// 此外可以使用 validator 库进行参数校验
	}
	// logs TODO
	if err := ctrl.userService.ResetPasswordByBuaaId(resetPasswordFrom, uint(changeUserId)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reset password success"})
}

func (ctrl *UserController) RecQuestion(c *gin.Context) {
	// 获取用户ID
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	questionList, err := ctrl.userService.RecQuestion(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回id, title, tags
	var returnQuestionList []map[string]interface{}
	for i := range questionList {
		returnQuestionList = append(returnQuestionList, map[string]interface{}{
			"problem_id":    questionList[i].ID,
			"problem_title": questionList[i].Title,
			"tags":          questionList[i].Tags,
		})
	}

	// 如果没有推荐的题目，返回空数组(不是string), 而不是null
	if returnQuestionList == nil {
		returnQuestionList = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"questionList": returnQuestionList})

}

func (ctrl *UserController) QuestionOverview(c *gin.Context) {
	// 获取用户ID
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	questionOverview, err := ctrl.userService.QuestionOverview(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"overview": questionOverview})
}

func (ctrl *UserController) GetWeaknessLabel(c *gin.Context) {
	// 获取用户ID
	userId, _ := utils.ParseStringToUint(c.Param("user_id"))

	weaknessLabel, err := ctrl.userService.GetWeaknessLabel(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weaknessLabel)
}
