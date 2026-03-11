package routes

import (
    "github.com/gin-gonic/gin"
    v1 "agentgo/internal/api/v1"
    "agentgo/internal/middleware"
)

// NewRouter 初始化 Gin 路由
func NewRouter() *gin.Engine {
    r := gin.Default()

    r.Use(middleware.Cors())
    // r.GET("/ping", func(c *gin.Context) {
    //     c.JSON(200, gin.H{"msg":"pong","status":"ok"})
    // })

    api := r.Group("/api/v1")

    // 公开路由
    userPublic := api.Group("/user")
    {
        userPublic.POST("/register", v1.UserRegister)
        userPublic.POST("/login", v1.UserLogin)
        userPublic.POST("/captcha", v1.UserSendCaptcha)
    }

    // 需要鉴权的路由
    authed := api.Group("")
    authed.Use(middleware.JWT())
    {
        // 用户
        authed.POST("/user/logout", v1.UserLogout)

        // 会话
        authed.POST("/session/list", v1.SessionList)
        authed.POST("/session/create", v1.SessionCreate)
        authed.POST("/session/history", v1.SessionHistory)
        authed.POST("/session/stream", v1.SessionStream)
        authed.POST("/session/create-and-stream", v1.SessionCreateAndStream)
    }

    return r
}
