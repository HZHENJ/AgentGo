package ctl

/*
This package provides a standardized way to handle API responses 
in the AgentGo application. It defines a Response struct for 
consistent response formatting and a Wrapper struct that encapsulates 
the Gin context for easy response handling.
*/

import (
	"net/http"
	"agentgo/pkg/e"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

type Wrapper struct {
	C *gin.Context
}

func NewWrapper(c *gin.Context) *Wrapper {
	return &Wrapper{C: c}
}

// Response sends a JSON response with the given HTTP status code, error code, and data.
func (w *Wrapper) Response(httpCode, errorCode int, data interface{}) {
	w.C.JSON(httpCode, Response{
		Code: errorCode,
		Msg:  e.GetMsg(errorCode),
		Data: data,
	})
}

// Success sends a successful response with the given data and a standard success code.
func (w *Wrapper) Success(data interface{}) {
	w.Response(http.StatusOK, e.SUCCESS, data)
}

// Error sends an error response with the given error code and error message.
func (w *Wrapper) Error(errCode int, err error) {
	// 只有在某种Debug开关开启时，才把 err.Error() 放入响应
	// 平时只返回 errCode 对应的 Msg
	msg := e.GetMsg(errCode)

	// 如果想要 track_id 或者日志打印，在这里进行 log.Println(err)

	w.C.JSON(http.StatusOK, Response{
		Code: errCode,
		Msg:  msg,
		Data: nil,
		// Error: err.Error(), // 生产环境建议注释掉这行，或者加开关判断
	})
}
