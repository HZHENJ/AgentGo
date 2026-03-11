package dao

import (
    "context"
    "testing"
    "time"

    "agentgo/internal/model"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupSessionTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open sqlite: %v", err)
    }
    if err := db.AutoMigrate(&model.Session{}); err != nil {
        t.Fatalf("failed to migrate: %v", err)
    }
    return db
}

func TestSessionDao_Create_List_Get_Delete(t *testing.T) {
    db := setupSessionTestDB(t)
    d := NewSessionDao(db)

    ctx := context.Background()

    s1 := &model.Session{Username: "alice", Title: "first"}
    if err := d.CreateSession(ctx, s1); err != nil {
        t.Fatalf("create s1 error: %v", err)
    }
    // 做一个更晚的创建时间
    time.Sleep(10 * time.Millisecond)
    s2 := &model.Session{Username: "alice", Title: "second"}
    if err := d.CreateSession(ctx, s2); err != nil {
        t.Fatalf("create s2 error: %v", err)
    }

    // List by username，应按 created_at DESC 排序，second 在前
    list, err := d.GetSessionsByUsername(ctx, "alice")
    if err != nil {
        t.Fatalf("list error: %v", err)
    }
    if len(list) != 2 || list[0].Title != "second" || list[1].Title != "first" {
        t.Fatalf("unexpected order or size: %+v", list)
    }

    // Get by id
    got, err := d.GetSessionByID(ctx, s1.ID)
    if err != nil || got.ID != s1.ID || got.Title != "first" {
        t.Fatalf("get by id unexpected: got=%+v err=%v", got, err)
    }

    // Delete with wrong username: 不应删除
    if err := d.DeleteSessionByID(ctx, s1.ID, "bob"); err != nil {
        t.Fatalf("delete with wrong user should not error: %v", err)
    }
    if _, err := d.GetSessionByID(ctx, s1.ID); err != nil {
        t.Fatalf("record should still exist, got err=%v", err)
    }

    // Delete with correct username: 应删除
    if err := d.DeleteSessionByID(ctx, s1.ID, "alice"); err != nil {
        t.Fatalf("delete error: %v", err)
    }
    if _, err := d.GetSessionByID(ctx, s1.ID); err == nil {
        t.Fatalf("expected not found after delete")
    }
}
