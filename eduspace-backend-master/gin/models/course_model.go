package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Teachers    []User     `json:"-" gorm:"many2many:course_teachers;"`
	Students    []User     `json:"-" gorm:"many2many:course_students;"`
	Exams       []Exam     `json:"-" gorm:"many2many:course_exams;"`
	Questions   []Question `json:"-" gorm:"many2many:course_questions;"`
}

func (c *Course) BeforeDelete(tx *gorm.DB) error {
	// Clear Teachers many2many association
	if err := tx.Model(c).Association("Teachers").Clear(); err != nil {
		return err
	}
	// Clear Students many2many association
	if err := tx.Model(c).Association("Students").Clear(); err != nil {
		return err
	}
	// Clear Exams many2many association
	if err := tx.Model(c).Association("Exams").Clear(); err != nil {
		return err
	}
	// Clear Questions many2many association
	if err := tx.Model(c).Association("Questions").Clear(); err != nil {
		return err
	}
	return nil
}

func (c *Course) MarshalJSON() ([]byte, error) {
	type Alias Course
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

// 另一种实现方法
//type CustomTime struct {
//	time.Time
//}
//
//func (ct *CustomTime) MarshalJSON() ([]byte, error) {
//	formatted := fmt.Sprintf("\"%s\"", ct.Format("2006-01-02 15:04"))
//	return []byte(formatted), nil
//}
//
//func (ct *CustomTime) Value() (driver.Value, error) {
//	return ct.Time, nil
//}
//
//func (ct *CustomTime) Scan(value interface{}) error {
//	if value == nil {
//		ct.Time = time.Time{}
//		return nil
//	}
//	t, ok := value.(time.Time)
//	if !ok {
//		return fmt.Errorf("invalid type for CustomTime: %T", value)
//	}
//	ct.Time = t
//	return nil
//}
//
//type Course struct {
//	ID          uint       `json:"id" gorm:"primary_key"`
//	Name        string     `json:"name"`
//	Description string     `json:"description"`
//	CreatedAt   CustomTime `json:"created_at" gorm:"type:timestamp"`
//	UpdatedAt   CustomTime `json:"updated_at"`
//	Teachers    []*User    `json:"teachers" gorm:"many2many:course_teachers;"`
//	Students    []*User    `json:"students" gorm:"many2many:course_students;"`
//}
