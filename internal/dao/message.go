package dao

import (
	"context"
	"agentgo/internal/model"
	"gorm.io/gorm"
)

type MessageDao interface {
	CreateMessage(ctx context.Context, msg *model.Message) error
	GetHistoryBySessionID(ctx context.Context, sessionID uint) ([]*model.Message, error)
}

type messageDao struct {
	db *gorm.DB
}

func NewMessageDao(db *gorm.DB) MessageDao {
	return &messageDao {db: db}
}

func (d *messageDao) CreateMessage(ctx context.Context, msg *model.Message) error {
	return d.db.WithContext(ctx).Create(msg).Error
}

func (d *messageDao) GetHistoryBySessionID(ctx context.Context, sessionID uint) ([]*model.Message, error) {
	var msgs []*model.Message
	err := d.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&msgs).Error
	return msgs, err
}