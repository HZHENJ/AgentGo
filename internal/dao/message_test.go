package dao

import (
    "context"
    "testing"
    "time"

    "agentgo/internal/model"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupMessageTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open sqlite: %v", err)
    }
    if err := db.AutoMigrate(&model.Message{}); err != nil {
        t.Fatalf("failed to migrate: %v", err)
    }
    return db
}

func TestMessageDao_Create_And_History(t *testing.T) {
    db := setupMessageTestDB(t)
    d := NewMessageDao(db)
    ctx := context.Background()

    // session 1
    m1 := &model.Message{SessionID: 1, Username: "alice", Content: "hello", IsUser: true}
    m1.CreatedAt = time.Now().Add(-2 * time.Second)
    if err := d.CreateMessage(ctx, m1); err != nil { t.Fatalf("create m1: %v", err) }

    m2 := &model.Message{SessionID: 1, Username: "alice", Content: "hi", IsUser: false}
    m2.CreatedAt = time.Now().Add(-1 * time.Second)
    if err := d.CreateMessage(ctx, m2); err != nil { t.Fatalf("create m2: %v", err) }

    // session 2 (noise)
    m3 := &model.Message{SessionID: 2, Username: "bob", Content: "ignore", IsUser: true}
    if err := d.CreateMessage(ctx, m3); err != nil { t.Fatalf("create m3: %v", err) }

    // fetch history for session 1, expect ascending by created_at
    hs, err := d.GetHistoryBySessionID(ctx, 1)
    if err != nil { t.Fatalf("history error: %v", err) }
    if len(hs) != 2 || hs[0].Content != "hello" || hs[1].Content != "hi" {
        t.Fatalf("unexpected history: %+v", hs)
    }
}
