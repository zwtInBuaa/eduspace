package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	Content   string    `json:"content" gorm:"size:5120"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentWithUser struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	Avatar    string    `json:"user_avatar"`
	Username  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) BeforeSave(tx *gorm.DB) error {
	maxContentLength := 5120
	if len(c.Content) > maxContentLength {
		return fmt.Errorf("帖子内容过长，最大长度为 %d 字节", maxContentLength)
	}
	return nil
}

func (c *Comment) MarshalJSON() ([]byte, error) {
	type Alias Comment
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04"),
		Alias:     (*Alias)(c),
	})
}

func (c *CommentWithUser) MarshalJSON() ([]byte, error) {
	type Alias CommentWithUser
	return json.Marshal(&struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04"),
		Alias:     (*Alias)(c),
	})
}
