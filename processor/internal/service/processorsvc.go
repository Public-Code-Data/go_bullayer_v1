package service

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"go_bullayer_v1/base/pkg/db"
	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/processor/internal/config"
	"go_bullayer_v1/processor/internal/processor"
)

// ProcessorService 数据处理服务
// 负责管理并执行链上处理任务
type ProcessorService struct {
	ctx        context.Context
	cancel     context.CancelFunc
	config     config.Config
	db         *sql.DB
	processors []processor.Processor
	wg         sync.WaitGroup
	mu         sync.Mutex
}

// NewProcessorService 创建数据处理服务
func NewProcessorService(ctx context.Context, cfg config.Config) *ProcessorService {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	svc := &ProcessorService{
		ctx:        ctxWithCancel,
		cancel:     cancel,
		config:     cfg,
		processors: make([]processor.Processor, 0),
	}
	svc.initDB()
	return svc
}

// Start 启动数据处理服务
func (s *ProcessorService) Start() {
	if !s.config.ProcessorEnabled {
		logger.Info("数据处理服务未启用")
		return
	}

	logger.Info("开始启动数据处理服务...")
	s.registerProcessors()

	for _, p := range s.processors {
		s.wg.Add(1)
		go s.runProcessor(p)
	}

	logger.Info("数据处理服务启动完成，共 %d 个处理任务", len(s.processors))
}

// Stop 停止数据处理服务
func (s *ProcessorService) Stop() {
	logger.Info("开始停止数据处理服务...")
	s.cancel()
	s.wg.Wait()

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			logger.Error("关闭数据库连接失败: %v", err)
		} else {
			logger.Info("数据库连接已关闭")
		}
	}

	logger.Info("数据处理服务已停止")
}

func (s *ProcessorService) initDB() {
	if s.config.Database.Host == "" {
		logger.Info("未配置数据库连接，跳过数据库初始化")
		return
	}

	dbConfig := db.DBConfig{
		Host:     s.config.Database.Host,
		Port:     s.config.Database.Port,
		User:     s.config.Database.User,
		Password: s.config.Database.Password,
		Database: s.config.Database.Database,
	}

	database, err := db.NewDB(dbConfig)
	if err != nil {
		logger.Error("数据库连接初始化失败: %v", err)
		return
	}

	s.db = database
	logger.Info("数据库连接初始化成功")
}

// registerProcessors 注册所有处理任务
func (s *ProcessorService) registerProcessors() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.config.BlockProcessor.Enabled {
		blockProcessor := processor.NewBlockProcessor(s.config)
		s.processors = append(s.processors, blockProcessor)
		logger.Info("已注册区块追踪解析任务")
	}
}

// runProcessor 循环执行单个处理任务
func (s *ProcessorService) runProcessor(p processor.Processor) {
	defer s.wg.Done()

	interval := s.config.Interval
	if interval <= 0 {
		interval = 10
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	s.executeProcessor(p)

	for {
		select {
		case <-ticker.C:
			s.executeProcessor(p)
		case <-s.ctx.Done():
			logger.Info("处理任务 %s 收到停止信号", p.Name())
			return
		}
	}
}

// executeProcessor 执行处理任务
func (s *ProcessorService) executeProcessor(p processor.Processor) {
	logger.Info("开始执行处理任务: %s", p.Name())
	startTime := time.Now()

	if err := p.Execute(s.ctx); err != nil {
		logger.Error("处理任务 %s 执行失败: %v", p.Name(), err)
		return
	}

	logger.Info("处理任务 %s 执行成功，耗时: %v", p.Name(), time.Since(startTime))
}
