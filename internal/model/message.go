package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SessionID uint   `gorm:"not null" json:"session_id"`
	Username  string `gorm:"type:varchar(50)" json:"username"`
	Content   string `gorm:"type:text;not null" json:"content"`
	IsUser    bool   `gorm:"not null" json:"is_user"`
}

// History 
type History struct {
	IsUser  bool   `json:"is_user"`
	Content string `json:"content"`
}