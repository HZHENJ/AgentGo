package dao

import (
	"context"
	"gorm.io/gorm"
	"agentgo/internal/model"
)

type UserDao interface {
	CheckUserExist(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{db: db}
}

// CheckUserExist checks if a user with the given email exists in the database.
func (dao *userDao) CheckUserExist(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := dao.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserByEmail retrieves a user record from the database based on the provided email.
func (dao *userDao) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := dao.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user record in the database.
func (dao *userDao) CreateUser(ctx context.Context, user *model.User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}