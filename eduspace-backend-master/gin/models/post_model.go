package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
	//CourseID  uint      `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Comments []Comment `gorm:"foreignKey:PostID" json:"comments"`
}

func (p *Post) BeforeSave(tx *gorm.DB) error {
	maxTitleLength := 200
	maxContentLength := 5120
	if len(p.Title) > maxTitleLength {
		return fmt.Errorf("帖子标题过长，最大长度为 %d 字节", maxTitleLength)
	}
	if len(p.Content) > maxContentLength {
		return fmt.Errorf("帖子内容过长，最大长度为 %d 字节", maxContentLength)
	}
	return nil
}

func (c *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
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

func (post *Post) BeforeDelete(tx *gorm.DB) error {
	// 删除所有相关的 Comment 记录
	if err := tx.Delete(&Comment{}, "post_id = ?", post.ID).Error; err != nil {
		return err
	}

	return nil
}
