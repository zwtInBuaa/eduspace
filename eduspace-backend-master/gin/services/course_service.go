package services

import (
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
)

type CourseService interface {
	CreateCourse(course *models.Course) error
	GetCourseByID(id uint) (*models.Course, error)
	UpdateCourse(course *models.Course) error
	DeleteCourse(id uint) error
	GetAllCourses() ([]*models.Course, error)
	GetTeachersByCourseID(courseID uint) ([]*models.User, error)
	GetStudentsByCourseID(courseID uint) ([]*models.User, error)
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

type courseService struct {
	courseRepo repositories.CourseRepository
}

func NewCourseService(courseRepo repositories.CourseRepository) CourseService {
	return &courseService{courseRepo: courseRepo}
}

func (s *courseService) CreateCourse(course *models.Course) error {
	if err := s.courseRepo.Create(course); err != nil {
		return err
	}
	return nil
}

func (s *courseService) GetCourseByID(id uint) (*models.Course, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s *courseService) UpdateCourse(courseFrom *models.Course) error {
	course, err := s.courseRepo.FindByID(courseFrom.ID)
	if err != nil {
		return err
	}

	course.Name = courseFrom.Name
	course.Description = courseFrom.Description

	if err := s.courseRepo.Update(course); err != nil {
		return err
	}
	return nil
}

func (s *courseService) DeleteCourse(id uint) error {
	if err := s.courseRepo.Delete(id); err != nil {
		return err
	}
	return nil
}

func (s *courseService) GetAllCourses() ([]*models.Course, error) {
	courses, err := s.courseRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s *courseService) GetTeachersByCourseID(courseID uint) ([]*models.User, error) {
	teachers, err := s.courseRepo.FindTeachersByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func (s *courseService) GetStudentsByCourseID(courseID uint) ([]*models.User, error) {
	students, err := s.courseRepo.FindStudentsByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (s *courseService) AddStudentToCourse(courseID uint, studentID uint) error {
	if err := s.courseRepo.AddStudentToCourse(courseID, studentID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) AddTeacherToCourse(courseID uint, teacherID uint) error {
	if err := s.courseRepo.AddTeacherToCourse(courseID, teacherID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) AddQuestionToCourse(courseID uint, questionID uint) error {
	if err := s.courseRepo.AddQuestionToCourse(courseID, questionID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) RemoveQuestionFromCourse(courseID uint, questionID uint) error {
	if err := s.courseRepo.RemoveQuestionFromCourse(courseID, questionID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) AddExamToCourse(courseID uint, examID uint) error {
	if err := s.courseRepo.AddExamToCourse(courseID, examID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) RemoveExamFromCourse(courseID uint, examID uint) error {
	if err := s.courseRepo.RemoveExamFromCourse(courseID, examID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) RemoveTeacherFromCourse(courseID uint, teacherID uint) error {
	if err := s.courseRepo.RemoveTeacherFromCourse(courseID, teacherID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) RemoveStudentFromCourse(courseID uint, studentID uint) error {
	if err := s.courseRepo.RemoveStudentFromCourse(courseID, studentID); err != nil {
		return err
	}
	return nil
}

func (s *courseService) FindQuestionsByCourseID(courseID uint) ([]*models.Question, error) {
	questions, err := s.courseRepo.FindQuestionsByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *courseService) FindExamsByCourseID(courseID uint) ([]*models.Exam, error) {
	exams, err := s.courseRepo.FindExamsByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return exams, nil
}
