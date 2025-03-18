package models

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserSubmitHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	QuestionID uint      `gorm:"not null" json:"question_id"`
	Tags       []string  `json:"tags" gorm:"-"`
	TagsStr    string    `json:"-" gorm:"column:tags"`
	Answer     string    `gorm:"not null" json:"answer"`
	IsCorrect  bool      `gorm:"not null" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User     User     `gorm:"foreignKey:UserID" json:"-"`
	Question Question `gorm:"foreignKey:QuestionID" json:"-"`
}

// 在创建表时为 Question 模型添加回调函数
func (q *UserSubmitHistory) BeforeSave(tx *gorm.DB) error {
	//if len(q.Tags) > 0 {
	tagsBytes, err := json.Marshal(q.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v", err)
	}
	q.TagsStr = string(tagsBytes)
	//}
	return nil
}

// 在查询数据时将 Tags 反序列化为字符串切片
func (q *UserSubmitHistory) AfterFind(tx *gorm.DB) error {
	//if q.TagsStr != "" {
	err := json.Unmarshal([]byte(q.TagsStr), &q.Tags)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tags: %v", err)
	}
	//}
	return nil
}

func (c *UserSubmitHistory) MarshalJSON() ([]byte, error) {
	type Alias UserSubmitHistory
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
