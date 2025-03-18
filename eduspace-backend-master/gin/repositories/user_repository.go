package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByBuaaId(buaaId string) (*models.User, error)
	FindAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetTeacherCourses(teacherID uint) ([]*models.Course, error)
	GetStudentCourses(studentID uint) ([]*models.Course, error)
	GetUserSubmitHistory(userID uint) ([]*models.UserSubmitHistory, error)
	GetQuestionByTag(tag string) ([]*models.Question, error)
	GetAllQuestionID() ([]uint, error)
	GetAllQuestion() ([]*models.Question, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Create(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := repo.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByBuaaId(buaaId string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("buaa_id = ?", buaaId).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("buaa_id not found")
	}
	return &user, nil
}

func (repo *userRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	//err := repo.db.Preload("Posts").Preload("Comments").Find(&users).Error  // 预加载
	err := repo.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) Update(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *userRepository) Delete(id uint) error {
	return repo.db.Delete(&models.User{ID: id}).Error
}

func (r *userRepository) GetTeacherCourses(teacherID uint) ([]*models.Course, error) {
	var courses []*models.Course
	err := r.db.Table("courses").Joins("JOIN course_teachers ON courses.id = course_teachers.course_id").
		Where("course_teachers.user_id = ?", teacherID).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *userRepository) GetStudentCourses(studentID uint) ([]*models.Course, error) {
	var courses []*models.Course
	err := r.db.Table("courses").Joins("JOIN course_students ON courses.id = course_students.course_id").
		Where("course_students.user_id = ?", studentID).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (repo *userRepository) GetUserSubmitHistory(userID uint) ([]*models.UserSubmitHistory, error) {
	var history []*models.UserSubmitHistory
	err := repo.db.Where("user_id = ?", userID).Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (repo *userRepository) GetQuestionByTag(tag string) ([]*models.Question, error) {
	var questions []*models.Question
	err := repo.db.Where("tags LIKE ?", "%"+tag+"%").Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *userRepository) GetAllQuestionID() ([]uint, error) {
	var questions []*models.Question
	var questionIDs []uint
	err := repo.db.Find(&questions).Error
	if err != nil {
		return nil, err
	}
	for _, question := range questions {
		questionIDs = append(questionIDs, question.ID)
	}
	return questionIDs, nil
}

func (repo *userRepository) GetAllQuestion() ([]*models.Question, error) {
	var questions []*models.Question
	err := repo.db.Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}
