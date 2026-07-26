// Package response 统一 API 的返回格式
//
//	{"code": 200, "data": ...}
//
// 出错时 data 为可以直接展示给用户的错误信息。
package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/model"
)

// OK 返回成功响应
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
}

// Fail 返回错误响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"code": code, "data": msg})
}

// AbortWith 返回错误响应并终止后续 handler
func AbortWith(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"code": code, "data": msg})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) { Fail(c, http.StatusBadRequest, msg) }

// Unauthorized 未登录
func Unauthorized(c *gin.Context, msg string) { Fail(c, http.StatusUnauthorized, msg) }

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) { Fail(c, http.StatusNotFound, msg) }

// ServerError 服务端错误
func ServerError(c *gin.Context, msg string) { Fail(c, http.StatusInternalServerError, msg) }

// FromError 按错误类型选择合适的状态码：
// 参数校验错误返回 400，其余返回 500
func FromError(c *gin.Context, err error, fallback string) {
	var validationErr model.ValidationError
	if errors.As(err, &validationErr) {
		BadRequest(c, validationErr.Msg)
		return
	}
	ServerError(c, fallback)
}
