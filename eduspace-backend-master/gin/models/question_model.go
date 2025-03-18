package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Question struct {
	ID        uint                   `json:"id" gorm:"primary_key"`
	Title     string                 `json:"title"` // 习题名
	Content   string                 `json:"content" gorm:"column:content"`
	Answer    string                 `json:"answer" gorm:"column:answer"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data" gorm:"-"`
	DataStr   string                 `json:"-" gorm:"column:data"`
	Tags      []string               `json:"tags" gorm:"-"`
	TagsStr   string                 `json:"-" gorm:"column:tags"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// 在创建表时为 Question 模型添加回调函数
func (q *Question) BeforeSave(tx *gorm.DB) error {
	//if len(q.Tags) > 0 {
	tagsBytes, err := json.Marshal(q.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %v", err)
	}
	q.TagsStr = string(tagsBytes)
	//}
	//if len(q.Data) > 0 {
	dataBytes, err := json.Marshal(q.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}
	q.DataStr = string(dataBytes)
	//}
	return nil
}

func (c *Question) BeforeDelete(tx *gorm.DB) error {
	tx.Table("exam_questions").Where("question_id = ?", c.ID).Delete(nil)
	tx.Table("course_questions").Where("question_id = ?", c.ID).Delete(nil)

	// 删除所有user_submit_history记录
	if err := tx.Delete(&UserSubmitHistory{}, "question_id = ?", c.ID).Error; err != nil {
		return err
	}
	return nil
}

// 在查询数据时将 Tags 反序列化为字符串切片
func (q *Question) AfterFind(tx *gorm.DB) error {
	//if q.TagsStr != "" {
	err := json.Unmarshal([]byte(q.TagsStr), &q.Tags)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tags: %v", err)
	}
	//}
	//if q.DataStr != "" {
	var data map[string]interface{}
	err = json.Unmarshal([]byte(q.DataStr), &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %v", err)
	}
	q.Data = data
	//}
	return nil
}

func (c *Question) MarshalJSON() ([]byte, error) {
	type Alias Question
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

// SubmitFrom
// 提交答案表单
// 参数： answer
type SubmitFrom struct {
	Answer string `json:"answer" binding:"required"`
}

func (r *SubmitFrom) UnmarshalJSON(data []byte) error {
	required := struct {
		Answer string `json:"answer"`
	}{}

	err := json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Answer == "" {
		err = errors.New("缺少必填字段answer")
	} else {
		r.Answer = required.Answer
	}

	return nil
}
