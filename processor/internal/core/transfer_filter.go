package core

import (
	"go_bullayer_v1/base/pkg/logger"
	"strings"
)

// AssetType 资产类型
type AssetType string

const (
	AssetTypeETH   AssetType = "ETH"
	AssetTypeERC20 AssetType = "ERC20"
)

// TransferRecord 统一转账记录结构
type TransferRecord struct {
	TxHash       string
	BlockNumber  int64
	From         string
	To           string
	Amount       string
	AssetType    AssetType
	TokenAddress string // 原生 ETH 可为空
	TokenSymbol  string // 例如 ETH/USDT/BTC/WBTC
}

// TransferFilter 按地址集合和资产白名单筛选“流入”交易。
type TransferFilter struct{}

func NewTransferFilter() *TransferFilter {
	return &TransferFilter{}
}

// DefaultTrackedAssets 默认资产白名单
func DefaultTrackedAssets() []string {
	return []string{"ETH", "USDT", "BTC", "WBTC"}
}

func (f *TransferFilter) FilterIncomingTransfers(
	targetAddresses []string,
	trackedAssets []string,
	transfers []TransferRecord,
) []TransferRecord {
	if len(targetAddresses) == 0 || len(transfers) == 0 {
		return nil
	}

	addressSet := makeSet(targetAddresses)
	if len(addressSet) == 0 {
		return nil
	}

	if len(trackedAssets) == 0 {
		trackedAssets = DefaultTrackedAssets()
	}
	assetSet := makeSet(trackedAssets)
	if len(assetSet) == 0 {
		return nil
	}

	results := make([]TransferRecord, 0, len(transfers))
	for _, t := range transfers {
		to := normalize(t.To)
		if _, ok := addressSet[to]; !ok {
			continue
		}

		symbol := normalize(t.TokenSymbol)
		if symbol == "" {
			if t.AssetType == AssetTypeETH {
				symbol = "eth"
			} else {
				continue
			}
		}

		if _, ok := assetSet[symbol]; !ok {
			continue
		}
		logger.Info("过滤到交易: %+v", t)
		results = append(results, t)
	}

	return results
}

func makeSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		n := normalize(v)
		if n == "" {
			continue
		}
		set[n] = struct{}{}
	}
	return set
}

func normalize(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}
