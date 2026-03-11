package dao

import (
	"context"
	"agentgo/internal/model"

	"gorm.io/gorm"
)

type SessionDao interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionsByUsername(ctx context.Context, username string) ([]*model.Session, error)
	GetSessionByID(ctx context.Context, id uint) (*model.Session, error)
	DeleteSessionByID(ctx context.Context, id uint, username string) error
}

type sessionDao struct {
	db *gorm.DB
}

func NewSessionDao(db *gorm.DB) SessionDao {
	return &sessionDao{db: db}
}

// CreateSession creates a new session record in the database.
func(dao *sessionDao) CreateSession(ctx context.Context, session *model.Session) error {
	return dao.db.WithContext(ctx).Create(session).Error
}

// GetSessionByUsername retrieves all session records from the database based on the provided username.
func(dao *sessionDao) GetSessionsByUsername(ctx context.Context, username string) ([]*model.Session, error) {
	var sessions []*model.Session
	err := dao.db.WithContext(ctx).
		Where("username = ?", username).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

// GetSessionByID retrieves a session record from the database based on the provided session ID.
func (dao *sessionDao) GetSessionByID(ctx context.Context, id uint) (*model.Session, error) {
	var session model.Session
	
	err := dao.db.WithContext(ctx).
		Where("id = ?", id).
		First(&session).Error
		
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteSessionByID deletes a session record from the database based on the provided session ID and username.
func (dao *sessionDao) DeleteSessionByID(ctx context.Context, id uint, username string) error {
	return dao.db.WithContext(ctx).
		Where("id = ? AND username = ?", id, username).
		Delete(&model.Session{}).Error
}