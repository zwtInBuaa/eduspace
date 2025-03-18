package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id uint) (*models.Comment, error)
	FindByPostID(postID uint) ([]models.CommentWithUser, error)
	Update(comment *models.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (repo *commentRepository) Create(comment *models.Comment) error {
	return repo.db.Create(comment).Error
}

func (repo *commentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := repo.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (repo *commentRepository) FindByPostID(postID uint) ([]models.CommentWithUser, error) {
	var comments []models.CommentWithUser
	err := repo.db.Table("comments").
		Select("comments.*, users.username").
		Where("post_id = ?", postID).
		Joins("JOIN users ON comments.user_id = users.id").
		Scan(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *commentRepository) Update(comment *models.Comment) error {
	return repo.db.Save(comment).Error
}

func (repo *commentRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Comment{ID: id}).Error
}
