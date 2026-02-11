package middleware

import (
	"net/http"

	"go_bullayer_v1/gateway/internal/limiter"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RateLimitMiddleware QPS限制中间件
// qps: 每秒允许的请求数
// 返回中间件函数
func RateLimitMiddleware(qps int) rest.Middleware {
	// 创建限流器
	rateLimiter := limiter.NewLimiter(qps)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 获取客户端IP作为限流key
			clientIP := getClientIP(r)

			// 检查是否允许请求
			if !rateLimiter.Allow(clientIP) {
				// 超过限制，返回429状态码
				httpx.ErrorCtx(r.Context(), w, &RateLimitError{
					Code:    429,
					Message: "请求过于频繁，请稍后再试",
				})
				return
			}

			// 允许请求，继续处理
			next(w, r)
		}
	}
}

// getClientIP 获取客户端IP地址
// 优先从 X-Forwarded-For 头获取，如果没有则使用 RemoteAddr
func getClientIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 头获取（经过代理的情况）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// 尝试从 X-Real-IP 头获取
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 使用 RemoteAddr
	return r.RemoteAddr
}

// RateLimitError 限流错误
type RateLimitError struct {
	Code    int
	Message string
}

func (e *RateLimitError) Error() string {
	return e.Message
}
