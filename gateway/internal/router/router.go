package router

import (
	"context"
	"io"
	"net/http"
	"time"

	"go_bullayer_v1/gateway/internal/config"
	"go_bullayer_v1/base/pkg/logger"
)

// Router 路由管理器
// 负责请求转发和路由管理
type Router struct {
	routes []config.RouteConfig // 路由配置列表
}

// NewRouter 创建路由管理器
// routes: 路由配置列表
// 返回路由管理器实例
func NewRouter(routes []config.RouteConfig) *Router {
	return &Router{
		routes: routes,
	}
}

// Forward 转发请求到后端服务
// r: 原始HTTP请求
// route: 路由配置
// 返回响应数据和错误信息
func (router *Router) Forward(r *http.Request, route config.RouteConfig) (interface{}, error) {
	// 构建后端服务URL
	backendURL := route.Backend + r.URL.Path

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(route.Timeout) * time.Second,
	}

	// 创建转发请求
	req, err := http.NewRequest(route.Method, backendURL, r.Body)
	if err != nil {
		return nil, err
	}

	// 复制请求头
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 设置请求超时
	if route.Timeout > 0 {
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(route.Timeout)*time.Second)
		defer cancel()
		req = req.WithContext(ctx)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("转发请求失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应 body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 返回响应数据
	// 这里简化处理，实际应该解析JSON等
	return map[string]interface{}{
		"status_code": resp.StatusCode,
		"body":        string(body),
	}, nil
}

// MatchRoute 匹配路由
// path: 请求路径
// method: HTTP方法
// 返回匹配的路由配置和是否匹配成功
func (router *Router) MatchRoute(path, method string) (*config.RouteConfig, bool) {
	for _, route := range router.routes {
		if route.Method == method && router.matchPath(route.Path, path) {
			return &route, true
		}
	}
	return nil, false
}

// matchPath 匹配路径
// pattern: 路径模式（支持通配符）
// path: 实际路径
// 返回是否匹配
func (router *Router) matchPath(pattern, path string) bool {
	// 简单的路径匹配实现
	// 支持通配符 * 和 **
	// TODO: 可以实现更复杂的路由匹配逻辑
	return pattern == path || pattern == "*" || pattern == "**"
}
