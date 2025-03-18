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

type CourseController struct {
	courseService services.CourseService
}

func NewCourseController(courseService services.CourseService) *CourseController {
	return &CourseController{courseService: courseService}
}

func (ctrl *CourseController) CreateCourse(c *gin.Context) {
	var courses []models.Course
	if err := c.ShouldBindJSON(&courses); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range courses {
		if err := ctrl.courseService.CreateCourse(&courses[i]); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户创建了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d create course %d", userId, courses[0].ID))

	c.JSON(http.StatusOK, gin.H{"message": "Courses created"})
}

func (ctrl *CourseController) GetCourseByID(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("course_id"))

	course, err := ctrl.courseService.GetCourseByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (ctrl *CourseController) UpdateCourse(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("course_id"))

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course.ID = id
	if err := ctrl.courseService.UpdateCourse(&course); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户更新了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d update course %d", userId, course.ID))

	c.JSON(http.StatusOK, gin.H{"message": "Course updated", "data": course})
}

func (ctrl *CourseController) DeleteCourse(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("course_id"))

	if err := ctrl.courseService.DeleteCourse(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户删除了哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete course %d", userId, id))

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted"})
}

func (ctrl *CourseController) GetAllCourses(c *gin.Context) {
	courses, err := ctrl.courseService.GetAllCourses()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (ctrl *CourseController) GetTeachersByCourseID(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("course_id"))

	teachers, err := ctrl.courseService.GetTeachersByCourseID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, teachers)
}

func (ctrl *CourseController) GetStudentsByCourseID(c *gin.Context) {
	id, _ := utils.ParseStringToUint(c.Param("course_id"))

	students, err := ctrl.courseService.GetStudentsByCourseID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (ctrl *CourseController) AddStudentsToCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var studentIDs []uint
	if err := c.ShouldBindJSON(&studentIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, studentID := range studentIDs {
		if err := ctrl.courseService.AddStudentToCourse(courseID, studentID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户添加了哪些学生到哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d add students %v to course %d", userId, studentIDs, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Students added to course"})
}

func (ctrl *CourseController) DeleteStudentFromCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := utils.ParseStringToUint(c.Param("student_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.courseService.RemoveStudentFromCourse(courseID, studentID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户从哪个题组删除了哪个学生
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete student %d from course %d", userId, studentID, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted from course"})
}

func (ctrl *CourseController) AddTeachersToCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var teacherIDs []uint
	if err := c.ShouldBindJSON(&teacherIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, teacherID := range teacherIDs {
		if err := ctrl.courseService.AddTeacherToCourse(courseID, teacherID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户添加了哪些老师到哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d add teachers %v to course %d", userId, teacherIDs, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Teachers added to course"})
}

func (ctrl *CourseController) DeleteTeacherFromCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID, err := utils.ParseStringToUint(c.Param("teacher_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.courseService.RemoveTeacherFromCourse(courseID, teacherID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户从哪个题组删除了哪个老师
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete teacher %d from course %d", userId, teacherID, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted from course"})
}

func (ctrl *CourseController) AddQuestionsToCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
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
		if err := ctrl.courseService.AddQuestionToCourse(courseID, questionID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户添加了哪些题目到哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d add questions %v to course %d", userId, questionIDs, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Questions added to course"})
}

func (ctrl *CourseController) DeleteQuestionFromCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questionID, err := utils.ParseStringToUint(c.Param("question_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.courseService.RemoveQuestionFromCourse(courseID, questionID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户从哪个题组删除了哪个题目
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete question %d from course %d", userId, questionID, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted from course"})
}

func (ctrl *CourseController) AddExamsToCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var examIDs []uint
	if err := c.ShouldBindJSON(&examIDs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, examID := range examIDs {
		if err := ctrl.courseService.AddExamToCourse(courseID, examID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// log: 哪个用户添加了哪些考试到哪个题组
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d add exams %v to course %d", userId, examIDs, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Exams added to course"})
}

func (ctrl *CourseController) DeleteExamFromCourse(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	examID, err := utils.ParseStringToUint(c.Param("exam_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.courseService.RemoveExamFromCourse(courseID, examID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// log: 哪个用户从哪个题组删除了哪个考试
	userId := c.GetInt("userId")
	logger.Info(fmt.Sprintf("user %d delete exam %d from course %d", userId, examID, courseID))

	c.JSON(http.StatusOK, gin.H{"message": "Exam deleted from course"})
}

func (ctrl *CourseController) FindQuestionsByCourseID(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := ctrl.courseService.FindQuestionsByCourseID(courseID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (ctrl *CourseController) FindExamsByCourseID(c *gin.Context) {
	courseID, err := utils.ParseStringToUint(c.Param("course_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exams, err := ctrl.courseService.FindExamsByCourseID(courseID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exams)
}
