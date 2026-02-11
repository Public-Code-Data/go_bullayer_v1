package core

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"go_bullayer_v1/base/pkg/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var transferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// ReceiptParser 票据解析器。
// 功能：按区块号拉取回执，并发解析交易，再过滤流入目标地址集合的资产转移记录。
type ReceiptParser struct {
	filter *TransferFilter
}

func NewReceiptParser() *ReceiptParser {
	return &ReceiptParser{
		filter: NewTransferFilter(),
	}
}

// ParseAndFilterByBlock 按区块解析并过滤转账记录。
//
// tokenSymbolsByAddress: ERC20合约地址 -> symbol，例如 {"0xdac17...":"USDT"}
func (p *ReceiptParser) ParseAndFilterByBlock(
	ctx context.Context,
	blockNumber int64,
	targetAddresses []string,
	trackedAssets []string,
	tokenSymbolsByAddress map[string]string,
	workerCount int,
) ([]TransferRecord, error) {
	if blockNumber < 0 {
		return nil, fmt.Errorf("invalid blockNumber: %d", blockNumber)
	}

	receipts, err := eth.BlockReceiptsByNumber(ctx, big.NewInt(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("query block receipts failed: %w", err)
	}

	normalizedTokenMap := normalizeTokenSymbolMap(tokenSymbolsByAddress)

	erc20Transfers := parseERC20TransfersFromReceipts(receipts, normalizedTokenMap)
	ethTransfers, err := p.parseNativeETHTransfersConcurrently(ctx, receipts, workerCount)
	if err != nil {
		return nil, err
	}

	allTransfers := make([]TransferRecord, 0, len(erc20Transfers)+len(ethTransfers))
	allTransfers = append(allTransfers, erc20Transfers...)
	allTransfers = append(allTransfers, ethTransfers...)

	return p.filter.FilterIncomingTransfers(targetAddresses, trackedAssets, allTransfers), nil
}

func parseERC20TransfersFromReceipts(
	receipts []*types.Receipt,
	tokenSymbolsByAddress map[string]string,
) []TransferRecord {
	results := make([]TransferRecord, 0, len(receipts))
	for _, receipt := range receipts {
		if receipt == nil {
			continue
		}
		for _, lg := range receipt.Logs {
			if lg == nil || len(lg.Topics) < 3 || lg.Topics[0] != transferTopic {
				continue
			}
			if len(lg.Data) != 32 {
				continue
			}

			from := topicToAddress(lg.Topics[1])
			to := topicToAddress(lg.Topics[2])
			amount := new(big.Int).SetBytes(lg.Data)
			tokenAddress := strings.ToLower(lg.Address.Hex())
			symbol := tokenSymbolsByAddress[tokenAddress]
			if symbol == "" {
				// 没在关注合约列表中时，用地址占位，留给后续过滤白名单处理。
				symbol = tokenAddress
			}

			results = append(results, TransferRecord{
				TxHash:       receipt.TxHash.Hex(),
				BlockNumber:  int64(receipt.BlockNumber.Uint64()),
				From:         from,
				To:           to,
				Amount:       amount.String(),
				AssetType:    AssetTypeERC20,
				TokenAddress: lg.Address.Hex(),
				TokenSymbol:  strings.ToUpper(symbol),
			})
		}
	}
	return results
}

func (p *ReceiptParser) parseNativeETHTransfersConcurrently(
	ctx context.Context,
	receipts []*types.Receipt,
	workerCount int,
) ([]TransferRecord, error) {
	if workerCount <= 0 {
		workerCount = 8
	}
	if workerCount > 64 {
		workerCount = 64
	}

	type parseResult struct {
		record TransferRecord
		ok     bool
		err    error
	}

	jobs := make(chan *types.Receipt)
	results := make(chan parseResult, len(receipts))

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for receipt := range jobs {
				if receipt == nil {
					continue
				}
				r, ok, err := parseETHTransferByReceipt(ctx, receipt)
				results <- parseResult{record: r, ok: ok, err: err}
			}
		}()
	}

	go func() {
		for _, receipt := range receipts {
			select {
			case <-ctx.Done():
				close(jobs)
				wg.Wait()
				close(results)
				return
			case jobs <- receipt:
			}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	records := make([]TransferRecord, 0, len(receipts))
	for result := range results {
		if result.err != nil {
			continue
		}
		if result.ok {
			records = append(records, result.record)
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func parseETHTransferByReceipt(
	ctx context.Context,
	receipt *types.Receipt,
) (TransferRecord, bool, error) {
	tx, err := eth.TransactionByHash(ctx, receipt.TxHash)
	if err != nil {
		return TransferRecord{}, false, err
	}
	if tx == nil || tx.To() == nil || tx.Value() == nil || tx.Value().Sign() <= 0 {
		return TransferRecord{}, false, nil
	}

	from, err := eth.TransactionSender(ctx, tx, receipt.BlockHash, uint(receipt.TransactionIndex))
	if err != nil {
		return TransferRecord{}, false, err
	}

	record := TransferRecord{
		TxHash:      receipt.TxHash.Hex(),
		BlockNumber: int64(receipt.BlockNumber.Uint64()),
		From:        from.Hex(),
		To:          tx.To().Hex(),
		Amount:      tx.Value().String(),
		AssetType:   AssetTypeETH,
		TokenSymbol: "ETH",
	}
	return record, true, nil
}

func topicToAddress(topic common.Hash) string {
	return common.BytesToAddress(topic.Bytes()[12:]).Hex()
}

func normalizeTokenSymbolMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		addr := strings.ToLower(strings.TrimSpace(k))
		symbol := strings.ToUpper(strings.TrimSpace(v))
		if addr == "" || symbol == "" {
			continue
		}
		out[addr] = symbol
	}
	return out
}
