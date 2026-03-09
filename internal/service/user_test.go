package service

import (
    "context"
    "errors"
    "testing"

    "agentgo/internal/cache"
    "agentgo/internal/dao"
    "agentgo/internal/model"
    "agentgo/internal/types"
    "agentgo/pkg/e"
)

// --- Mocks ---
type mockUserDao struct {
    byEmail map[string]*model.User
}

func newMockUserDao() *mockUserDao { return &mockUserDao{byEmail: map[string]*model.User{}} }

func (m *mockUserDao) CheckUserExist(ctx context.Context, email string) (bool, error) {
    _, ok := m.byEmail[email]
    return ok, nil
}
func (m *mockUserDao) CreateUser(ctx context.Context, user *model.User) error {
    if _, ok := m.byEmail[user.Email]; ok {
        return errors.New("exists")
    }
    // emulate auto-increment id
    user.ID = uint(len(m.byEmail) + 1)
    cp := *user
    m.byEmail[user.Email] = &cp
    return nil
}
func (m *mockUserDao) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    u, ok := m.byEmail[email]
    if !ok {
        return nil, errors.New("not found")
    }
    cp := *u
    return &cp, nil
}

type mockUserCache struct {
    store map[string]string
}

func newMockUserCache() *mockUserCache { return &mockUserCache{store: map[string]string{}} }

func (m *mockUserCache) SetCaptchaForEmail(ctx context.Context, email, captcha string) error {
    m.store[email] = captcha
    return nil
}
func (m *mockUserCache) CheckCaptchaForEmail(ctx context.Context, email, captcha string) (bool, error) {
    v, ok := m.store[email]
    if !ok {
        return false, nil
    }
    if v == captcha {
        delete(m.store, email)
        return true, nil
    }
    return false, nil
}

// Ensure mocks satisfy interfaces
var _ dao.UserDao = (*mockUserDao)(nil)
var _ cache.UserCacheDao = (*mockUserCache)(nil)

func TestRegister_Success(t *testing.T) {
    ud := newMockUserDao()
    cd := newMockUserCache()
    svc := NewUserService(ud, cd)

    // prepare captcha
    _ = cd.SetCaptchaForEmail(context.Background(), "u@example.com", "123456")

    req := &types.UserRegisterRequest{Email: "u@example.com", Captcha: "123456", Password: "abc123"}
    data, code := svc.Register(context.Background(), req)
    if code != e.SUCCESS {
        t.Fatalf("expected SUCCESS, got %d", code)
    }
    if data == nil {
        t.Fatalf("expected response data, got nil")
    }
    // ensure user created and password hashed
    u, _ := ud.GetUserByEmail(context.Background(), "u@example.com")
    if u.Password == "abc123" || u.Password == "" {
        t.Fatalf("password should be hashed")
    }
}

func TestRegister_ExistingUser(t *testing.T) {
    ud := newMockUserDao()
    cd := newMockUserCache()
    svc := NewUserService(ud, cd)

    // seed existing user
    u := &model.User{Email: "u@example.com", Username: "u", Nickname: "u"}
    _ = u.SetPassword("any")
    _ = ud.CreateUser(context.Background(), u)

    req := &types.UserRegisterRequest{Email: "u@example.com", Captcha: "x", Password: "abc123"}
    _, code := svc.Register(context.Background(), req)
    if code != e.ERROR_USER_EXIST {
        t.Fatalf("expected ERROR_USER_EXIST, got %d", code)
    }
}

func TestRegister_InvalidCaptcha(t *testing.T) {
    ud := newMockUserDao()
    cd := newMockUserCache()
    svc := NewUserService(ud, cd)

    // set different captcha
    _ = cd.SetCaptchaForEmail(context.Background(), "u@example.com", "654321")

    req := &types.UserRegisterRequest{Email: "u@example.com", Captcha: "123456", Password: "abc123"}
    _, code := svc.Register(context.Background(), req)
    if code != e.ERROR_INVALID_CAPTCHA {
        t.Fatalf("expected ERROR_INVALID_CAPTCHA, got %d", code)
    }
}

func TestLogin_Scenarios(t *testing.T) {
    ud := newMockUserDao()
    cd := newMockUserCache()
    svc := NewUserService(ud, cd)

    // seed a user
    u := &model.User{Email: "u@example.com", Username: "u", Nickname: "u"}
    _ = u.SetPassword("abc123")
    _ = ud.CreateUser(context.Background(), u)

    // ok
    if _, code := svc.Login(context.Background(), &types.UserLoginRequest{Email: "u@example.com", Password: "abc123"}); code != e.SUCCESS {
        t.Fatalf("login ok expected SUCCESS, got %d", code)
    }
    // wrong pwd
    if _, code := svc.Login(context.Background(), &types.UserLoginRequest{Email: "u@example.com", Password: "bad"}); code != e.ERROR_USER_WRONG_PWD {
        t.Fatalf("wrong pwd expected ERROR_USER_WRONG_PWD, got %d", code)
    }
    // not exist
    if _, code := svc.Login(context.Background(), &types.UserLoginRequest{Email: "none@example.com", Password: "x"}); code != e.ERROR_USER_NOT_EXIST {
        t.Fatalf("not exist expected ERROR_USER_NOT_EXIST, got %d", code)
    }
}
