package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	// circuitBreaker 熔断器实例
	// 使用 go-zero 的 breaker 包实现熔断功能
	circuitBreaker = breaker.NewBreaker(
		breaker.WithName("gateway"),
	)
)

// CircuitBreakerMiddleware 熔断降级中间件
// 当后端服务异常或超时时，自动熔断并返回降级响应
// timeout: 请求超时时间，<=0 则不设置超时
// 返回中间件函数
func CircuitBreakerMiddleware(timeout time.Duration) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 使用熔断器执行请求
			var err error
			if timeout > 0 {
				ctx, cancel := context.WithTimeout(r.Context(), timeout)
				defer cancel()
				err = circuitBreaker.Do(func() error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						next(w, r.WithContext(ctx))
						return nil
					}
				})
			} else {
				err = circuitBreaker.Do(func() error {
					next(w, r)
					return nil
				})
			}

			// 如果熔断器触发，返回降级响应
			if err != nil {
				httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
					"code":    503,
					"message": "服务暂时不可用，请稍后重试",
					"data":    nil,
				})
			}
		}
	}
}
