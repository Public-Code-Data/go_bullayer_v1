package config

import "github.com/zeromicro/go-zero/rest"

// Config Gateway服务配置结构
type Config struct {
	rest.RestConf // go-zero RESTful 服务器配置

	// QPS限制配置
	RateLimit struct {
		Enabled bool `json:"enabled"` // 是否启用QPS限制
		QPS     int  `json:"qps"`     // 每秒请求数限制
	} `json:"RateLimit"`

	// 熔断降级配置
	CircuitBreaker struct {
		Enabled bool `json:"enabled"` // 是否启用熔断器
		Timeout int  `json:"timeout"` // 超时时间（秒）
	} `json:"CircuitBreaker"`

	// 路由配置
	Routes []RouteConfig `json:"Routes"`
}

// RouteConfig 路由配置
// 定义请求路由规则和后端服务映射
type RouteConfig struct {
	Path     string `json:"path"`     // 请求路径（支持通配符）
	Backend  string `json:"backend"`  // 后端服务地址
	Method   string `json:"method"`   // HTTP方法（GET, POST, PUT, DELETE等）
	Fallback string `json:"fallback"` // 降级后的响应内容
	Timeout  int    `json:"timeout"`  // 超时时间（秒）
}
