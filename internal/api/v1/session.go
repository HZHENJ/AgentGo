package v1

import (
	"fmt"
	db "agentgo/internal/common/mysql"
	"agentgo/internal/dao"
	"agentgo/internal/service"
	"agentgo/internal/types"
	"agentgo/pkg/ctl"
	"agentgo/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SessionList(c *gin.Context) {
	app := ctl.NewWrapper(c)
	username := c.GetString("Username") 

	sessionDao := dao.NewSessionDao(db.DB)
	messageDao := dao.NewMessageDao(db.DB)
	svc := service.NewSessionService(sessionDao, messageDao)
	data, code := svc.GetSessionList(c.Request.Context(), &types.GetSessionListRequest{Username: username})
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}

// SessionCreateAndStream 第一次提问：创建会话并直接流式输出
func SessionCreateAndStream(c *gin.Context) {
	app := ctl.NewWrapper(c)
	username := c.GetString("Username")

	// TODO: 这里休要调整一下
	var req types.StreamChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	sessionDao := dao.NewSessionDao(db.DB)
	messageDao := dao.NewMessageDao(db.DB)
	svc := service.NewSessionService(sessionDao, messageDao)

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		app.Response(http.StatusOK, e.ERROR_STREAM_RESPONSE_FAIL, nil)
		return
	}

	// 1. 先通过 Service 创建 Session (标题暂时用用户的问题)
	createReq := &types.CreateSessionRequest{Username: username, Title: req.Question}
	data, code := svc.CreateSession(c.Request.Context(), createReq)
	if code != e.SUCCESS {
		c.Writer.Write([]byte("event: error\ndata: Failed to create session\n\n"))
		flusher.Flush()
		return
	}
	sessionID := data.(*types.CreateSessionResponse).SessionID

	c.Writer.Write([]byte(fmt.Sprintf("data: {\"sessionId\": %d}\n\n", sessionID)))
	flusher.Flush()

	req.SessionID = sessionID

	write := func(delta string) {
		c.Writer.Write([]byte("data: " + delta + "\n\n"))
		flusher.Flush()
	}

	if _, code := svc.StreamChat(c.Request.Context(), &req, write); code != e.SUCCESS {
		c.Writer.Write([]byte("event: error\ndata: failed\n\n"))
		flusher.Flush()
		return
	}

	c.Writer.Write([]byte("event: end\ndata: done\n\n"))
	flusher.Flush()
}

// SessionCreate 创建会话
func SessionCreate(c *gin.Context) {
	app := ctl.NewWrapper(c)

	var req types.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	sessionDao := dao.NewSessionDao(db.DB)
	messageDao := dao.NewMessageDao(db.DB)
	svc := service.NewSessionService(sessionDao, messageDao)

	data, code := svc.CreateSession(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}

// SessionHistory 获取会话历史
func SessionHistory(c *gin.Context) {
	app := ctl.NewWrapper(c)

	var req types.GetHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	sessionDao := dao.NewSessionDao(db.DB)
	messageDao := dao.NewMessageDao(db.DB)
	svc := service.NewSessionService(sessionDao, messageDao)

	data, code := svc.GetChatHistory(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}

// SessionStream
func SessionStream(c *gin.Context) {
	var req types.StreamChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app := ctl.NewWrapper(c)
		app.Error(e.INVALID_PARAMS, err)
		return
	}

	// 服务实例
	sessionDao := dao.NewSessionDao(db.DB)
	messageDao := dao.NewMessageDao(db.DB)
	svc := service.NewSessionService(sessionDao, messageDao)

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Accel-Buffering", "no")

	// 直接写增量数据给客户端
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		app := ctl.NewWrapper(c)
		app.Response(http.StatusOK, e.ERROR_STREAM_RESPONSE_FAIL, nil)
		return
	}

	write := func(delta string) {
		// SSE 格式：data: <chunk>\n\n
		_, _ = c.Writer.Write([]byte("data: "))
		_, _ = c.Writer.Write([]byte(delta))
		_, _ = c.Writer.Write([]byte("\n\n"))
		flusher.Flush()
	}

	if _, code := svc.StreamChat(c.Request.Context(), &req, write); code != e.SUCCESS {
		// 失败时返回一个完结事件，前端可据此关闭连接
		_, _ = c.Writer.Write([]byte("event: error\n"))
		_, _ = c.Writer.Write([]byte("data: failed\n\n"))
		flusher.Flush()
		return
	}

	// 成功完结事件
	_, _ = c.Writer.Write([]byte("event: end\n"))
	_, _ = c.Writer.Write([]byte("data: done\n\n"))
	flusher.Flush()
}
