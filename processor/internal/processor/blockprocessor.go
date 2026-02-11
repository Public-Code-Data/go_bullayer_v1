package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go_bullayer_v1/base/pkg/eth"
	"go_bullayer_v1/base/pkg/logger"
	"go_bullayer_v1/processor/internal/config"
	"go_bullayer_v1/processor/internal/core"
)

// BlockProcessor 区块处理任务
// 负责追踪链上区块高度并解析区块数据
type BlockProcessor struct {
	config         config.Config
	mu             sync.Mutex
	currentHeight  int64
	mockLatestHead int64
}

// NewBlockProcessor 创建区块处理任务
func NewBlockProcessor(cfg config.Config) *BlockProcessor {
	startHeight := cfg.Chain.StartHeight
	if startHeight < 0 {
		startHeight = 0
	}

	return &BlockProcessor{
		config:         cfg,
		currentHeight:  startHeight,
		mockLatestHead: startHeight + 50,
	}
}

// Name 返回任务名称
func (p *BlockProcessor) Name() string {
	return "区块追踪解析任务"
}

// Execute 执行区块追踪和解析逻辑
func (p *BlockProcessor) Execute(ctx context.Context) error {
	latestHeight, err := p.fetchLatestHeight(ctx)
	if err != nil {
		return err
	}

	p.mu.Lock()
	fromHeight := p.currentHeight + 1
	safeHeight := latestHeight - p.config.Chain.Confirmations
	if safeHeight < fromHeight {
		p.mu.Unlock()
		logger.Info("暂无可处理区块，当前=%d, 链上=%d, 安全高度=%d", p.currentHeight, latestHeight, safeHeight)
		return nil
	}

	maxPerRound := p.config.Chain.MaxBlocksPerRound
	if maxPerRound <= 0 {
		maxPerRound = 20
	}

	toHeight := safeHeight
	limitHeight := p.currentHeight + maxPerRound
	if toHeight > limitHeight {
		toHeight = limitHeight
	}
	p.mu.Unlock()

	for h := fromHeight; h <= toHeight; h++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := p.parseBlock(ctx, h); err != nil {
			return err
		}
	}

	p.mu.Lock()
	p.currentHeight = toHeight
	p.mu.Unlock()

	logger.Info("区块处理完成，已更新到高度 %d", toHeight)
	return nil
}

// fetchLatestHeight 获取链上最新区块高度
func (p *BlockProcessor) fetchLatestHeight(ctx context.Context) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if eth.IsInitialized() {
		latest, err := eth.LatestBlockNumber(ctx)
		if err != nil {
			return 0, err
		}
		return int64(latest), nil
	}

	// 降级路径：未初始化 ETH 客户端时使用本地模拟高度，保障流程可运行。
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mockLatestHead += 3
	return p.mockLatestHead, nil
}

// parseBlock 解析指定高度区块数据
func (p *BlockProcessor) parseBlock(ctx context.Context, height int64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Info("开始解析区块 %d", height)

	if p.config.BlockProcessor.ParseTx {
		receiptParser := core.NewReceiptParser()
		incomingTransfers, err := receiptParser.ParseAndFilterByBlock(
			ctx,
			height,
			p.config.BlockProcessor.TargetAddresses,
			p.config.BlockProcessor.TrackedAssets,
			p.config.BlockProcessor.TokenContracts,
			p.config.BlockProcessor.ParseWorkers,
		)
		if err != nil {
			return err
		}
		logger.Info("区块 %d 命中流入交易 %d 条", height, len(incomingTransfers))
	}

	if p.config.BlockProcessor.ParseEvent {
		logger.Info("解析区块 %d 事件日志", height)
	}

	// TODO: 在这里落库并更新业务索引
	time.Sleep(20 * time.Millisecond)
	logger.Info("区块 %d 解析完成", height)
	return nil
}

func (p *BlockProcessor) String() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return fmt.Sprintf("current_height=%d", p.currentHeight)
}
