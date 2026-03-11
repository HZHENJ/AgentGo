package middleware

import (
	"net/http"
	"strings"
	"agentgo/pkg/utils"
	"agentgo/pkg/e"
	"github.com/gin-gonic/gin"
)

// JWT
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_FAIL,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// TODO Redis blacklist check

		// parse Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT,
				"msg":  e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set("UserId", claims.UserId)
		c.Set("Username", claims.Username)

		c.Next()
	}
}
