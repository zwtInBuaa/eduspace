package repositories

import (
	"EDU_TH_2_backend/gin/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id uint) (*models.Post, error)
	FindAll() ([]*models.Post, error)
	FindByUserID(userID uint) ([]*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (repo *postRepository) Create(post *models.Post) error {
	return repo.db.Create(post).Error
}

func (repo *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := repo.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *postRepository) FindAll() ([]*models.Post, error) {
	var posts []*models.Post
	err := repo.db.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepository) FindByUserID(userID uint) ([]*models.Post, error) {
	var posts []*models.Post
	err := repo.db.Where("user_id = ?", userID).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepository) Update(post *models.Post) error {
	return repo.db.Save(post).Error
}

func (repo *postRepository) Delete(id uint) error {
	return repo.db.Delete(&models.Post{ID: id}).Error
}

//func (repo *postRepository) Delete(id uint) error {
//	// 预加载 Comments 字段
//	var post models.Post
//	if err := repo.db.Preload("Comments").First(&post, id).Error; err != nil {
//		return err
//	}
//
//	// 删除所有相关的 Comment 记录
//	for _, comment := range post.Comments {
//		if err := repo.db.Delete(&comment).Error; err != nil {
//			return err
//		}
//	}
//
//	// 删除 Post 记录
//	if err := repo.db.Delete(&post).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
