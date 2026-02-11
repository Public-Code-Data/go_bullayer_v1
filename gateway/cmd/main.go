package main

import (
	"flag"
	"fmt"
	"time"

	"go_bullayer_v1/gateway/internal/config"
	"go_bullayer_v1/gateway/internal/handler"
	"go_bullayer_v1/gateway/internal/middleware"
	baseconfig "go_bullayer_v1/base/pkg/config"
	"go_bullayer_v1/base/pkg/logger"

	"github.com/zeromicro/go-zero/rest"
)

// configFile 配置文件路径参数
var configFile = flag.String("f", "etc/gateway.yaml", "配置文件路径")

// main 程序入口函数
func main() {
	// 解析命令行参数
	flag.Parse()

	// 使用 base 模块的配置加载功能
	var c config.Config
	baseconfig.MustLoadConfig(*configFile, &c)

	// 使用 base 模块的日志初始化功能
	logger.InitLogger("gateway-service")
	logger.Info("Gateway服务开始启动...")

	// 创建 RESTful 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 注册中间件
	// 注意：中间件的注册顺序很重要，先注册的会先执行
	if c.RateLimit.Enabled {
		server.Use(middleware.RateLimitMiddleware(c.RateLimit.QPS))
		logger.Info("已启用QPS限制中间件，限制: %d QPS", c.RateLimit.QPS)
	}

	if c.CircuitBreaker.Enabled {
		timeout := time.Duration(c.CircuitBreaker.Timeout) * time.Second
		server.Use(middleware.CircuitBreakerMiddleware(timeout))
		logger.Info("已启用熔断降级中间件，超时: %s", timeout)
	}

	// 创建服务上下文
	ctx := handler.NewServiceContext(c)

	// 注册路由处理器
	handler.RegisterHandlers(server, ctx)

	// 输出启动信息
	fmt.Printf("Gateway服务启动成功，监听地址: %s:%d\n", c.Host, c.Port)
	logger.Info("Gateway服务启动成功，监听地址: %s:%d", c.Host, c.Port)

	// 启动服务器
	server.Start()
}
