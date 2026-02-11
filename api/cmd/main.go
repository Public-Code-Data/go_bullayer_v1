package main

import (
	"flag"
	"fmt"

	"go_bullayer_v1/api/internal/config"
	"go_bullayer_v1/api/internal/handler"
	"go_bullayer_v1/api/internal/svc"
	baseconfig "go_bullayer_v1/base/pkg/config"
	"go_bullayer_v1/base/pkg/logger"

	"github.com/zeromicro/go-zero/rest"
)

// configFile 配置文件路径参数
var configFile = flag.String("f", "etc/api.yaml", "配置文件路径")

// main 程序入口函数
func main() {
	// 解析命令行参数
	flag.Parse()

	// 使用 base 模块的配置加载功能
	var c config.Config
	baseconfig.MustLoadConfig(*configFile, &c)

	// 使用 base 模块的日志初始化功能
	logger.InitLogger("api-service")
	logger.Info("API服务开始启动...")

	// 创建 RESTful 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 注册路由处理器
	handler.RegisterHandlers(server, ctx)

	// 输出启动信息
	fmt.Printf("API服务启动成功，监听地址: %s:%d\n", c.Host, c.Port)
	logger.Info("API服务启动成功，监听地址: %s:%d", c.Host, c.Port)

	// 启动服务器
	server.Start()
}
