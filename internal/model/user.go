package model

import (
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username 	string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password 	string `gorm:"type:varchar(255);not null" json:"-"`
	Nickname 	string `gorm:"type:varchar(50);not null" json:"nickname"`
	Email 		string `gorm:"type:varchar(100);unique;not null;index" json:"email"`
}

// SetPassword password hash and set to user
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword validate password
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}