package logger

import (
	"log"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
)

// InitLogger 初始化日志系统
// serviceName: 服务名称，用于日志标识
func InitLogger(serviceName string) {
	logx.MustSetup(logx.LogConf{
		ServiceName:         serviceName,
		Mode:                "file",
		Path:                "logs",
		Level:               "info",
		Compress:            true,
		KeepDays:            7,
		StackCooldownMillis: 100,
	})
}

// GetLogger 获取标准日志实例
// 返回标准库的 log.Logger，用于简单的日志输出
func GetLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}

// Info 输出信息日志
func Info(format string, v ...interface{}) {
	logx.Infof(format, v...)
}

// Error 输出错误日志
func Error(format string, v ...interface{}) {
	logx.Errorf(format, v...)
}

// Debug 输出调试日志
func Debug(format string, v ...interface{}) {
	logx.Infof(format, v...)
}
