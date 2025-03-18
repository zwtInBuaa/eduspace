package services

import (
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
)

type ExamService interface {
	CreateExam(exam *models.Exam) error
	GetExamByID(id uint) (*models.Exam, error)
	UpdateExam(exam *models.Exam) error
	DeleteExam(id uint) error
	AddQuestion(examID uint, questionID uint) error
	RemoveQuestion(examID uint, questionID uint) error
	GetAllExams() ([]models.Exam, error)
	FindQuestionsByExamID(examID uint) ([]*models.Question, error)
	CreateExamWithQuestions(exam *models.Exam, questions []interface{}) error
	UpdateExamWithQuestions(exam *models.Exam, questions []interface{}) error
}

type examService struct {
	examRepo repositories.ExamRepository
}

func NewExamService(examRepo repositories.ExamRepository) ExamService {
	return &examService{examRepo: examRepo}
}

func (s *examService) CreateExam(exam *models.Exam) error {
	if err := s.examRepo.Create(exam); err != nil {
		return err
	}
	return nil
}

func (s *examService) GetExamByID(id uint) (*models.Exam, error) {
	exam, err := s.examRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return exam, nil
}

func (s *examService) UpdateExam(examFrom *models.Exam) error {
	exam, err := s.examRepo.FindByID(examFrom.ID)
	if err != nil {
		return err
	}

	exam.Name = examFrom.Name
	exam.Description = examFrom.Description

	if err := s.examRepo.Update(exam); err != nil {
		return err
	}
	return nil
}

func (s *examService) DeleteExam(id uint) error {
	if err := s.examRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *examService) GetAllExams() ([]models.Exam, error) {
	exams, err := s.examRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return exams, nil
}

func (s *examService) AddQuestion(examID uint, questionID uint) error {
	if err := s.examRepo.AddQuestion(examID, questionID); err != nil {
		return err
	}
	return nil
}

func (s *examService) RemoveQuestion(examID uint, questionID uint) error {
	if err := s.examRepo.RemoveQuestion(examID, questionID); err != nil {
		return err
	}
	return nil
}

func (s *examService) FindQuestionsByExamID(examID uint) ([]*models.Question, error) {
	questions, err := s.examRepo.FindQuestionsByExamID(examID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *examService) CreateExamWithQuestions(exam *models.Exam, questions []interface{}) error {
	if err := s.examRepo.Create(exam); err != nil {
		return err
	}

	for _, q := range questions {
		if err := s.examRepo.AddQuestion(exam.ID, uint(q.(float64))); err != nil {
			return err
		}
	}

	return nil
}

func (s *examService) UpdateExamWithQuestions(examFrom *models.Exam, questions []interface{}) error {
	exam, err := s.examRepo.FindByID(examFrom.ID)
	if err != nil {
		return err
	}

	exam.Name = examFrom.Name
	exam.Description = examFrom.Description

	if err := s.examRepo.Update(exam); err != nil {
		return err
	}

	err = s.examRepo.RemoveAllQuestions(exam.ID)
	if err != nil {
		return err
	}

	for _, q := range questions {
		if err := s.examRepo.AddQuestion(exam.ID, uint(q.(float64))); err != nil {
			return err
		}
	}

	return nil
}
