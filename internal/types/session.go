package types

import "agentgo/internal/model"

type CreateSessionRequest struct {
	Username string `json:"username" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

type CreateSessionResponse struct {
	SessionID uint `json:"session_id"`
}

type GetHistoryRequest struct {
	SessionID uint `json:"session_id" binding:"required"`
}

type GetHistoryResponse struct {
	History []model.History `json:"history"`
}

type StreamChatRequest struct {
	Username  string `json:"username" binding:"required"`
	SessionID uint   `json:"session_id" binding:"required"`
	Question  string `json:"question" binding:"required"`
	ModelType string `json:"model_type" binding:"required"` // such as: "openai", "ollama"
}

type GetSessionListRequest struct {
	Username string `json:"username" binding:"required"`
}

type SessionInfo struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type GetSessionListResponse struct {
	Sessions []SessionInfo `json:"sessions"`
}