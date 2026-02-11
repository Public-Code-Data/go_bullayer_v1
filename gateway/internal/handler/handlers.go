package handler

import (
	"net/http"

	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/gateway/internal/config"
	"go_bullayer_v1/gateway/internal/router"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ServiceContext 服务上下文
// 包含Gateway服务运行所需的所有依赖和配置
type ServiceContext struct {
	Config config.Config  // 服务配置
	Router *router.Router // 路由管理器
}

// NewServiceContext 创建服务上下文
// c: 服务配置
// 返回服务上下文实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Router: router.NewRouter(c.Routes),
	}
}

// RegisterHandlers 注册所有路由处理器
// server: RESTful 服务器实例
// ctx: 服务上下文
func RegisterHandlers(server *rest.Server, ctx *ServiceContext) {
	// 健康检查接口
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/gateway/health",
		Handler: HealthHandler(ctx),
	})

	// 根据配置动态注册路由
	for _, route := range ctx.Config.Routes {
		server.AddRoute(rest.Route{
			Method:  route.Method,
			Path:    route.Path,
			Handler: ProxyHandler(ctx, route),
		})
		logger.Info("已注册路由: %s %s -> %s", route.Method, route.Path, route.Backend)
	}
}

// HealthHandler 健康检查处理器
// ctx: 服务上下文
// 返回 HTTP 处理器函数
func HealthHandler(ctx *ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
			"status":  "ok",
			"message": "Gateway服务运行正常",
			"routes":  len(ctx.Config.Routes),
		})
	}
}

// ProxyHandler 代理处理器
// 将请求转发到后端服务
// ctx: 服务上下文
// route: 路由配置
// 返回 HTTP 处理器函数
func ProxyHandler(ctx *ServiceContext, route config.RouteConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 使用路由管理器转发请求
		resp, err := ctx.Router.Forward(r, route)
		if err != nil {
			// 转发失败，返回降级响应
			logger.Error("请求转发失败: %v", err)
			httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{
				"code":    503,
				"message": route.Fallback,
				"data":    nil,
			})
			return
		}

		// 转发成功，返回后端响应
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
