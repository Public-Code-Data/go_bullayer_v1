package handler

import (
	"net/http"

	"go_bullayer_v1/api/internal/logic"
	"go_bullayer_v1/api/internal/svc"
	"go_bullayer_v1/api/internal/types"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RegisterHandlers 注册所有路由处理器
// server: RESTful 服务器实例
// ctx: 服务上下文
func RegisterHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/health",
				Handler: HealthHandler(ctx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/user/:id",
				Handler: UserHandler(ctx),
			},
		},
	)
}

// HealthHandler 健康检查处理器
// ctx: 服务上下文
// 返回 HTTP 处理器函数
func HealthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HealthRequest
		// 解析请求参数（如果有）
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 创建业务逻辑处理器
		l := logic.NewHealthLogic(r.Context(), ctx)
		// 执行业务逻辑
		resp, err := l.Health(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// UserHandler 用户信息处理器
// ctx: 服务上下文
// 返回 HTTP 处理器函数
func UserHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRequest
		// 解析路径参数
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 创建业务逻辑处理器
		l := logic.NewUserLogic(r.Context(), ctx)
		// 执行业务逻辑
		resp, err := l.GetUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
