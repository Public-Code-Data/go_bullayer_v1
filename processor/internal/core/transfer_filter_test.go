package core

import "testing"

func TestFilterIncomingTransfers(t *testing.T) {
	filter := NewTransferFilter()

	targets := []string{
		"0xAbC0000000000000000000000000000000000001",
		"0xAbC0000000000000000000000000000000000002",
	}
	assets := []string{"ETH", "USDT", "WBTC"}

	input := []TransferRecord{
		{
			TxHash:      "0x1",
			To:          "0xabc0000000000000000000000000000000000001",
			AssetType:   AssetTypeETH,
			TokenSymbol: "ETH",
		},
		{
			TxHash:      "0x2",
			To:          "0xabc0000000000000000000000000000000000002",
			AssetType:   AssetTypeERC20,
			TokenSymbol: "USDT",
		},
		{
			TxHash:      "0x3",
			To:          "0xabc0000000000000000000000000000000000002",
			AssetType:   AssetTypeERC20,
			TokenSymbol: "BTC",
		},
		{
			TxHash:      "0x4",
			To:          "0xabc0000000000000000000000000000000000003",
			AssetType:   AssetTypeERC20,
			TokenSymbol: "USDT",
		},
	}

	got := filter.FilterIncomingTransfers(targets, assets, input)
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}

	if got[0].TxHash != "0x1" || got[1].TxHash != "0x2" {
		t.Fatalf("unexpected tx order or tx hash: %+v", got)
	}
}

func TestFilterIncomingTransfers_UseDefaultAssets(t *testing.T) {
	filter := NewTransferFilter()

	targets := []string{"0xabc0000000000000000000000000000000000001"}
	input := []TransferRecord{
		{TxHash: "0x1", To: "0xabc0000000000000000000000000000000000001", AssetType: AssetTypeETH, TokenSymbol: "ETH"},
		{TxHash: "0x2", To: "0xabc0000000000000000000000000000000000001", AssetType: AssetTypeERC20, TokenSymbol: "USDT"},
		{TxHash: "0x3", To: "0xabc0000000000000000000000000000000000001", AssetType: AssetTypeERC20, TokenSymbol: "WBTC"},
		{TxHash: "0x4", To: "0xabc0000000000000000000000000000000000001", AssetType: AssetTypeERC20, TokenSymbol: "DAI"},
	}

	got := filter.FilterIncomingTransfers(targets, nil, input)
	if len(got) != 3 {
		t.Fatalf("expected 3 records with default assets, got %d", len(got))
	}
}
