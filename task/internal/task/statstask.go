package task

import (
	"context"
	"time"

	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/base/pkg/utils"
	"go_bullayer_v1/task/internal/config"
)

// StatsTask 统计任务
// 负责执行数据统计相关的后台任务
type StatsTask struct {
	config config.Config // 任务配置
}

// NewStatsTask 创建统计任务
// cfg: 服务配置
// 返回统计任务实例
func NewStatsTask(cfg config.Config) *StatsTask {
	return &StatsTask{
		config: cfg,
	}
}

// Name 返回任务名称
func (t *StatsTask) Name() string {
	return "统计任务"
}

// Execute 执行统计任务
// ctx: 上下文
// 返回错误信息
func (t *StatsTask) Execute(ctx context.Context) error {
	// 检查是否在指定时间执行
	if !t.shouldExecute() {
		return nil
	}

	logger.Info("开始执行统计任务...")

	// TODO: 在这里实现具体的统计逻辑
	// 例如：
	// 1. 统计用户数据
	// 2. 统计订单数据
	// 3. 生成统计报表
	// 4. 更新统计数据到数据库

	// 示例：模拟统计操作
	now := utils.Now()
	currentTime := utils.FormatDateTime(now)
	logger.Info("统计任务执行中，当前时间: %s", currentTime)

	// 模拟数据处理
	time.Sleep(100 * time.Millisecond)

	logger.Info("统计任务执行完成")
	return nil
}

// shouldExecute 判断是否应该执行任务
// 根据配置的执行时间判断
func (t *StatsTask) shouldExecute() bool {
	// 如果配置了执行时间，检查当前时间是否匹配
	if t.config.StatsTask.Hour >= 0 && t.config.StatsTask.Minute >= 0 {
		now := time.Now()
		return now.Hour() == t.config.StatsTask.Hour &&
			now.Minute() == t.config.StatsTask.Minute
	}

	// 如果没有配置执行时间，每次都执行
	return true
}
