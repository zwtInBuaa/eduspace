package services

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
	"fmt"
)

type QuestionService interface {
	CreateQuestion(question *models.Question) error
	CreateQuestionAndAddSource(courseID uint, question *models.Question) error
	GetQuestionByID(id uint) (*models.Question, error)
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(id uint) error
	GetAllQuestions() ([]models.Question, error)
	Submit(userId uint, quesitionId uint, submitForm *models.SubmitFrom) (map[string]interface{}, error)
}

type questionService struct {
	questionRepo repositories.QuestionRepository
}

func NewQuestionService(questionRepo repositories.QuestionRepository) QuestionService {
	return &questionService{questionRepo: questionRepo}
}

func (s *questionService) CreateQuestion(question *models.Question) error {
	if err := s.questionRepo.Create(question); err != nil {
		return err
	}
	return nil
}

func (s *questionService) CreateQuestionAndAddSource(courseID uint, question *models.Question) error {
	if err := s.questionRepo.Create(question); err != nil {
		return err
	}
	if err := s.questionRepo.AddQuestionToCourse(courseID, question); err != nil {
		_ = s.questionRepo.Delete(question.ID)
		return err
	}
	return nil
}

func (s *questionService) GetQuestionByID(id uint) (*models.Question, error) {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionService) UpdateQuestion(questionFrom *models.Question) error {
	question, err := s.questionRepo.FindByID(questionFrom.ID)
	if err != nil {
		return err
	}

	// 更新字段
	question.Title = questionFrom.Title
	question.Content = questionFrom.Content
	question.Answer = questionFrom.Answer
	question.Tags = questionFrom.Tags
	question.Type = questionFrom.Type
	question.Data = questionFrom.Data

	if err := s.questionRepo.Update(question); err != nil {
		return err
	}

	return nil
}

func (s *questionService) DeleteQuestion(id uint) error {
	if err := s.questionRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *questionService) GetAllQuestions() ([]models.Question, error) {
	questions, err := s.questionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return questions, nil
}

// 根据question_id获取question,把submitForm中的数据和question中的数据进行正则匹配
func (s *questionService) Submit(userId uint, quesition_id uint, submitForm *models.SubmitFrom) (map[string]interface{}, error) {
	question, err := s.questionRepo.FindByID(quesition_id)
	if err != nil {
		return nil, err
	}
	right_answer := question.Answer
	submit_answer := submitForm.Answer
	// right_answer 存的是正则表达式字符串，submit_answer是字符串
	// 正则表达式匹配

	submitInfo := map[string]interface{}{}

	//if ok, _ := regexp.MatchString(right_answer, submit_answer); ok {
	if right_answer == submit_answer {
		// 匹配成功
		submitInfo["result"] = true
	} else {
		// 匹配失败
		submitInfo["result"] = false
	}

	submitInfo["result_desc"] = ""

	// 创建UserSubmitHistory模型并保存到数据库中
	submitHistory := &models.UserSubmitHistory{
		UserID:     userId,
		QuestionID: quesition_id,
		Tags:       question.Tags,
		Answer:     submit_answer,
		IsCorrect:  submitInfo["result"].(bool),
	}
	if err := s.questionRepo.CreateSubmitHistory(submitHistory); err != nil {
		return nil, err
	}

	// logs: 哪个用户提交了什么题目，提交的结果是什么
	logger.Info(fmt.Sprintf("user_id: %d, question_id: %d, submit_answer: %s, result: %v", userId, quesition_id, submit_answer, submitInfo["result"]))

	return submitInfo, nil
}
