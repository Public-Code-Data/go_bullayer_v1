package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	initOnce sync.Once

	client    *ethclient.Client
	clientURL string
	initErr   error
)

// InitClient 初始化 ETH 客户端，只会初始化一次。
// 如果后续用不同 RPCURL 再次初始化，会返回错误。
func InitClient(rpcURL string) error {
	rpcURL = strings.TrimSpace(rpcURL)
	if rpcURL == "" {
		return errors.New("rpcURL is required")
	}

	initOnce.Do(func() {
		c, err := ethclient.Dial(rpcURL)
		if err != nil {
			initErr = fmt.Errorf("dial eth rpc failed: %w", err)
			return
		}
		client = c
		clientURL = rpcURL
	})

	if initErr != nil {
		return initErr
	}
	if clientURL != rpcURL {
		return fmt.Errorf("eth client already initialized with %s", clientURL)
	}
	return nil
}

// IsInitialized 返回客户端是否已初始化成功。
func IsInitialized() bool {
	return client != nil
}

// LatestBlockNumber 获取链上最新区块高度。
func LatestBlockNumber(ctx context.Context) (uint64, error) {
	if client == nil {
		return 0, errors.New("eth client is not initialized")
	}

	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if header == nil || header.Number == nil {
		return 0, errors.New("latest header is nil")
	}
	return header.Number.Uint64(), nil
}

// BlockReceiptsByNumber 按区块号查询区块内全部回执。
func BlockReceiptsByNumber(ctx context.Context, blockNumber *big.Int) ([]*types.Receipt, error) {
	if client == nil {
		return nil, errors.New("eth client is not initialized")
	}
	if blockNumber == nil {
		return nil, errors.New("blockNumber is required")
	}

	blockRef := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64()))
	return client.BlockReceipts(ctx, blockRef)
}

// TransactionByHash 按交易哈希查询交易对象。
func TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction, error) {
	if client == nil {
		return nil, errors.New("eth client is not initialized")
	}
	tx, _, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// TransactionSender 查询交易发送方地址。
func TransactionSender(ctx context.Context, tx *types.Transaction, blockHash common.Hash, index uint) (common.Address, error) {
	if client == nil {
		return common.Address{}, errors.New("eth client is not initialized")
	}
	if tx == nil {
		return common.Address{}, errors.New("transaction is nil")
	}
	return client.TransactionSender(ctx, tx, blockHash, index)
}
