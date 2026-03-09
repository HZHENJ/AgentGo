package v1

import (
	"net/http"
	"agentgo/internal/dao"
	cache "agentgo/internal/cache"
	redis "agentgo/internal/common/redis"
	db "agentgo/internal/common/mysql"

	"agentgo/internal/service"
	"agentgo/internal/types"
	"agentgo/pkg/ctl"
	"agentgo/pkg/e"
	"github.com/gin-gonic/gin"
	"strings"
)

// UserRegister handles user registration requests.
// It validates the input, calls the UserService to
// register the user, and returns the appropriate response.
func UserRegister(c *gin.Context) {
	app := ctl.NewWrapper(c)
	
	var req types.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return 
	}

	userDao := dao.NewUserDao(db.DB)
	userCacheDao := cache.NewUserCacheDao(redis.RDB)
	userService := service.NewUserService(userDao, userCacheDao)

	data, code := userService.Register(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}

	app.Success(data)
}

// UserLogin handles user login requests. It validates the input, calls the UserService to
func UserLogin(c *gin.Context) {
	app := ctl.NewWrapper(c)

	var req types.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(e.INVALID_PARAMS, err)
		return 
	}

	userDao := dao.NewUserDao(db.DB)
	userCacheDao := cache.NewUserCacheDao(redis.RDB)
	userService := service.NewUserService(userDao, userCacheDao)

	data, code := userService.Login(c.Request.Context(), &req)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}

	app.Success(data)
}

// UserLogout handles user logout requests. It retrieves the JWT token from the Authorization header,
func UserLogout(c *gin.Context) {
	app := ctl.NewWrapper(c)
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		app.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		app.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	token := parts[1]
	userDao := dao.NewUserDao(db.DB)
	userCacheDao := cache.NewUserCacheDao(redis.RDB)
	userService := service.NewUserService(userDao, userCacheDao)

	data, code := userService.Logout(c.Request.Context(), token)
	if code != e.SUCCESS {
		app.Response(http.StatusOK, code, nil)
		return
	}
	app.Success(data)
}