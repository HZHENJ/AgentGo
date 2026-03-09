package dao

import (
    "context"
    "testing"

    "agentgo/internal/model"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open sqlite: %v", err)
    }
    if err := db.AutoMigrate(&model.User{}); err != nil {
        t.Fatalf("failed to migrate: %v", err)
    }
    return db
}

func TestUserDao_CreateAndGet(t *testing.T) {
    db := setupTestDB(t)
    d := NewUserDao(db)

    u := &model.User{Email: "test@example.com", Username: "user1", Nickname: "user1"}
    if err := u.SetPassword("pass123"); err != nil {
        t.Fatalf("set password error: %v", err)
    }

    if err := d.CreateUser(context.Background(), u); err != nil {
        t.Fatalf("create user error: %v", err)
    }

    got, err := d.GetUserByEmail(context.Background(), "test@example.com")
    if err != nil {
        t.Fatalf("get user error: %v", err)
    }
    if got.Email != u.Email || got.Username != u.Username || got.Nickname != u.Nickname {
        t.Fatalf("unexpected user: %+v", got)
    }
    if got.Password == "pass123" || got.Password == "" {
        t.Fatalf("password should be hashed and non-empty")
    }
}

func TestUserDao_CheckUserExist(t *testing.T) {
    db := setupTestDB(t)
    d := NewUserDao(db)

    exists, err := d.CheckUserExist(context.Background(), "none@example.com")
    if err != nil || exists {
        t.Fatalf("expected not exist, got exists=%v err=%v", exists, err)
    }

    u := &model.User{Email: "exist@example.com", Username: "user2", Nickname: "user2"}
    if err := u.SetPassword("pass123"); err != nil {
        t.Fatalf("set password error: %v", err)
    }
    if err := d.CreateUser(context.Background(), u); err != nil {
        t.Fatalf("create user error: %v", err)
    }

    exists, err = d.CheckUserExist(context.Background(), "exist@example.com")
    if err != nil || !exists {
        t.Fatalf("expected exist, got exists=%v err=%v", exists, err)
    }
}
