package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(question *models.Question) error
	FindByID(id uint) (*models.Question, error)
	Update(question *models.Question) error
	Delete(id uint) error
	FindAll() ([]models.Question, error)
	AddQuestionToCourse(courseID uint, question *models.Question) error
	CreateSubmitHistory(history *models.UserSubmitHistory) error
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (repo *questionRepository) Create(question *models.Question) error {
	return repo.db.Create(question).Error
}

func (repo *questionRepository) FindByID(id uint) (*models.Question, error) {
	var question models.Question
	err := repo.db.Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (repo *questionRepository) Update(question *models.Question) error {
	return repo.db.Save(question).Error
}

func (repo *questionRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Question{ID: id}).Error
}

func (repo *questionRepository) FindAll() ([]models.Question, error) {
	var questions []models.Question
	err := repo.db.Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *questionRepository) AddQuestionToCourse(courseID uint, question *models.Question) error {
	// 首先，根据课程ID查询课程
	var course models.Course
	err := repo.db.First(&course, courseID).Error
	if err != nil {
		return err
	}

	// 然后，将题目添加到课程的题目列表中
	course.Questions = append(course.Questions, *question)

	// 最后，保存更新后的课程
	return repo.db.Save(&course).Error
}

func (repo *questionRepository) CreateSubmitHistory(history *models.UserSubmitHistory) error {
	return repo.db.Create(history).Error
}
