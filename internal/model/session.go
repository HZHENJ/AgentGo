package model

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null" json:"username"`
	Title   string `gorm:"type:varchar(100);not null" json:"title"`
}

type SessionInfo struct {
	SessionID uint   	`json:"session_id"`
	Title	 string 	`json:"title"`
}

