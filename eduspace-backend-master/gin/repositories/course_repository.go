package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"fmt"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindByID(id uint) (*models.Course, error)
	Update(course *models.Course) error
	Delete(id uint) error
	FindAll() ([]*models.Course, error)
	FindTeachersByCourseID(courseID uint) ([]*models.User, error)
	FindStudentsByCourseID(courseID uint) ([]*models.User, error)
	AddStudentToCourse(courseID uint, studentID uint) error
	AddTeacherToCourse(courseID uint, teacherID uint) error
	AddQuestionToCourse(courseID uint, questionID uint) error
	RemoveQuestionFromCourse(courseID uint, questionID uint) error
	AddExamToCourse(courseID uint, examID uint) error
	RemoveExamFromCourse(courseID uint, examID uint) error
	RemoveTeacherFromCourse(courseID uint, teacherID uint) error
	RemoveStudentFromCourse(courseID uint, studentID uint) error
	FindQuestionsByCourseID(courseID uint) ([]*models.Question, error)
	FindExamsByCourseID(courseID uint) ([]*models.Exam, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (repo *courseRepository) Create(course *models.Course) error {
	return repo.db.Create(course).Error
}

func (repo *courseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course
	err := repo.db.Where("id = ?", id).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (repo *courseRepository) Update(course *models.Course) error {
	return repo.db.Save(course).Error
}

func (repo *courseRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Course{ID: id}).Error
}

func (repo *courseRepository) FindAll() ([]*models.Course, error) {
	var courses []*models.Course
	err := repo.db.Preload("Teachers").Preload("Students").Preload("Questions").Preload("Exams").Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (repo *courseRepository) FindTeachersByCourseID(courseID uint) ([]*models.User, error) {
	var teachers []*models.User
	err := repo.db.Model(&models.Course{ID: courseID}).Association("Teachers").Find(&teachers)
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func (repo *courseRepository) FindStudentsByCourseID(courseID uint) ([]*models.User, error) {
	var students []*models.User
	err := repo.db.Model(&models.Course{ID: courseID}).Association("Students").Find(&students)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (repo *courseRepository) AddStudentToCourse(courseID uint, studentID uint) error {
	var course models.Course
	if err := repo.db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return err
	}
	var student models.User
	if err := repo.db.Where("id = ?", studentID).First(&student).Error; err != nil {
		return err
	}
	return repo.db.Model(&course).Association("Students").Append(&student)
}

func (repo *courseRepository) AddTeacherToCourse(courseID uint, teacherID uint) error {
	var course models.Course
	if err := repo.db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return err
	}
	var teacher models.User
	if err := repo.db.Where("id = ?", teacherID).First(&teacher).Error; err != nil {
		return err
	}
	return repo.db.Model(&course).Association("Teachers").Append(&teacher)
}

func (repo *courseRepository) AddQuestionToCourse(courseID uint, questionID uint) error {
	var course models.Course
	if err := repo.db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return err
	}
	var question models.Question
	if err := repo.db.Where("id = ?", questionID).First(&question).Error; err != nil {
		return err
	}
	return repo.db.Model(&course).Association("Questions").Append(&question)
}

func (repo *courseRepository) RemoveQuestionFromCourse(courseID uint, questionID uint) error {
	// 从数据库中找到对应的 Course 记录
	var course models.Course
	if err := repo.db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("failed to find Course record: %v", err)
	}

	// 从 Course 的 questions 关联中找到对应的 Question 记录
	var question models.Question
	if err := repo.db.First(&question, questionID).Error; err != nil {
		return fmt.Errorf("failed to find Question record: %v", err)
	}

	// 从 Course 的 questions 关联中删除对应的 Question 记录
	if err := repo.db.Model(&course).Association("Questions").Delete(&question); err != nil {
		return fmt.Errorf("failed to remove question from course: %v", err)
	}

	return nil
}

func (repo *courseRepository) AddExamToCourse(courseID uint, examID uint) error {
	var course models.Course
	if err := repo.db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return err
	}
	var exam models.Exam
	if err := repo.db.Where("id = ?", examID).First(&exam).Error; err != nil {
		return err
	}
	return repo.db.Model(&course).Association("Exams").Append(&exam)
}

func (repo *courseRepository) RemoveExamFromCourse(courseID uint, examID uint) error {
	// 从数据库中找到对应的 Course 记录
	var course models.Course
	if err := repo.db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("failed to find Course record: %v", err)
	}

	// 从 Course 的 exams 关联中找到对应的 Exam 记录
	var exam models.Exam
	if err := repo.db.First(&exam, examID).Error; err != nil {
		return fmt.Errorf("failed to find Exam record: %v", err)
	}

	// 从 Course 的 exams 关联中删除对应的 Exam 记录
	if err := repo.db.Model(&course).Association("Exams").Delete(&exam); err != nil {
		return fmt.Errorf("failed to remove exam from course: %v", err)
	}

	return nil
}

func (repo *courseRepository) RemoveTeacherFromCourse(courseID uint, teacherID uint) error {
	// 从数据库中找到对应的 Course 记录
	var course models.Course
	if err := repo.db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("failed to find Course record: %v", err)
	}

	// 从 Course 的 teachers 关联中找到对应的 Teacher 记录
	var teacher models.User
	if err := repo.db.First(&teacher, teacherID).Error; err != nil {
		return fmt.Errorf("failed to find Teacher record: %v", err)
	}

	// 从 Course 的 teachers 关联中删除对应的 Teacher 记录
	if err := repo.db.Model(&course).Association("Teachers").Delete(&teacher); err != nil {
		return fmt.Errorf("failed to remove teacher from course: %v", err)
	}

	return nil
}

func (repo *courseRepository) RemoveStudentFromCourse(courseID uint, studentID uint) error {
	// 从数据库中找到对应的 Course 记录
	var course models.Course
	if err := repo.db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("failed to find Course record: %v", err)
	}

	// 从 Course 的 students 关联中找到对应的 Student 记录
	var student models.User
	if err := repo.db.First(&student, studentID).Error; err != nil {
		return fmt.Errorf("failed to find Student record: %v", err)
	}

	// 从 Course 的 students 关联中删除对应的 Student 记录
	if err := repo.db.Model(&course).Association("Students").Delete(&student); err != nil {
		return fmt.Errorf("failed to remove student from course: %v", err)
	}

	return nil
}

func (repo *courseRepository) FindQuestionsByCourseID(courseID uint) ([]*models.Question, error) {
	var questions []*models.Question
	err := repo.db.Model(&models.Course{ID: courseID}).Association("Questions").Find(&questions)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *courseRepository) FindExamsByCourseID(courseID uint) ([]*models.Exam, error) {
	var exams []*models.Exam
	err := repo.db.Model(&models.Course{ID: courseID}).Association("Exams").Find(&exams)
	if err != nil {
		return nil, err
	}
	return exams, nil
}
