package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Exam struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	Name        string     `json:"name" gorm:"column:name"`
	Description string     `json:"description" gorm:"column:description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Questions   []Question `json:"-" gorm:"many2many:exam_questions;"`
}

func (c *Exam) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Model(c).Association("Questions").Clear(); err != nil {
		return err
	}

	tx.Table("course_exams").Where("exam_id = ?", c.ID).Delete(nil)
	return nil
}

func (c *Exam) MarshalJSON() ([]byte, error) {
	type Alias Exam
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
