package services

import (
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
	"errors"
)

type CommentService interface {
	CreateComment(content string, userID, postID uint) error
	GetCommentByID(id uint) (*models.Comment, error)
	GetCommentsByPostID(postID uint) ([]models.CommentWithUser, error)
	UpdateComment(id uint, content string, changeUserId, changeUserRole int) error
	DeleteComment(id uint, changeUserId, changeUserRole int) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
}

func NewCommentService(commentRepo repositories.CommentRepository) CommentService {
	return &commentService{commentRepo: commentRepo}
}

func (s *commentService) CreateComment(content string, userID, postID uint) error {
	comment := &models.Comment{Content: content, UserID: userID, PostID: postID}
	if err := s.commentRepo.Create(comment); err != nil {
		return err
	}
	return nil
}

func (s *commentService) GetCommentByID(id uint) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentService) GetCommentsByPostID(postID uint) ([]models.CommentWithUser, error) {
	commentsWithUser, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}

	return commentsWithUser, nil
}

func (s *commentService) UpdateComment(id uint, content string, changeUserId, changeUserRole int) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return err
	}

	if changeUserRole == 3 {
		if comment.UserID != uint(changeUserId) {
			return errors.New("你沒有权限修改此评论")
		}
	}

	comment.Content = content
	if err := s.commentRepo.Update(comment); err != nil {
		return err
	}
	return nil
}

func (s *commentService) DeleteComment(id uint, changeUserId, changeUserRole int) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return err
	}

	if changeUserRole == 3 {
		if comment.UserID != uint(changeUserId) {
			return errors.New("你沒有权限删除此评论")
		}
	}

	if err := s.commentRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
