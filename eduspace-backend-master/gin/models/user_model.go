package models

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"column:username"`
	BuaaId    string    `json:"buaa_id" gorm:"column:buaa_id;unique"`
	Password  string    `json:"-"` // 密码在JSON序列化时不显示
	Role      int64     `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Posts    []Post    `json:"-" gorm:"foreignKey:UserID"`
	Comments []Comment `json:"-" gorm:"foreignKey:UserID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) BeforeDelete(tx *gorm.DB) error {
	// 删除所有相关的 Comment 记录
	if err := tx.Delete(&Comment{}, "user_id = ?", user.ID).Error; err != nil {
		return err
	}

	// 删除所有相关的 Post 记录
	if err := tx.Delete(&Post{}, "user_id = ?", user.ID).Error; err != nil {
		return err
	}

	// 删除所有user_submit_history记录
	if err := tx.Delete(&UserSubmitHistory{}, "user_id = ?", user.ID).Error; err != nil {
		return err
	}

	// Clear Courses many2many association
	tx.Table("course_teachers").Where("user_id = ?", user.ID).Delete(nil)
	tx.Table("course_students").Where("user_id = ?", user.ID).Delete(nil)

	return nil
}

func (c *User) MarshalJSON() ([]byte, error) {
	type Alias User
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

// ResetPasswordForm
// 重置密码表单
// 参数：user_id，newPassword
type ResetPasswordForm struct {
	BuaaID      string `json:"buaa_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (r *ResetPasswordForm) UnmarshalJSON(data []byte) error {
	required := struct {
		BuaaID      string `json:"buaa_id"`
		NewPassword string `json:"password"`
	}{}
	err := json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if len(required.BuaaID) == 0 {
		err = errors.New("缺少必填字段buaa_id")
	} else if len(required.NewPassword) == 0 {
		err = errors.New("缺少必填字段password")
	} else {
		r.BuaaID = required.BuaaID
		r.NewPassword = required.NewPassword
	}
	return nil
}
