package llm

import (
	"context"
	"log"
	"sync"
	"agentgo/internal/model"
)

// Helper
type Helper interface {
	SetSaveFunc(saveFunc func(*model.Message) error)
	AddMessage(content string, username string, isUser bool, save bool)
	GetMessages() []*model.Message
	GenerateResponse(ctx context.Context, username string, userQuestion string) (*model.Message, error)
	StreamResponse(ctx context.Context, username string, userQuestion string, cb func(string)) (*model.Message, error)
}

// helper 结构体：私有化，隐藏内部锁和切片细节
type helper struct {
	model     Model
	messages  []*model.Message
	mu        sync.RWMutex
	sessionID uint
	saveFunc  func(*model.Message) error
}

// NewHelper 构造函数：返回接口类型
func NewHelper(m Model, sessionID uint) Helper {
	return &helper{
		model:     m,
		messages:  make([]*model.Message, 0),
		sessionID: sessionID,
	}
}

// SetSaveFunc 
func (h *helper) SetSaveFunc(saveFunc func(*model.Message) error) {
	h.saveFunc = saveFunc
}

// AddMessage
func (h *helper) AddMessage(content string, username string, isUser bool, save bool) {
	msg := &model.Message{
		SessionID: h.sessionID,
		Username:  username,
		Content:   content,
		IsUser:    isUser,
	}

	// 1. 写锁保护内存切片，防止并发读写 Panic
	h.mu.Lock()
	h.messages = append(h.messages, msg)
	h.mu.Unlock()

	// 2. 持久化（不带锁执行，防止阻塞）
	if save && h.saveFunc != nil {
		if err := h.saveFunc(msg); err != nil {
			log.Printf("[Helper] 持久化失败: %v", err)
		}
	}
}

// GetMessages
func (h *helper) GetMessages() []*model.Message {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	out := make([]*model.Message, len(h.messages))
	copy(out, h.messages)
	return out
}

// StreamResponse
func (h *helper) StreamResponse(ctx context.Context, username string, userQuestion string, cb func(string)) (*model.Message, error) {
	// 1. 记录并保存用户的提问
	h.AddMessage(userQuestion, username, true, true)

	// 2. 获取当前所有历史上下文
	history := h.GetMessages()

	// 3. 调用底层适配器请求大模型
	fullContent, err := h.model.StreamResponse(ctx, history, cb)
	if err != nil {
		return nil, err
	}

	// 4. 记录并保存 AI 的回复
	h.AddMessage(fullContent, username, false, true)

	return &model.Message{
		SessionID: h.sessionID,
		Username:  username,
		Content:   fullContent,
		IsUser:    false,
	}, nil
}

// GenerateResponse 同步响应逻辑
func (h *helper) GenerateResponse(ctx context.Context, username string, userQuestion string) (*model.Message, error) {
	h.AddMessage(userQuestion, username, true, true)
	history := h.GetMessages()

	fullContent, err := h.model.GenerateResponse(ctx, history)
	if err != nil {
		return nil, err
	}

	h.AddMessage(fullContent, username, false, true)
	return &model.Message{
		SessionID: h.sessionID,
		Username:  username,
		Content:   fullContent,
		IsUser:    false,
	}, nil
}