package service

import (
    "context"
    "errors"
    "strings"
    "testing"

    "agentgo/internal/dao"
    "agentgo/internal/llm"
    "agentgo/internal/model"
    "agentgo/internal/types"
    "agentgo/pkg/conf"
    "agentgo/pkg/e"
)

// --- Fakes / Mocks ---

type fakeSessionDao struct{ nextID uint }

func (f *fakeSessionDao) CreateSession(ctx context.Context, s *model.Session) error {
    f.nextID++
    s.ID = f.nextID
    return nil
}
func (f *fakeSessionDao) GetSessionsByUsername(ctx context.Context, username string) ([]*model.Session, error) {
    return nil, nil
}
func (f *fakeSessionDao) GetSessionByID(ctx context.Context, id uint) (*model.Session, error) {
    return nil, errors.New("not implemented")
}
func (f *fakeSessionDao) DeleteSessionByID(ctx context.Context, id uint, username string) error { return nil }

var _ dao.SessionDao = (*fakeSessionDao)(nil)

type fakeMessageDao struct{
    saved []*model.Message
    history map[uint][]*model.Message
}

func newFakeMessageDao() *fakeMessageDao { return &fakeMessageDao{history: map[uint][]*model.Message{}} }

func (d *fakeMessageDao) CreateMessage(ctx context.Context, msg *model.Message) error {
    d.saved = append(d.saved, msg)
    return nil
}
func (d *fakeMessageDao) GetHistoryBySessionID(ctx context.Context, sid uint) ([]*model.Message, error) {
    return d.history[sid], nil
}

var _ dao.MessageDao = (*fakeMessageDao)(nil)

// fake model implements llm.Model
type fakeModel struct{ }

func (m *fakeModel) GenerateResponse(ctx context.Context, msgs []*model.Message) (string, error) {
    return "reply", nil
}
func (m *fakeModel) StreamResponse(ctx context.Context, msgs []*model.Message, cb func(string)) (string, error) {
    parts := []string{"hel", "lo"}
    var b strings.Builder
    for _, p := range parts {
        cb(p)
        b.WriteString(p)
    }
    return b.String(), nil
}
func (m *fakeModel) GetModelType() string { return "mock-session-svc" }

// --- Tests ---

func TestSessionService_CreateSession(t *testing.T) {
    sd := &fakeSessionDao{}
    md := newFakeMessageDao()
    svc := NewSessionService(sd, md)

    data, code := svc.CreateSession(context.Background(), &types.CreateSessionRequest{Username: "alice", Title: "t1"})
    if code != e.SUCCESS { t.Fatalf("expected SUCCESS, got %d", code) }
    resp, ok := data.(*types.CreateSessionResponse)
    if !ok || resp.SessionID == 0 { t.Fatalf("invalid resp: %#v", data) }
}

func TestSessionService_GetChatHistory(t *testing.T) {
    sd := &fakeSessionDao{}
    md := newFakeMessageDao()
    svc := NewSessionService(sd, md)

    md.history[1] = []*model.Message{
        {SessionID:1, Username:"alice", Content:"hi", IsUser:true},
        {SessionID:1, Username:"alice", Content:"hello", IsUser:false},
    }

    data, code := svc.GetChatHistory(context.Background(), &types.GetHistoryRequest{SessionID: 1})
    if code != e.SUCCESS { t.Fatalf("expected SUCCESS, got %d", code) }
    resp := data.(*types.GetHistoryResponse)
    if len(resp.History) != 2 || !resp.History[0].IsUser || resp.History[1].IsUser {
        t.Fatalf("unexpected history: %#v", resp.History)
    }
}

func TestSessionService_StreamChat(t *testing.T) {
    // register a unique mock model type
    mockType := "mock-session-svc"
    llm.Register(mockType, func(ctx context.Context, cfg *conf.LLMConfig) (llm.Model, error) { return &fakeModel{}, nil })

    // set config (not strictly required by our mock, but keep structure)
    conf.Config = &conf.Configuration{LLM: conf.LLMConfig{Type: mockType}}

    sd := &fakeSessionDao{}
    md := newFakeMessageDao()
    // seed some history
    md.history[10] = []*model.Message{{SessionID:10, Username:"alice", Content:"prev", IsUser:true}}

    svc := NewSessionService(sd, md)

    var got strings.Builder
    _, code := svc.StreamChat(context.Background(), &types.StreamChatRequest{
        Username:  "alice",
        SessionID: 10,
        Question:  "ping",
        ModelType: mockType,
    }, func(delta string){ got.WriteString(delta) })

    if code != e.SUCCESS { t.Fatalf("expected SUCCESS, got %d", code) }
    if got.String() != "hello" { t.Fatalf("unexpected streamed: %q", got.String()) }
    // should save 2 messages: user question + assistant reply
    if len(md.saved) != 2 || !md.saved[0].IsUser || md.saved[1].IsUser {
        t.Fatalf("unexpected saved messages: %#v", md.saved)
    }
}
