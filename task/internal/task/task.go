package task

import "context"

// Task 任务接口
// 所有后台任务都需要实现此接口
type Task interface {
	// Name 返回任务名称
	Name() string

	// Execute 执行任务
	// ctx: 上下文，用于控制任务执行和取消
	// 返回错误信息
	Execute(ctx context.Context) error
}
