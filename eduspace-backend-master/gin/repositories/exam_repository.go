package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"gorm.io/gorm"
)

type ExamRepository interface {
	Create(exam *models.Exam) error
	FindByID(id uint) (*models.Exam, error)
	Update(exam *models.Exam) error
	Delete(id uint) error
	FindAll() ([]models.Exam, error)
	AddQuestion(examID uint, questionID uint) error
	RemoveQuestion(examID uint, questionID uint) error
	RemoveAllQuestions(examID uint) error
	FindQuestionsByExamID(examID uint) ([]*models.Question, error)
}

type examRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) ExamRepository {
	return &examRepository{db: db}
}

func (repo *examRepository) Create(exam *models.Exam) error {
	return repo.db.Create(exam).Error
}

func (repo *examRepository) FindByID(id uint) (*models.Exam, error) {
	var exam models.Exam
	err := repo.db.Where("id = ?", id).First(&exam).Error
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func (repo *examRepository) Update(exam *models.Exam) error {
	return repo.db.Save(exam).Error
}

func (repo *examRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Exam{ID: id}).Error
}

func (repo *examRepository) FindAll() ([]models.Exam, error) {
	var exams []models.Exam
	err := repo.db.Find(&exams).Error
	if err != nil {
		return nil, err
	}
	return exams, nil
}

func (repo *examRepository) AddQuestion(examID uint, questionID uint) error {
	return repo.db.Model(&models.Exam{ID: examID}).Association("Questions").Append(&models.Question{ID: questionID})
}

func (repo *examRepository) RemoveQuestion(examID uint, questionID uint) error {
	return repo.db.Model(&models.Exam{ID: examID}).Association("Questions").Delete(&models.Question{ID: questionID})
}

func (repo *examRepository) RemoveAllQuestions(examID uint) error {
	exam, err := repo.FindByID(examID)
	if err != nil {
		return err
	}
	return repo.db.Model(exam).Association("Questions").Clear()
}

func (repo *examRepository) FindQuestionsByExamID(examID uint) ([]*models.Question, error) {
	var questions []*models.Question
	err := repo.db.Model(&models.Exam{ID: examID}).Association("Questions").Find(&questions)
	if err != nil {
		return nil, err
	}
	return questions, nil
}
