package service

import (
	"time"
	"agentgo/internal/dao"
	"agentgo/internal/cache"
	"agentgo/internal/model"
	"agentgo/internal/types"
	"agentgo/pkg/e"
	"agentgo/pkg/utils"
	"context"
)

type UserService struct {
	userDao dao.UserDao
	userCacheDao cache.UserCacheDao
}

func NewUserService(userDao dao.UserDao, userCacheDao cache.UserCacheDao) *UserService {
	return &UserService{
		userDao: userDao,
		userCacheDao: userCacheDao,
	}
}

// Register handles user registration. It checks if the user already exists, generates a random username, and creates a new user record in the database.
func (s *UserService) Register(ctx context.Context, req *types.UserRegisterRequest) (interface{}, int) {
	// Check if the user already exists
	exists, err := s.userDao.CheckUserExist(ctx, req.Email)
	if err != nil {
		return nil, e.ERROR
	}
	if exists {
		return nil, e.ERROR_USER_EXIST
	}

	// Check the validity of the captcha
	isValid, err := s.userCacheDao.CheckCaptchaForEmail(ctx, req.Email, req.Captcha)
	if err != nil {
		return nil, e.ERROR
	}
	if !isValid {
		return nil, e.ERROR_INVALID_CAPTCHA
	}

	randomUsername := utils.GenerateDefaultUsername()

	newUser := &model.User {
		Email: req.Email,
		Username: randomUsername,
		Nickname: randomUsername,
	}

	if err := newUser.SetPassword(req.Password); err != nil {
		return nil, e.ERROR
	}

	// TODO: send welcome email to user

	// Create the new user in the database
	if err := s.userDao.CreateUser(ctx, newUser); err != nil {
		return nil, e.ERROR
	}

	// Generate JWT token for the new user
	token, err := utils.GenerateToken(newUser.ID, newUser.Email, newUser.Username)
	if err != nil {
		return nil, e.ERROR
	}

	resp := &types.UserRegisterResponse{
		Token: token,
		User: types.UserInfo{
			Id:       newUser.ID,
			Email:    newUser.Email,
			Username: newUser.Username,
			Nickname: newUser.Nickname,
		},
	}

	return resp, e.SUCCESS
}

// Login handles user login. It checks if the user exists, verifies the password, and generates a JWT token for the user.
func (s *UserService) Login(ctx context.Context, req *types.UserLoginRequest) (interface{}, int) {
	// Check User existence
	user, err := s.userDao.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, e.ERROR_USER_NOT_EXIST
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, e.ERROR_USER_WRONG_PWD
	}

	// Generate JWT token for the user
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, e.ERROR
	}

	resp := &types.UserLoginResponse{
		Token: token,
		User: types.UserInfo{
			Id: user.ID,
			Email: user.Email,
			Username: user.Username,
			Nickname: user.Nickname,
		},	
	}

	return resp, e.SUCCESS
}

// Logout handles user logout. It parses the JWT token to get the claims and stores the invalidated token in Redis with an expiration time equal to the remaining validity period of the token.
func (s *UserService) Logout(ctx context.Context, token string) (interface{}, int) {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, e.SUCCESS
	}

	now := time.Now().Unix()
	exp := claims.ExpiresAt.Unix()
	remain := exp - now

	if remain <= 0 {
		return nil, e.SUCCESS
	}

	// TODO : store the invalidated token in Redis with an expiration time equal to the remaining validity period of the token
	
	return nil, e.SUCCESS
}

// SendCaptcha generates a captcha code for the given email and stores it in Redis 
// with an expiration time. It returns an error code if the operation fails.
func (s *UserService) SendCaptcha(ctx context.Context, email string) (interface{}, int) {
	// 1. Generate a random captcha code
	captcha := utils.GenerateRandomString(6)
	if err := s.userCacheDao.SetCaptchaForEmail(ctx, email, captcha); err != nil {
		return nil, e.ERROR
	}	

	// 2. Send the captcha code to the user's email address using an email service
	if err := utils.SendEmail(email, "Your Captcha Code", "Your captcha code is: " + captcha); err != nil {
		return nil, e.ERROR_SEND_EMAIL
	}

	return nil, e.SUCCESS
}