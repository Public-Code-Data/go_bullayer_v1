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
	"go_bullayer_v1/processor/internal/config"
	"go_bullayer_v1/processor/internal/service"
)

// configFile 配置文件路径参数
var configFile = flag.String("f", "etc/processor.yaml", "配置文件路径")

// main 程序入口函数
func main() {
	flag.Parse()

	var c config.Config
	baseconfig.MustLoadConfig(*configFile, &c)

	logger.InitLogger("processor-service")
	logger.Info("数据处理服务开始启动...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processorService := service.NewProcessorService(ctx, c)
	processorService.Start()

	fmt.Println("数据处理服务启动成功")
	logger.Info("数据处理服务启动成功")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("收到退出信号，正在关闭服务...")
	logger.Info("收到退出信号，正在关闭服务")

	processorService.Stop()

	fmt.Println("数据处理服务已关闭")
	logger.Info("数据处理服务已关闭")
}
