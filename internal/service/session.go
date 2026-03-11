package service

import (
	"agentgo/internal/dao"
	"agentgo/internal/llm"
	"agentgo/internal/model"
	"agentgo/internal/types"
	"agentgo/pkg/conf"
	"agentgo/pkg/e"
	"context"
	"log"
)

type SessionService struct {
	sessionDao dao.SessionDao
	messageDao dao.MessageDao
}

func NewSessionService(sessionDao dao.SessionDao, messageDao dao.MessageDao) *SessionService {
	return &SessionService{
		sessionDao: sessionDao,
		messageDao: messageDao,
	}
}

// CreateSession creates a new session for the user and returns the session ID.
func (s *SessionService) CreateSession(ctx context.Context, req *types.CreateSessionRequest) (interface{}, int) {
	session := &model.Session{
		Username: req.Username,
		Title:    req.Title, // using user first message as session title for better experience
	}

	if err := s.sessionDao.CreateSession(ctx, session); err != nil {
		log.Printf("[CreateSession] DB Error: %v", err)
		return nil, e.ERROR_SESSION_CREATE_FAIL
	}

	resp := &types.CreateSessionResponse{
		SessionID: session.ID,
	}

	return resp, e.SUCCESS
}

// GetChatHistory 加载历史聊天记录
func (s *SessionService) GetChatHistory(ctx context.Context, req *types.GetHistoryRequest) (interface{}, int) {
	messages, err := s.messageDao.GetHistoryBySessionID(ctx, req.SessionID)
	if err != nil {
		log.Printf("[GetChatHistory] DB Error: %v", err)
		return nil, e.ERROR_HISTORY_LOAD_FAIL
	}

	history := make([]model.History, 0, len(messages))
	for _, msg := range messages {
		history = append(history, model.History{
			IsUser:  msg.IsUser,
			Content: msg.Content,
		})
	}

	resp := &types.GetHistoryResponse{
		History: history,
	}

	return resp, e.SUCCESS
}

func (s *SessionService) StreamChat(ctx context.Context, req *types.StreamChatRequest, cb func(string)) (interface{}, int) {
	cfg := conf.Config.LLM

	modelObj, err := llm.CreateModel(ctx, req.ModelType, &cfg)
	if err != nil {
		log.Printf("[StreamChat] CreateModel error: %v", err)
		return nil, e.ERROR_LLM_CREATE_FAIL
	}

	helper := llm.NewHelper(modelObj, req.SessionID)
	helper.SetSaveFunc(func(msg *model.Message) error {
		return s.messageDao.CreateMessage(ctx, msg)
	})

	historyMsgs, err := s.messageDao.GetHistoryBySessionID(ctx, req.SessionID)
	if err == nil {
		for _, m := range historyMsgs {
			helper.AddMessage(m.Content, m.Username, m.IsUser, false)
		}
	}

	_, err = helper.StreamResponse(ctx, req.Username, req.Question, cb)
	if err != nil {
		log.Printf("[StreamChat] StreamResponse error: %v", err)
		return nil, e.ERROR_STREAM_RESPONSE_FAIL
	}

	return nil, e.SUCCESS
}

// GetSessionList
func (s *SessionService) GetSessionList(ctx context.Context, req *types.GetSessionListRequest) (interface{}, int) {
	// 调用 DAO 获取历史会话
	sessions, err := s.sessionDao.GetSessionsByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("[GetSessionList] DB Error: %v", err)
		return nil, e.ERROR // 或者定义具体的 e.ERROR_SESSION_LIST_FAIL
	}

	var sessionInfos []types.SessionInfo
	for _, session := range sessions {
		sessionInfos = append(sessionInfos, types.SessionInfo{
			SessionID: session.ID,
			Title:     session.Title,
			CreatedAt: session.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp := &types.GetSessionListResponse{
		Sessions: sessionInfos,
	}

	return resp, e.SUCCESS
}