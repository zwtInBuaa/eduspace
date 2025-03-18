package services

import (
	"EDU_TH_2_backend/gin/logger"
	"EDU_TH_2_backend/gin/models"
	"EDU_TH_2_backend/gin/repositories"
	"errors"
	"fmt"
	"unicode/utf8"
)

type PostService interface {
	CreatePost(title, content string, userID uint) error
	GetPostByID(id uint) (map[string]interface{}, error)
	GetAllPosts() ([]map[string]interface{}, error)
	GetPostsByUserID(userID uint) ([]*models.Post, error)
	UpdatePost(id uint, title, content string, changeUserId, changeUserRole int) error
	DeletePost(id uint, changeUserId, changeUserRole int) error
}

type postService struct {
	postRepo repositories.PostRepository
	userRepo repositories.UserRepository
}

func NewPostService(postRepo repositories.PostRepository, userRepo repositories.UserRepository) PostService {
	return &postService{postRepo: postRepo, userRepo: userRepo}
}

func (s *postService) CreatePost(title, content string, userID uint) error {
	post := &models.Post{Title: title, Content: content, UserID: userID}
	if err := s.postRepo.Create(post); err != nil {
		return err
	}

	// log: 哪个用户创建了哪篇文章
	logger.Info(fmt.Sprintf("user %d created post %d", userID, post.ID))

	return nil
}

func (s *postService) GetPostByID(id uint) (map[string]interface{}, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	userId := post.UserID
	user, err := s.userRepo.FindByID(userId)
	if err != nil {
		return nil, err
	}

	postData := map[string]interface{}{
		"id":          post.ID,
		"title":       post.Title,
		"content":     post.Content,
		"created_at":  post.CreatedAt.Format("2006-01-02 15:04"),
		"updated_at":  post.UpdatedAt.Format("2006-01-02 15:04"),
		"poster_id":   user.ID,
		"poster_name": user.Username,
	}

	return postData, nil
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	// 将字符串转换为rune切片
	runes := []rune(s)

	// 确保我们不截断一个中文字符或emoji
	if len(runes) > maxLength {
		runes = runes[0:maxLength]
		for i := len(runes) - 1; i >= 0; i-- {
			if utf8.RuneStart(byte(runes[i])) {
				runes = runes[0:i]
				break
			}
		}

		// 在字符串末尾添加省略号
		runes = append(runes, []rune("...")...)
	}

	// 将rune切片转换为字符串
	return string(runes)
}

func (s *postService) GetAllPosts() ([]map[string]interface{}, error) {
	posts, err := s.postRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var postsData []map[string]interface{}

	for _, post := range posts {
		userId := post.UserID
		user, err := s.userRepo.FindByID(userId)
		if err != nil {
			return nil, err
		}

		content := post.Content
		infro := truncateString(content, 50)

		newPost := map[string]interface{}{
			"id":          post.ID,
			"title":       post.Title,
			"created_at":  post.CreatedAt.Format("2006-01-02 15:04"),
			"updated_at":  post.UpdatedAt.Format("2006-01-02 15:04"),
			"poster_name": user.Username,
			"poster_id":   user.ID,
			"intro":       infro,
		}

		postsData = append(postsData, newPost)
	}

	return postsData, nil
}

func (s *postService) GetPostsByUserID(userID uint) ([]*models.Post, error) {
	posts, err := s.postRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) UpdatePost(id uint, title, content string, changeUserId, changeUserRole int) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if changeUserRole == 3 {
		if changeUserId != int(post.UserID) {
			return errors.New("你没有权限修改该文章")
		}
	}

	post.Title = title
	post.Content = content
	if err := s.postRepo.Update(post); err != nil {
		return err
	}

	// log: 哪个用户更新了哪篇文章
	logger.Info(fmt.Sprintf("user %d updated post %d", post.UserID, post.ID))

	return nil
}

func (s *postService) DeletePost(id uint, changeUserId, changeUserRole int) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return err
	}

	if changeUserRole == 3 {
		if changeUserId != int(post.UserID) {
			return errors.New("你没有权限删除该文章")
		}
	}

	if err := s.postRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
