package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	baseconfig "go_bullayer_v1/base/pkg/config"
	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/task/internal/config"
	"go_bullayer_v1/task/internal/service"
)

// configFile 配置文件路径参数
var configFile = flag.String("f", "etc/task.yaml", "配置文件路径")

// main 程序入口函数
func main() {
	// 解析命令行参数
	flag.Parse()

	// 使用 base 模块的配置加载功能
	var c config.Config
	baseconfig.MustLoadConfig(*configFile, &c)

	// 使用 base 模块的日志初始化功能
	logger.InitLogger("task-service")
	logger.Info("后台任务服务开始启动...")

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建任务服务
	taskService := service.NewTaskService(ctx, c)

	// 启动任务服务
	taskService.Start()

	// 输出启动信息
	fmt.Println("后台任务服务启动成功")
	logger.Info("后台任务服务启动成功")

	// 监听系统信号，优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待退出信号
	<-sigChan
	fmt.Println("收到退出信号，正在关闭服务...")
	logger.Info("收到退出信号，正在关闭服务")

	// 停止任务服务
	taskService.Stop()

	fmt.Println("后台任务服务已关闭")
	logger.Info("后台任务服务已关闭")
}
