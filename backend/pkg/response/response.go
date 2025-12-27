package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 响应码，0 表示成功
	Message string      `json:"message"`           // 响应消息
	Data    interface{} `json:"data"`              // 响应数据，成功时返回数据，失败时为 null
}

const (
	CodeSuccess        = 0     // 成功
	CodeBadRequest    = 400   // 请求参数错误
	CodeUnauthorized  = 401   // 未授权
	CodeNotFound      = 404   // 资源不存在
	CodeConflict      = 409   // 资源冲突
	CodeInternalError = 500   // 服务器内部错误
)

const (
	MessageSuccess        = "success"
	MessageBadRequest    = "Bad request"
	MessageUnauthorized  = "Unauthorized"
	MessageNotFound      = "Not found"
	MessageConflict      = "Conflict"
	MessageInternalError = "Internal server error"
)

// Success 成功响应 (200)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    data,
	})
}

// Created 创建成功响应 (201)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MessageUnauthorized
	}
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MessageNotFound
	}
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// Conflict 409 错误
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = MessageConflict
	}
	Error(c, http.StatusConflict, CodeConflict, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, CodeInternalError, MessageInternalError)
}
