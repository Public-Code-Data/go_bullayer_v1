package common

import "fmt"

// BaseError 基础错误类型
// 用于统一错误处理，包含错误码和错误信息
type BaseError struct {
	Code    int    // 错误码
	Message string // 错误信息
}

// Error 实现 error 接口
func (e *BaseError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewError 创建基础错误
// code: 错误码
// message: 错误信息
func NewError(code int, message string) *BaseError {
	return &BaseError{
		Code:    code,
		Message: message,
	}
}

// CommonErrorCode 通用错误码定义
const (
	ErrCodeSuccess      = 0    // 成功
	ErrCodeInvalidParam = 1001 // 参数错误
	ErrCodeNotFound     = 1002 // 未找到
	ErrCodeInternal     = 1003 // 内部错误
	ErrCodeUnauthorized = 1004 // 未授权
	ErrCodeForbidden    = 1005 // 禁止访问
)
