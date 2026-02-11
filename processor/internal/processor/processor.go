package processor

import "context"

// Processor 数据处理任务接口
type Processor interface {
	// Name 返回任务名称
	Name() string

	// Execute 执行任务
	Execute(ctx context.Context) error
}
