package service

import (
	"context"
	"sync"
	"time"

	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/task/internal/config"
	"go_bullayer_v1/task/internal/task"
)

// TaskService 任务服务
// 负责管理和执行所有后台任务
type TaskService struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 取消函数
	config config.Config      // 配置
	tasks  []task.Task        // 任务列表
	wg     sync.WaitGroup     // 等待组，用于等待所有任务完成
	mu     sync.Mutex         // 互斥锁
}

// NewTaskService 创建任务服务
// ctx: 上下文
// cfg: 服务配置
// 返回任务服务实例
func NewTaskService(ctx context.Context, cfg config.Config) *TaskService {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	return &TaskService{
		ctx:    ctxWithCancel,
		cancel: cancel,
		config: cfg,
		tasks:  make([]task.Task, 0),
	}
}

// Start 启动任务服务
func (s *TaskService) Start() {
	if !s.config.TaskEnabled {
		logger.Info("任务服务未启用")
		return
	}

	logger.Info("开始启动任务服务...")

	// 注册所有任务
	s.registerTasks()

	// 启动所有任务
	for _, t := range s.tasks {
		s.wg.Add(1)
		go s.runTask(t)
	}

	logger.Info("任务服务启动完成，共 %d 个任务", len(s.tasks))
}

// Stop 停止任务服务
func (s *TaskService) Stop() {
	logger.Info("开始停止任务服务...")

	// 取消上下文，通知所有任务停止
	s.cancel()

	// 等待所有任务完成
	s.wg.Wait()

	logger.Info("任务服务已停止")
}

// registerTasks 注册所有任务
func (s *TaskService) registerTasks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 注册统计任务
	if s.config.StatsTask.Enabled {
		statsTask := task.NewStatsTask(s.config)
		s.tasks = append(s.tasks, statsTask)
		logger.Info("已注册统计任务")
	}

	// 可以在这里注册更多任务
	// 例如：数据清理任务、报表生成任务等
}

// runTask 运行单个任务
// t: 任务实例
func (s *TaskService) runTask(t task.Task) {
	defer s.wg.Done()

	ticker := time.NewTicker(time.Duration(s.config.Interval) * time.Second)
	defer ticker.Stop()

	// 立即执行一次
	s.executeTask(t)

	// 定时执行
	for {
		select {
		case <-ticker.C:
			s.executeTask(t)
		case <-s.ctx.Done():
			logger.Info("任务 %s 收到停止信号", t.Name())
			return
		}
	}
}

// executeTask 执行任务
// t: 任务实例
func (s *TaskService) executeTask(t task.Task) {
	logger.Info("开始执行任务: %s", t.Name())

	startTime := time.Now()

	// 执行任务
	if err := t.Execute(s.ctx); err != nil {
		logger.Error("任务 %s 执行失败: %v", t.Name(), err)
	} else {
		duration := time.Since(startTime)
		logger.Info("任务 %s 执行成功，耗时: %v", t.Name(), duration)
	}
}
